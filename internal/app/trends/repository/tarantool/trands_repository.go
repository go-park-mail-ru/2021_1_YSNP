package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends"
	"github.com/tarantool/go-tarantool"
)

type TrendsRepository struct {
	dbConn     *tarantool.Connection
	dbConnPsql *sql.DB
}

func NewTrendsRepository(conn *tarantool.Connection, conn2 *sql.DB) trends.TrendsRepository {
	return &TrendsRepository{
		dbConn:     conn,
		dbConnPsql: conn2,
	}
}

func (tr *TrendsRepository) getNewTrendsProductID(userID uint64, oldModel models.Trends, orderBy string, categoryID uint64) ([]uint64, error) {
	selectQuery := `
				SELECT p.id
				FROM product as p
				INNER JOIN users as u ON p.owner_id=u.id and u.id<>$1
				WHERE `
	var step int

	var trendValue []interface{}
	trendValue = append(trendValue, userID)
	if categoryID != 0 {
		selectQuery += `p.category_id = $2 and p.close = false and LOWER(p.name) LIKE ANY(ARRAY[LOWER($3)`
		trendValue = append(trendValue, categoryID)
		step = 3
	} else {
		selectQuery += `p.close = false and LOWER(p.name) LIKE ANY(ARRAY[LOWER($2)`
		step = 2
	}

	titleLike := ""
	for i, trend := range oldModel.Popular {
		if i == 0 {
			trendValue = append(trendValue, "%"+trend.Title+"%")
			continue
		}
		titleLike += ", LOWER($" + strconv.Itoa(i+step) + ")"
		trendValue = append(trendValue, "%"+trend.Title+"%")
	}
	selectQuery += titleLike
	if orderBy == "date" {
		selectQuery += `])
				ORDER BY p.date DESC 
				LIMIT 30
				`
	} else if orderBy == "likes" {
		selectQuery += `])
				ORDER BY p.likes DESC 
				LIMIT 30
				`
	} else {
		selectQuery += `]) 
				LIMIT 30
				`
	}

	query, err := tr.dbConnPsql.Query(selectQuery, trendValue...)
	if err != nil {
		return nil, err
	}

	var productsID []uint64
	defer query.Close()
	for query.Next() {
		var id uint64
		err := query.Scan(
			&id)
		if err != nil {
			return nil, err
		}
		productsID = append(productsID, id)
	}
	return productsID, nil
}

func (tr *TrendsRepository) CreateTrendsProducts(userID uint64) error {
	val, _ := tr.dbConn.Call("get_user_trend", []interface{}{userID})
	d := fmt.Sprintf("%v", val.Data)
	oldModel := &models.Trends{}
	json.Unmarshal([]byte(removeLastChar(d)), &oldModel)

	productsID, err := tr.getNewTrendsProductID(userID, *oldModel, "date", 0)
	if err != nil {
		return err
	}

	if len(productsID) < 10 {
		index := len(productsID)
		var pId uint64
		b := true
		if index < 1 {
			if len(oldModel.Popular) > 0 {
				query := tr.dbConnPsql.QueryRow(
					`
				SELECT p.id
				FROM product AS p
				WHERE LOWER(p.name) LIKE LOWER($1)
				`,
					"%"+oldModel.Popular[rand.Intn(len(oldModel.Popular))].Title+"%")
				err := query.Scan(&pId)
				if err != nil {
					return err
				}
			} else {
				b = false
			}
		} else {
			pId = productsID[rand.Intn(len(productsID))]
		}
		if b {
			addProd, err := tr.getProdIdSameCategory(pId, userID)
			if err != nil {
				return err
			}
			for _, val := range addProd {
				productsID = append(productsID, val)
			}
		}
	}

	oldProducts := &models.TrendProducts{}
	for _, id := range productsID {
		var prod models.PopularProduct
		prod.ProductID = id
		prod.Time = time.Now()
		oldProducts.Popular = append(oldProducts.Popular, prod)
	}

	data, err := json.Marshal(oldProducts)

	if err != nil {
		return err
	}

	dataStr := string(data)
	resp, _ := tr.dbConn.Insert("trends_products", []interface{}{userID, dataStr})
	if resp.Code == 3 {
		val, err = tr.dbConn.Call("get_user_trends_products", []interface{}{userID})
		d = fmt.Sprintf("%v", val.Data)
		json.Unmarshal([]byte(removeLastChar(d)), &oldProducts)

		oldProducts, err = replaceNewTrends(productsID, *oldProducts, userID)
		if err != nil {
			return err
		}

		data, err := json.Marshal(oldProducts)

		if err != nil {
			return err
		}

		dataStr := string(data)
		_, err1 := tr.dbConn.Replace("trends_products", []interface{}{userID, dataStr})

		return err1
	}
	return nil
}

func (tr *TrendsRepository) getProdIdSameCategory(productID uint64, userID uint64) ([]uint64, error) {
	query := tr.dbConnPsql.QueryRow(
		`
				SELECT p.category_id
				FROM product AS p
				WHERE p.id=$1
				`,
		productID)
	var categoryId uint64
	err := query.Scan(&categoryId)
	if err != nil {
		return nil, err
	}
	selectQuery := `
				SELECT p.id
				FROM product as p
				INNER JOIN users as u ON p.owner_id=u.id and u.id<>$1
				WHERE p.category_id=$2 and p.id<>$3
				LIMIT 5`

	var trendValue []interface{}
	trendValue = append(trendValue, userID)
	trendValue = append(trendValue, categoryId)
	trendValue = append(trendValue, productID)
	queryProd, err := tr.dbConnPsql.Query(selectQuery, trendValue...)
	if err != nil {
		return nil, err
	}

	var productsID []uint64
	defer queryProd.Close()
	for queryProd.Next() {
		var id uint64
		err := queryProd.Scan(
			&id)
		if err != nil {
			return nil, err
		}
		productsID = append(productsID, id)
	}
	return productsID, nil
}

