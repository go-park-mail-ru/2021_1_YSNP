package repository

import (
	"encoding/json"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"unicode/utf8"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/trends"
	"sort"
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

func (tr *TrandsRepository) InsertOrUpdate(ui *models.Trands) error {
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
		oldModel := &models.Trands{}
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


