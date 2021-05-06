package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends"
	"github.com/tarantool/go-tarantool"
)

type TrendsRepository struct {
	dbConn *tarantool.Connection
	dbConnPsql *sql.DB
}

func NewTrendsRepository(conn *tarantool.Connection, conn2 *sql.DB) trends.TrendsRepository {
	return &TrendsRepository{
		dbConn: conn,
		dbConnPsql: conn2,
	}
}

func (tr *TrendsRepository) getProductIDForTrend(userID uint64, title string, limit int, offset int) ([]uint64, error) {
	selectQuery := `
				SELECT p.id
				FROM product as p
				INNER JOIN users as u ON p.owner_id=u.id and u.id<>$1
				WHERE LOWER(p.name) LIKE LOWER($2)
				ORDER BY p.date DESC 
				LIMIT $3 OFFSET $4`

	fmt.Println("SELECT", selectQuery, userID, title, limit, offset)
	query, err := tr.dbConnPsql.Query(selectQuery, userID, title, limit, offset)
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

func (tr *TrendsRepository) getNewTrendsProductID(userID uint64, oldModel models.Trends, reshuffle bool) ([]uint64, error) {
	sort.Sort(models.PopularSorter(oldModel.Popular))
	var productsID []uint64
	limit := 30 / len(oldModel.Popular)
	if limit <= 1 {
		limit = 1
	}
	offset := 0
	if reshuffle {
		offset = int(30 % oldModel.Popular[0].Count)
	}

	for _, trend := range oldModel.Popular {
		trendsID, err := tr.getProductIDForTrend(userID, trend.Title, limit, offset)
		if err != nil {
			return nil, err
		}
		fmt.Println("EACH TREND", trendsID)
		for _, id := range trendsID {
			productsID = append(productsID, id)
		}
	}
	return productsID, nil
}

func (tr *TrendsRepository) CreateTrendsProducts(userID uint64) error {
	val, _ := tr.dbConn.Call("get_user_trend", []interface{}{userID})
	d  := fmt.Sprintf("%v", val.Data)
	oldModel := &models.Trends{}
	json.Unmarshal([]byte(removeLastChar(d)), &oldModel)

	fmt.Println("TRENDS", oldModel.Popular)

	productsID, err := tr.getNewTrendsProductID(userID, *oldModel, false)
	if err != nil {
		return err
	}
	fmt.Println("ProductsID", productsID)

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
		d  = fmt.Sprintf("%v", val.Data)
		json.Unmarshal([]byte(removeLastChar(d)), &oldProducts)


		oldProducts, err = tr.replaceNewTrends(productsID, *oldProducts, userID, *oldModel)
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

func (tr *TrendsRepository) replaceNewTrends(productsID []uint64, oldProducts models.TrendProducts, userID uint64, oldTrends models.Trends) (*models.TrendProducts, error) {
	for i := 0; i < len(productsID); i++ {
		for _, prod := range oldProducts.Popular {
			if productsID[i] == prod.ProductID {
				prod.Time = time.Now()
				if i + 1 != len(productsID) {
					productsID = append(productsID[:i], productsID[i+1:]...)
					i -= 1
					break
				}
			}
		}
	}
	var err error
	if len(productsID) <= 3 {
		productsID, err = tr.getNewTrendsProductID(userID, oldTrends, true)
		if err != nil {
			return nil, err
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
			oldProducts.Popular[len(oldProducts.Popular)- i - 1].ProductID = id
			oldProducts.Popular[len(oldProducts.Popular)- i - 1].Time = time.Now()
		}
	}
	oldProducts.UserID = userID
	return &oldProducts, nil
}

func (tr *TrendsRepository) GetTrendsProducts(userID uint64) ([]uint64, error) {
	products := &models.TrendProducts{}
	val, err := tr.dbConn.Call("get_user_trends_products", []interface{}{userID})
	if err != nil {
		return nil, err
	}
	d  := fmt.Sprintf("%v", val.Data)
	json.Unmarshal([]byte(removeLastChar(d)), &products)

	var productsID []uint64
	for _, val := range products.Popular {
		productsID = append(productsID, val.ProductID)
	}
	return productsID, nil
}

func (tr *TrendsRepository) updateData(ui1 *models.Trends, ui2 *models.Trends) {
	for _, item := range ui1.Popular{
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
		if ui2.Popular[i].Date.Unix() < time.Now().Add(-10 * time.Hour).Unix() {
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
              return str[2:len(str)-size - 1]
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
		d  := fmt.Sprintf("%v", val.Data)
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
