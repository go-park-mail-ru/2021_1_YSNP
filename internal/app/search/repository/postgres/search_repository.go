package repository

import (
	"database/sql"
	"math"
	"strings"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search"
)

func NewProductRepository(conn *sql.DB) search.SearchRepository {
	return & SearchRepository{
		dbConn: conn,
	}
}

type SearchRepository struct {
	dbConn *sql.DB
}

func getMaxMinAmount(data *models.Search) (int, int) {
	maxAmount := data.ToAmount

	if maxAmount == 0 {
		maxAmount = math.MaxInt32
	}
	return  data.FromAmount, maxAmount

}

func getDateSorting(data *models.Search) string {
	switch data.Date {
		case "За 24 часа":
			return "AND date BETWEEN now() - INTERVAL '1 DAY' AND now()"
		case "За 7 дней":
			return "AND date BETWEEN now() - INTERVAL '7 DAY' AND now()"
		default:
			return ""
	}
}

func getOrderQuery(data *models.Search) string {
	switch data.Sorting {
	case "По возрастанию цены":
		return "ORDER BY amount "
	case "По убыванию цены":
		return "ORDER BY amount DESC"
	case "По имени":
		return "ORDER BY name"
	case "По дате добавления":
		return "ORDER BY date"
	default:
		return ""
	}
}

func (s SearchRepository) SelectByFilter(data *models.Search) ([]*models.ProductListData, error) {
	var products []*models.ProductListData
	minAmount, maxAmount := getMaxMinAmount(data);
	var values []interface{}
	selectQuery := `
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link)
				FROM product as p
				left join product_images as pi on pi.product_id=p.id
				left join category as cat on cat.id=p.category_id
				WHERE LOWER(name) LIKE LOWER($1) AND
					  cat.title LIKE $2  AND
				      amount BETWEEN $3 AND $4   `
	values = append(values, "%" + data.Search + "%", "%" + data.Category + "%", minAmount, maxAmount)
	dateQuery := getDateSorting(data)
	var orderQuery string
	orderQuery += "GROUP BY p.id " + getOrderQuery(data)
	resultQuery := strings.Join([]string{
		selectQuery,
		dateQuery,
		orderQuery,
	}, " ")
	query, err := s.dbConn.Query(resultQuery, values...)

	if err != nil {
		return nil, err
	}

	defer query.Close()

	var linkStr string
	for query.Next() {
		product := &models.ProductListData{}

		err := query.Scan(
			&product.ID,
			&product.Name,
			&product.Date,
			&product.Amount,
			&linkStr,
			)

		if err != nil {
			return nil, err
		}

		products = append(products, product)

		linkStr = linkStr[1 : len(linkStr)-1]
		if linkStr != "NULL" {
			product.LinkImages = strings.Split(linkStr, ",")
		}
	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	return products, err
}
