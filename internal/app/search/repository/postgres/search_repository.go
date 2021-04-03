package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search"
	"math"
	"strings"
)

type SearchRepository struct {
	dbConn *sql.DB
}

func (s SearchRepository) SelectByFilter(data *models.Search) ([]*models.ProductListData, error) {
	var products []*models.ProductListData

	minAmount := data.FromAmount
	maxAmount := data.ToAmount

	if maxAmount == 0 {
		maxAmount = math.MaxInt32
	}

	var values []interface{}
	selectQuery := `
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link)
				FROM product as p
				left join product_images as pi on pi.product_id=p.id
				WHERE name LIKE $1 AND
				      category LIKE $2  AND
				      amount BETWEEN $3 AND $4   `
	values = append(values, "%" + data.Search + "%", "%" + data.Category + "%", minAmount, maxAmount)
	var pgntQuery string
	switch data.Date {
		case "За 24 часа":
			pgntQuery = "AND date BETWEEN now() - INTERVAL '1 DAY' AND now()"
		case "За 7 дней":
			pgntQuery = "AND date BETWEEN now() - INTERVAL '7 DAY' AND now()"
		default:
			break
	}
	var orderQuery string
	orderQuery += "GROUP BY p.id "
	switch data.Sorting {
	case "По возрастанию цены":
		orderQuery += "ORDER BY amount "
	case "По убыванию цены":
		orderQuery += "ORDER BY amount DESC"
	case "По имени":
		orderQuery += "ORDER BY name"
	case "По дате добавления":
		orderQuery += "ORDER BY date"
	default:
		break
	}
	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
		orderQuery,
	}, " ")
	query, err := s.dbConn.Query(resultQuery, values...)

	if err != nil {
		return nil, err
	}

	defer query.Close()

	//var linkStr string

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

func NewProductRepository(conn *sql.DB) search.SearchRepository {
	return & SearchRepository{
		dbConn: conn,
	}
}