func replaceNewTrends(productsID []uint64, oldProducts models.TrendProducts, userID uint64) (*models.TrendProducts, error) {
	for i := 0; i < len(productsID); i++ {
		for _, prod := range oldProducts.Popular {
			if productsID[i] == prod.ProductID {
				prod.Time = time.Now()
				if i+1 != len(productsID) {
					productsID = append(productsID[:i], productsID[i+1:]...)
					i -= 1
					break
				}
			}
		}
	}
	sort.Sort(models.ProductSorter(oldProducts.Popular))
	for i, id := range productsID {
		if i >= len(oldProducts.Popular) || len(oldProducts.Popular) < 30 {
			var prod models.PopularProduct
			prod.ProductID = id
			prod.Time = time.Now()
			oldProducts.Popular = append(oldProducts.Popular, prod)
		} else {
			oldProducts.Popular[len(oldProducts.Popular)-i-1].ProductID = id
			oldProducts.Popular[len(oldProducts.Popular)-i-1].Time = time.Now()
		}
	}
	oldProducts.UserID = userID
	return &oldProducts, nil
}

func (tr *TrendsRepository) checkSuffix(word string) bool {
	if len(word) < 5 {
		return true
	}
	stop := []string{"ими", "ыми", "его", "ого", "ему", "ому", "ее", "ие",
		"ые", "ое", "ей", "ий", "ый", "ой", "ем", "им", "ым",
		"ом", "их", "ых", "ую", "юю", "ая", "яя", "ою", "ею"}

	for _, item := range stop {
		suf := word[len(word)-len(item):]
		if suf == item {
			return false
		}
	}
	return true
}

func (tr *TrendsRepository) GetRecommendationProducts(productID uint64, userID uint64) ([]uint64, error) {
	query := tr.dbConnPsql.QueryRow(
		`
				SELECT p.name, p.category_id
				FROM product AS p
				WHERE p.id=$1
				`,
		productID)
	var title string
	var categoryId uint64
	err := query.Scan(&title, &categoryId)
	if err != nil {
		return nil, err
	}

	var words models.Trends
	words.UserID = userID
	cleanContent := stopwords.CleanString(title, "ru", true)
	sn := strings.TrimSpace(cleanContent)
	s := strings.FieldsFunc(sn, func(r rune) bool { return strings.ContainsRune(" .,:-", r) })
	for _, item := range s {
		if !tr.checkSuffix(item) {
			continue
		}

		stemmed, err := snowball.Stem(item, "russian", true)
		if err == nil {
			words.Popular = append(words.Popular, models.Popular{
				Title: stemmed,
				Count: 1,
				Date:  time.Now(),
			})
		}
	}

	productsID, err := tr.getNewTrendsProductID(userID, words, "likes", categoryId)
	if err != nil {
		return nil, err
	}
	productsID = append(productsID)
	if len(productsID) < 10 {
		addProd, err := tr.getProdIdSameCategory(productID, userID)
		if err != nil {
			return nil, err
		}
		for _, val := range addProd {
			productsID = append(productsID, val)
		}
	}
	return productsID, nil
}

func (tr *TrendsRepository) GetTrendsProducts(userID uint64) ([]uint64, error) {
	products := &models.TrendProducts{}
	val, err := tr.dbConn.Call("get_user_trends_products", []interface{}{userID})
	if err != nil {
		return nil, err
	}
	d := fmt.Sprintf("%v", val.Data)
	json.Unmarshal([]byte(removeLastChar(d)), &products)

	var productsID []uint64
	for _, val := range products.Popular {
		productsID = append(productsID, val.ProductID)
	}
	return productsID, nil
}

func (tr *TrendsRepository) updateData(ui1 *models.Trends, ui2 *models.Trends) {
	for _, item := range ui1.Popular {
		found := false
		for i, item_2 := range ui2.Popular {
			if item.Title == item_2.Title {
				ui2.Popular[i].Count += 1
				ui2.Popular[i].Date = item.Date
				found = true
				break
			}
			found = false
		}

		if !found {
			ui2.Popular = append(ui2.Popular, item)
		}
	}

	i := 0
	for i < len(ui2.Popular) {
		if ui2.Popular[i].Date.Unix() < time.Now().Add(-10*time.Hour).Unix() {
			ui2.Popular = remove(ui2.Popular, i)
			i -= 1
		}
		i += 1
	}

	sort.Sort(models.PopularSorter(ui2.Popular))
}

func remove(slice []models.Popular, i int) []models.Popular {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func removeLastChar(str string) string {
	for len(str) > 0 {
		_, size := utf8.DecodeLastRuneInString(str)
		return str[2 : len(str)-size-1]
	}
	return str
}

func (tr *TrendsRepository) InsertOrUpdate(ui *models.Trends) error {
	data, err := json.Marshal(ui)

	if err != nil {
		return err
	}

	dataStr := string(data)

	resp, _ := tr.dbConn.Insert("trends", []interface{}{ui.UserID, dataStr})

	if resp.Code == 3 {
		val, _ := tr.dbConn.Call("get_user_trend", []interface{}{ui.UserID})
		d := fmt.Sprintf("%v", val.Data)
		oldModel := &models.Trends{}
		json.Unmarshal([]byte(removeLastChar(d)), &oldModel)

		tr.updateData(ui, oldModel)

		data, err = json.Marshal(oldModel)

		if err != nil {
			return err
		}

		dataStr = string(data)
		_, err1 := tr.dbConn.Replace("trends", []interface{}{ui.UserID, dataStr})
		return err1
	}
	return errors.New(resp.Error)
}
