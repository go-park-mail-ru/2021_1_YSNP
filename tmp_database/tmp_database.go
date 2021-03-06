package tmp_database

import (
	"2021_1_YSNP/models"
	"errors"
	"math/rand"
	"strconv"
	"sync"
)

var newDB map[string]map[string]interface{}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

var mtx sync.Mutex

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func InitDB() {//map[string][]interface{}{
	newDB = make(map[string]map[string]interface{})

	newDB["users"] = make(map[string]interface{})
	newDB["products"] = make(map[string]interface{})
	newDB["session"] = make(map[string]interface{})

	newDB["users"]["89990009900"] = models.SignUpData{
		ID:         0,
		Name:       "Sergey",
		Surname:    "Alehin",
		Sex:        "male",
		Email:      "alehin@mail.ru",
		Telephone:  "89990009900",
		Password:   "Qwerty12",
		DateBirth:  0,
		Day:        "",
		Month:      "",
		Year:       "",
		LinkImages: nil,
	}

	newDB["products"]["0"] = models.ProductData{
		ID:          0,
		Name:        "iphone",
		Date:        2000,
		Amount:      12000,
		LinkImages:  nil,
		Description: "eto iphone",
		OwnerID:     0,
		Views:       1000,
		Likes:       20,
	}

	newDB["products"]["1"] = models.ProductData{
		ID:          1,
		Name:        "iphone 10",
		Date:        2001,
		Amount:      12001,
		LinkImages:  nil,
		Description: "eto iphone 12",
		OwnerID:     0,
		Views:       1000,
		Likes:       20,
	}
}

func checkLogin(number string) bool {
	_, exist := newDB["users"][number]
	return exist
}

func GetUserByLogin(login string) (models.LoginData,error){
	user, ok := newDB["users"][login]
	if !ok {
		return models.LoginData{}, errors.New(`no user`)
	}
	return models.LoginData{
		Telephone:  user.(models.SignUpData).Telephone,
		Password:   user.(models.SignUpData).Password,
		IsLoggedIn: false,
	}, nil
}

func getUserBySession(session string) (models.SignUpData){
	login := newDB["session"][session]
	return newDB["users"][login.(string)].(models.SignUpData)
}

func GetProducts() map[string]models.ProductListData {
	products := make(map[string]models.ProductListData)

	for k, v := range newDB["products"]{
		products[k] = models.ProductListData{
			ID:         v.(models.ProductData).ID,
			Name:       v.(models.ProductData).Name,
			Date:       v.(models.ProductData).Date,
			Amount:     v.(models.ProductData).Amount,
			LinkImages: v.(models.ProductData).LinkImages,
		}
	}
	return products
}

func GetProduct(id string) (models.ProductData,error) {
	product, ok := newDB["products"][id]
	if !ok {
		return models.ProductData{}, errors.New("no product")
	}

	return product.(models.ProductData), nil
}

func NewUser(user *models.SignUpData) error {
	defer mtx.Unlock()
	mtx.Lock()

	if checkLogin(user.Telephone) {
		return errors.New("user with this phone number exists")
	} else {
		var id uint64 = 0
		if len(newDB["users"]) > 0 {
			id = newDB["users"][strconv.Itoa(len(newDB["users"])-1)].(models.SignUpData).ID + 1
		}
		user.ID = id
		newDB["users"][user.Telephone] = user
	}
	return nil
}

func ChangeUserData (session string, newData *models.SignUpData) error{
	defer mtx.Unlock()
	mtx.Lock()

	user := getUserBySession(session)
	newData.ID = user.ID
	delete(newDB["users"], user.Telephone)
	newDB["users"][newData.Telephone] = newData
	return nil
}

func NewProduct(product *models.ProductData) error {
	defer mtx.Unlock()
	mtx.Lock()
	var id uint64 = 0
	if len(newDB["products"]) > 0 {
		id = newDB["products"][strconv.Itoa(len(newDB["products"])-1)].(models.ProductData).ID + 1
	}
	product.ID = id
	newDB["products"][strconv.Itoa(int(id))] = product
	return nil
}

func NewSession(number string) string {
	defer mtx.Unlock()
	SID := RandStringRunes(32)
	mtx.Lock()
	newDB["session"][SID] = number
	return SID
}

func CheckSession(sessionValue string) bool {
	_, auth := newDB["session"][sessionValue]
	return auth
}



