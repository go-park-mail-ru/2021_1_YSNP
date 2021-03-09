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

//const Url = "http://89.208.199.170:8080"
const Url = "http://localhost:8080"

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
		DateBirth:  "1991-11-11",
		LinkImages: []string{Url + "/static/avatar/test-avatar.jpg"},
	}

	newDB["products"]["0"] = models.ProductData{
		ID:           0,
		Name:         "iphone",
		Date:         "2000-10-10",
		Amount:       5994,
		LinkImages:   []string{Url + "/static/product/pic10.jpeg", Url + "/static/product/pic7.jpeg", Url + "/static/product/pic3.jpeg"},
		Description:  "Ясность нашей позиции очевидна: перспективное планирование играет определяющее значение для благоприятных перспектив. Противоположная точка зрения подразумевает, что сторонники тоталитаризма в науке неоднозначны и будут объективно рассмотрены соответствующими инстанциями.",
		Category:     "Автомобили",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        23,
		Likes:        634,
	}

	newDB["products"]["1"] = models.ProductData{
		ID:           1,
		Name:         "iphone 10",
		Date:         "2000-10-10",
		Amount:       68788,
		LinkImages:   []string{Url + "/static/product/pic15.jpeg", Url + "/static/product/pic6.jpeg", Url + "/static/product/pic8.jpeg", Url + "/static/product/pic3.jpeg"},
		Description:  "Безусловно, сложившаяся структура организации обеспечивает широкому кругу (специалистов) участие в формировании форм воздействия. Высокий уровень вовлечения представителей целевой аудитории является четким доказательством простого факта: граница обучения кадров обеспечивает актуальность кластеризации усилий. Таким образом, глубокий уровень погружения позволяет выполнить важные задания по разработке первоочередных требований. Таким образом, реализация намеченных плановых заданий создаёт необходимость включения в производственный план целого ряда внеочередных мероприятий с учётом комплекса переосмысления внешнеэкономических политик. Приятно, граждане, наблюдать, как многие известные личности заблокированы в рамках своих собственных рациональных ограничений.",
		Category:     "Автомобили",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        653,
		Likes:        234,
	}

	newDB["products"]["2"] = models.ProductData{
		ID:           2,
		Name:         "Загородный дом",
		Date:         "2012-1-12",
		Amount:       78538,
		LinkImages:   []string{Url + "/static/product/pic11.jpeg", Url + "/static/product/pic5.jpeg", Url + "/static/product/pic14.jpeg", Url + "/static/product/pic13.jpeg", Url + "/static/product/pic6.jpeg"},
		Description:  "Повседневная практика показывает, что разбавленное изрядной долей эмпатии, рациональное мышление способствует повышению качества экономической целесообразности принимаемых решений. Противоположная точка зрения подразумевает, что действия представителей оппозиции формируют глобальную экономическую сеть и при этом - превращены в посмешище, хотя само их существование приносит несомненную пользу обществу. В целом, конечно, новая модель организационной деятельности предполагает независимые способы реализации экспериментов, поражающих по своей масштабности и грандиозности. Современные технологии достигли такого уровня, что перспективное планирование, в своём классическом представлении, допускает внедрение вывода текущих активов.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        234,
		Likes:        54,
	}

	newDB["products"]["3"] = models.ProductData{
		ID:           3,
		Name:         "AirPods",
		Date:         "2012-1-12",
		Amount:       23579,
		LinkImages:   []string{Url + "/static/product/pic3.jpeg", Url + "/static/product/pic7.jpeg", Url + "/static/product/pic12.jpeg"},
		Description:  "Задача организации, в особенности же синтетическое тестирование требует анализа приоритизации разума над эмоциями. Противоположная точка зрения подразумевает, что сторонники тоталитаризма в науке набирают популярность среди определенных слоев населения, а значит, должны быть объективно рассмотрены соответствующими инстанциями. В своём стремлении улучшить пользовательский опыт мы упускаем, что непосредственные участники технического прогресса объединены в целые кластеры себе подобных.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        5432,
		Likes:        432,
	}

	newDB["products"]["4"] = models.ProductData{
		ID:           4,
		Name:         "Наушники",
		Date:         "2012-1-12",
		Amount:       79297,
		LinkImages:   []string{Url + "/static/product/pic12.jpeg", Url + "/static/product/pic3.jpeg", Url + "/static/product/pic9.jpeg", Url + "/static/product/pic6.jpeg"},
		Description:  "Есть над чем задуматься: активно развивающиеся страны третьего мира, вне зависимости от их уровня, должны быть объективно рассмотрены соответствующими инстанциями! Значимость этих проблем настолько очевидна, что высококачественный прототип будущего проекта не даёт нам иного выбора, кроме определения направлений прогрессивного развития. Новая модель организационной деятельности создаёт предпосылки для экспериментов, поражающих по своей масштабности и грандиозности.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        635,
		Likes:        23,
	}

	newDB["products"]["5"] = models.ProductData{
		ID:           5,
		Name:         "Звонилка",
		Date:         "2012-1-12",
		Amount:       56644,
		LinkImages:   []string{Url + "/static/product/pic14.jpeg", Url + "/static/product/pic2.jpeg", Url + "/static/product/pic4.jpeg"},
		Description:  "Также как понимание сути ресурсосберегающих технологий однозначно определяет каждого участника как способного принимать собственные решения касаемо благоприятных перспектив. Господа, повышение уровня гражданского сознания является качественно новой ступенью поставленных обществом задач. В своём стремлении повысить качество жизни, они забывают, что граница обучения кадров требует определения и уточнения благоприятных перспектив.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        43,
		Likes:        654,
	}

	newDB["products"]["6"] = models.ProductData{
		ID:           6,
		Name:         "Онлайн магазин",
		Date:         "2012-1-12",
		Amount:       45339,
		LinkImages:   []string{Url + "/static/product/pic13.jpeg", Url + "/static/product/pic14.jpeg", Url + "/static/product/pic7.jpeg", Url + "/static/product/pic12.jpeg"},
		Description:  "Господа, современная методология разработки не даёт нам иного выбора, кроме определения экспериментов, поражающих по своей масштабности и грандиозности. Следует отметить, что дальнейшее развитие различных форм деятельности играет важную роль в формировании приоритизации разума над эмоциями. Каждый из нас понимает очевидную вещь: глубокий уровень погружения не оставляет шанса для направлений прогрессивного развития!",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        10432000,
		Likes:        54,
	}

	newDB["products"]["7"] = models.ProductData{
		ID:           7,
		Name:         "Телефон",
		Date:         "2012-1-12",
		Amount:       83753,
		LinkImages:   []string{Url + "/static/product/pic9.jpeg", Url + "/static/product/pic6.jpeg", Url + "/static/product/pic5.jpeg"},
		Description:  "Разнообразный и богатый опыт говорит нам, что повышение уровня гражданского сознания играет важную роль в формировании форм воздействия. А также элементы политического процесса, которые представляют собой яркий пример континентально-европейского типа политической культуры, будут функционально разнесены на независимые элементы. Предварительные выводы неутешительны: внедрение современных методик прекрасно подходит для реализации распределения внутренних резервов и ресурсов.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        543,
		Likes:        123,
	}

	newDB["products"]["8"] = models.ProductData{
		ID:           8,
		Name:         "Тоже гараж",
		Date:         "2012-1-12",
		Amount:       98873,
		LinkImages:   []string{Url + "/static/product/pic12.jpeg", Url + "/static/product/pic7.jpeg", Url + "/static/product/pic6.jpeg", Url + "/static/product/pic10.jpeg", Url + "/static/product/pic1.jpeg"},
		Description:  "И нет сомнений, что диаграммы связей, превозмогая сложившуюся непростую экономическую ситуацию, призваны к ответу. С другой стороны, реализация намеченных плановых заданий предоставляет широкие возможности для стандартных подходов. С учётом сложившейся международной обстановки, убеждённость некоторых оппонентов требует определения и уточнения поставленных обществом задач.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        103200,
		Likes:        432,
	}

	newDB["products"]["9"] = models.ProductData{
		ID:           9,
		Name:         "Гараж",
		Date:         "2012-1-12",
		Amount:       75204,
		LinkImages:   []string{Url + "/static/product/pic8.jpeg", Url + "/static/product/pic11.jpeg", Url + "/static/product/pic5.jpeg", Url + "/static/product/pic13.jpeg"},
		Description:  "Но дальнейшее развитие различных форм деятельности напрямую зависит от стандартных подходов. В целом, конечно, реализация намеченных плановых заданий влечет за собой процесс внедрения и модернизации системы массового участия. Принимая во внимание показатели успешности, постоянный количественный рост и сфера нашей активности обеспечивает широкому кругу (специалистов) участие в формировании новых принципов формирования материально-технической и кадровой базы.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        876,
		Likes:        543,
	}

	newDB["products"]["10"] = models.ProductData{
		ID:           10,
		Name:         "iPhone",
		Date:         "2012-1-12",
		Amount:       66085,
		LinkImages:   []string{Url + "/static/product/pic10.jpeg", Url + "/static/product/pic2.jpeg", Url + "/static/product/pic5.jpeg", Url + "/static/product/pic6.jpeg", Url + "/static/product/pic14.jpeg", Url + "/static/product/pic7.jpeg"},
		Description:  "Но глубокий уровень погружения не оставляет шанса для поставленных обществом задач! А ещё многие известные личности набирают популярность среди определенных слоев населения, а значит, должны быть объективно рассмотрены соответствующими инстанциями. Равным образом, консультация с широким активом напрямую зависит от инновационных методов управления процессами.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        234,
		Likes:        6754,
	}

	newDB["products"]["11"] = models.ProductData{
		ID:           11,
		Name:         "Просто телефон",
		Date:         "2012-1-12",
		Amount:       9591,
		LinkImages:   []string{Url + "/static/product/pic1.jpeg", Url + "/static/product/pic6.jpeg", Url + "/static/product/pic2.jpeg", Url + "/static/product/pic3.jpeg"},
		Description:  "Мы вынуждены отталкиваться от того, что существующая теория обеспечивает актуальность новых принципов формирования материально-технической и кадровой базы. В частности, укрепление и развитие внутренней структуры является качественно новой ступенью переосмысления внешнеэкономических политик. В своём стремлении улучшить пользовательский опыт мы упускаем, что интерактивные прототипы неоднозначны и будут преданы социально-демократической анафеме.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        3123,
		Likes:        78,
	}

	newDB["products"]["12"] = models.ProductData{
		ID:           12,
		Name:         "Четкий велосипед",
		Date:         "2012-1-12",
		Amount:       73001,
		LinkImages:   []string{Url + "/static/product/pic13.jpeg", Url + "/static/product/pic6.jpeg", Url + "/static/product/pic9.jpeg", Url + "/static/product/pic13.jpeg"},
		Description:  "Прежде всего, консультация с широким активом влечет за собой процесс внедрения и модернизации новых принципов формирования материально-технической и кадровой базы. Банальные, но неопровержимые выводы, а также некоторые особенности внутренней политики формируют глобальную экономическую сеть и при этом - разоблачены. Равным образом, синтетическое тестирование позволяет оценить значение поставленных обществом задач.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        674,
		Likes:        123,
	}

	newDB["products"]["13"] = models.ProductData{
		ID:           13,
		Name:         "Класная машинка",
		Date:         "2012-1-12",
		Amount:       12453,
		LinkImages:   []string{Url + "/static/product/pic1.jpeg", Url + "/static/product/pic2.jpeg", Url + "/static/product/pic6.jpeg", Url + "/static/product/pic6.jpeg", Url + "/static/product/pic3.jpeg"},
		Description:  "Высокий уровень вовлечения представителей целевой аудитории является четким доказательством простого факта: разбавленное изрядной долей эмпатии, рациональное мышление представляет собой интересный эксперимент проверки стандартных подходов. В своём стремлении улучшить пользовательский опыт мы упускаем, что представители современных социальных резервов могут быть объективно рассмотрены соответствующими инстанциями. Задача организации, в особенности же граница обучения кадров, в своём классическом представлении, допускает внедрение направлений прогрессивного развития.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        54252,
		Likes:        123,
	}

	newDB["products"]["14"] = models.ProductData{
		ID:           14,
		Name:         "Красивый дом",
		Date:         "2012-1-12",
		Amount:       1247564,
		LinkImages:   []string{Url + "/static/product/pic11.jpeg", Url + "/static/product/pic2.jpeg", Url + "/static/product/pic1.jpeg", Url + "/static/product/pic3.jpeg"},
		Description:  "Сложно сказать, почему многие известные личности, инициированные исключительно синтетически, в равной степени предоставлены сами себе. Таким образом, внедрение современных методик прекрасно подходит для реализации приоритизации разума над эмоциями. Как принято считать, непосредственные участники технического прогресса, которые представляют собой яркий пример континентально-европейского типа политической культуры, будут ассоциативно распределены по отраслям.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        124,
		Likes:        523523,
	}

	newDB["products"]["15"] = models.ProductData{
		ID:           15,
		Name:         "Шуба",
		Date:         "2012-1-12",
		Amount:       1704,
		LinkImages:   []string{Url + "/static/product/pic1.jpeg", Url + "/static/product/pic3.jpeg", Url + "/static/product/pic7.jpeg"},
		Description:  "Но элементы политического процесса призваны к ответу. Являясь всего лишь частью общей картины, некоторые особенности внутренней политики представляют собой не что иное, как квинтэссенцию победы маркетинга над разумом и должны быть своевременно верифицированы! Учитывая ключевые сценарии поведения, дальнейшее развитие различных форм деятельности позволяет выполнить важные задания по разработке поэтапного и последовательного развития общества.",
		Category:     "Хобби",
		OwnerID:      0,
		OwnerName:    "Sergey",
		OwnerSurname: "Alehin",
		Views:        1234,
		Likes:        65,
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
		id := RandStringRunes(32)
		user.ID, _ = strconv.ParseUint(id, 10, 64)
		newDB["users"][user.Telephone] = *user
	}
	return nil
}

func ChangeUserPassword(session string, newData *models.PasswordChange) error {
	defer mtx.Unlock()
	mtx.Lock()
	user := GetUserBySession(session)
	if newData.OldPassword == user.Password {
		user.Password = newData.NewPassword
		newDB["users"][user.Telephone] = user
		return nil
	} else {
		return errors.New("Old password didn't match.")
	}
}

func ChangeUserData(session string, newData *models.SignUpData) error {
	defer mtx.Unlock()
	mtx.Lock()
	user := GetUserBySession(session)
	newData.ID = user.ID
	newData.Password = user.Password
	delete(newDB["users"], user.Telephone)
	newDB["users"][newData.Telephone] = *newData
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
	product.Date = time.Now().UTC().Format("2006-01-02")
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
