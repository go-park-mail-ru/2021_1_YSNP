package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends"
	"github.com/tarantool/go-tarantool"
	"sort"
	"strconv"
	"time"
	"unicode/utf8"
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

func (tr *TrendsRepository) CreateTrendsProducts(userID uint64) error {
	val, _ := tr.dbConn.Call("get_user_trend", []interface{}{userID})
	d  := fmt.Sprintf("%v", val.Data)
	oldModel := &models.Trends{}
	json.Unmarshal([]byte(RemoveLastChar(d)), &oldModel)

	fmt.Println("RECEIVED TRENDS ", oldModel.Popular)

	selectQuery := `
				SELECT p.id
				FROM product as p
				INNER JOIN users as u ON p.owner_id=u.id and u.id<>$1
				WHERE LOWER(p.name) LIKE ANY(ARRAY[LOWER($2)`

	titleLike := ""
	var trendValue []interface{}
	trendValue = append(trendValue, userID)
	for i, trend := range oldModel.Popular {
		if i == 0 {
			trendValue = append(trendValue, "%"+trend.Title+"%")
			continue
		}
		titleLike += ", LOWER($" + strconv.Itoa(i + 2) + ")"
		trendValue = append(trendValue, "%"+trend.Title+"%")
	}
	selectQuery += titleLike
	selectQuery += `])
				ORDER BY p.date DESC 
				LIMIT 30
				`

	fmt.Println("QUERY SELECT ", selectQuery)
	fmt.Println("TRENDS LIKE ", trendValue)
	query, err := tr.dbConnPsql.Query(selectQuery, trendValue...)
	if err != nil {
		fmt.Println("ERROR", err)
		return err
	}

	var productsID []uint64
	defer query.Close()
	for query.Next() {
		var id uint64
		err := query.Scan(
			&id)
		if err != nil {
			return err
		}
		productsID = append(productsID, id)
	}
	fmt.Println("NEW SELECTED PRODUCTS ", productsID)

	oldProducts := &models.TrendProducts{}
	for i, id := range productsID {
		oldProducts.Popular[len(oldProducts.Popular) - i].ProductID = id
		oldProducts.Popular[len(oldProducts.Popular) - i].Time = time.Now()
	}

	fmt.Println("NEW TREND PRODUCTS", oldProducts.Popular)

	data, err := json.Marshal(oldProducts)

	if err != nil {
		return err
	}

	dataStr := string(data)
	resp, _ := tr.dbConn.Insert("trends_products", []interface{}{userID, dataStr})
	if resp.Code == 3 {
		val, err = tr.dbConn.Call("get_user_trends_products", []interface{}{userID})
		d  = fmt.Sprintf("%v", val.Data)
		json.Unmarshal([]byte(RemoveLastChar(d)), &oldProducts)
		fmt.Println("OLD TREND PRODUCTS", oldProducts.Popular)

		for i, item := range productsID {
			for _, prod := range oldProducts.Popular {
				if item == prod.ProductID {
					productsID = append(productsID[:i], productsID[i + 1:]...)
					prod.Time = time.Now()
				}
			}
		}
		sort.Sort(models.ProductSorter(oldProducts.Popular))
		for i, id := range productsID {
			oldProducts.Popular[len(oldProducts.Popular) - i].ProductID = id
			oldProducts.Popular[len(oldProducts.Popular) - i].Time = time.Now()
		}

		fmt.Println("NEW TREND PRODUCTS", oldProducts.Popular)

		data, err := json.Marshal(oldProducts)

		if err != nil {
			return err
		}

		dataStr := string(data)
		_, err1 := tr.dbConn.Replace("trends_products", []interface{}{userID, dataStr})

		fmt.Println("NEW TREND PRODUCTS IN STR", dataStr)
		return err1
	}
	return nil
}

func (tr *TrendsRepository) updateData(ui1 *models.Trends, ui2 *models.Trends) {
	for _, item := range ui1.Popular{
		found := false
		for i, item_2 := range ui2.Popular {
			if item.Title == item_2.Title {
				ui2.Popular[i].Count += 1
				found = true
				break
			}
			found = false
		}

		if !found {
			ui2.Popular = append(ui2.Popular, item)
		}
	}
	sort.Sort(models.PopularSorter(ui2.Popular))
	if len(ui2.Popular) > 10 {
		ui2.Popular = ui2.Popular[:10]
	}
	fmt.Println(ui2)
}

func RemoveLastChar(str string) string {
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
	fmt.Println(resp.Data...)
	if resp.Code == 3 {
		val, _ := tr.dbConn.Call("get_user_trend", []interface{}{ui.UserID})
		d  := fmt.Sprintf("%v", val.Data)
		oldModel := &models.Trends{}
		json.Unmarshal([]byte(RemoveLastChar(d)), &oldModel)
		tr.updateData(ui, oldModel)
		data, err = json.Marshal(oldModel)

		if err != nil {
			return err
		}

		dataStr = string(data)
		_, err1 := tr.dbConn.Replace("trends", []interface{}{ui.UserID, dataStr})
		return err1
	}
	return nil
}


