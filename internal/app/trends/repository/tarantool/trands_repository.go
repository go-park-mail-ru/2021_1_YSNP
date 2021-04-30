package repository

import (
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

type TrandsRepository struct {
	dbConn *tarantool.Connection
}

func NewTrandsRepository(conn *tarantool.Connection) trends.TrandsRepository {
	return &TrandsRepository{
		dbConn: conn,
	}
}

func (tr *TrandsRepository) updateData(ui1 *models.Trands, ui2 *models.Trands) {
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
		if ui2.Popular[i].Date.Unix() < time.Now().Add(-10 * time.Hour).Unix() {
			ui2.Popular = remove(ui2.Popular, i)
			i -= 1
		} 
		i += 1
	}

	sort.Sort(models.PopularSorter(ui2.Popular))


	fmt.Println(ui2)
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

func (tr *TrandsRepository) InsertOrUpdate(ui *models.Trands) error {
	data, err := json.Marshal(ui)
	
	if err != nil {
		return err
	}

	dataStr := string(data)

	resp, _ := tr.dbConn.Insert("trends", []interface{}{ui.UserID, dataStr})

	if resp.Code == 3 {
		val, _ := tr.dbConn.Call("get_user_trend", []interface{}{ui.UserID})
		d  := fmt.Sprintf("%v", val.Data)
		oldModel := &models.Trands{}
		json.Unmarshal([]byte(removeLastChar(d)), &oldModel)
		fmt.Println(oldModel)
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
