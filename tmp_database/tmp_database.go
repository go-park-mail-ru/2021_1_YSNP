package tmp_database

import (
	"2021_1_YSNP/models"
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"time"
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

func InitDB() { //map[string][]interface{}{
	newDB = make(map[string]map[string]interface{})

	newDB["users"] = make(map[string]interface{})
	newDB["products"] = make(map[string]interface{})
	newDB["session"] = make(map[string]interface{})

	newDB["users"]["+79990009900"] = models.SignUpData{
		ID:         0,
		Name:       "Sergey",
		Surname:    "Alehin",
		Sex:        "male",
		Email:      "alehin@mail.ru",
		Telephone:  "+79990009900",
		Password:   "Qwerty12",
		DateBirth:  "",
		LinkImages: []string{"http://89.208.199.170:8080/static/avatar/test-avatar.jpg"},
	}

	newDB["products"]["0"] = models.ProductData{
		ID:           0,
		Name:         "iphone",
		Date:         "2000-10-10",
		Amount:       12000,
		LinkImages:   []string{"http://89.208.199.170:8080/static/product/pic4.jpeg", "http://89.208.199.170:8080/static/product/pic5.jpeg", "http://89.208.199.170:8080/static/product/pic6.jpeg"},
		Description:  "eto iphone",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        1000,
		Likes:        20,
	}

	newDB["products"]["1"] = models.ProductData{
		ID:           1,
		Name:         "iphone 10",
		Date:         "2000-10-10",
		Amount:       12001,
		LinkImages:   []string{"http://89.208.199.170:8080/static/product/pic1.jpeg", "http://89.208.199.170:8080/static/product/pic2.jpeg", "http://89.208.199.170:8080/static/product/pic3.jpeg"},
		Description:  "eto iphone 12",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        1000,
		Likes:        20,
	}
}

func checkLogin(number string) bool {
	_, exist := newDB["users"][number]
	return exist
}

func GetUserByLogin(login string) (models.LoginData, error) {
	user, ok := newDB["users"][login]
	if !ok {
		return models.LoginData{}, errors.New(`No user with this number`)
	}
	return models.LoginData{
		Telephone:  user.(models.SignUpData).Telephone,
		Password:   user.(models.SignUpData).Password,
		IsLoggedIn: false,
	}, nil
}

func GetUserBySession(session string) models.SignUpData {
	login := newDB["session"][session]
	return newDB["users"][login.(string)].(models.SignUpData)
}

func GetProducts() map[string][]models.ProductListData {
	products := make(map[string][]models.ProductListData)

	for _, v := range newDB["products"] {
		products["product_list"] = append(products["product_list"], models.ProductListData{
			ID:         v.(models.ProductData).ID,
			Name:       v.(models.ProductData).Name,
			Date:       v.(models.ProductData).Date,
			Amount:     v.(models.ProductData).Amount,
			LinkImages: v.(models.ProductData).LinkImages,
		})
	}
	return products
}

func GetProduct(id string) (models.ProductData, error) {
	product, ok := newDB["products"][id]
	if !ok {
		return models.ProductData{}, errors.New("No product with this id.")
	}

	return product.(models.ProductData), nil
}

func NewUser(user *models.SignUpData) error {
	defer mtx.Unlock()
	mtx.Lock()

	if checkLogin(user.Telephone) {
		return errors.New("User with this phone number exists.")
	} else {
		// var id uint64 = 0
		// if len(newDB["users"]) > 0 {
		id := RandStringRunes(32)
		// }
		user.ID, _ = strconv.ParseUint(id, 10, 64)
		newDB["users"][user.Telephone] = *user
	}
	return nil
}

func ChangeUserData(session string, newData *models.SignUpData) error {
	defer mtx.Unlock()
	mtx.Lock()
	user := GetUserBySession(session)
	newData.ID = user.ID
	delete(newDB["users"], user.Telephone)
	newDB["users"][newData.Telephone] = newData
	return nil
}

func NewProduct(product *models.ProductData, session string) error {
	defer mtx.Unlock()
	mtx.Lock()
	var id uint64 = 0
	if len(newDB["products"]) > 0 {
		id = newDB["products"][strconv.Itoa(len(newDB["products"])-1)].(models.ProductData).ID + 1
	}
	product.ID = id
	user := GetUserBySession(session)
	product.OwnerID = user.ID
	product.OwnerName = user.Name
	product.OwnerSurname = user.Surname
	product.Date = time.Now().UTC().String()
	newDB["products"][strconv.Itoa(int(id))] = *product
	return nil
}

func NewSession(number string) string {
	defer mtx.Unlock()
	SID := RandStringRunes(32)
	mtx.Lock()
	newDB["session"][SID] = number
	return SID
}

func DeleteSession(session string) {
	delete(newDB["session"], session)
}

func CheckSession(sessionValue string) bool {
	_, auth := newDB["session"][sessionValue]
	return auth
}
