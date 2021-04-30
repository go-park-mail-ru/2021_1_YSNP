package repository

import (
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/search"
)

func NewSearchRepository(conn *sql.DB) search.SearchRepository {
	return &SearchRepository{
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
	return data.FromAmount, maxAmount

}

func getDateSorting(data *models.Search) string {
	switch data.Date {
	case "За 24 часа":
		return "AND date BETWEEN now() - INTERVAL '2 DAY' AND now()"
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

func (s SearchRepository) SelectByFilter(userID *uint64, data *models.Search) ([]*models.ProductListData, error) {
	var products []*models.ProductListData
	minAmount, maxAmount := getMaxMinAmount(data)
	var values []interface{}
	selectQuery := `
				SELECT p.id, p.name, p.date, p.amount, array_agg(pi.img_link), uf.user_id, p.tariff
				FROM product as p
				left join product_images as pi on pi.product_id=p.id
				left join category as cat on cat.id=p.category_id
				left join user_favorite uf on p.id = uf.product_id and uf.user_id = $1
				WHERE LOWER(name) LIKE LOWER($2) AND
					  cat.title LIKE $3  AND
				      amount BETWEEN $4 AND $5
			  `
	var limit string

	if data.Radius == 0 {
		limit = "LIMIT $6 OFFSET $7"
		values = append(values, *userID, "%"+data.Search+"%", "%"+data.Category+"%", minAmount, maxAmount, data.Count, data.From*data.Count)
	} else {
		limit = "LIMIT $8 OFFSET $9"
		selectQuery += `AND
		ST_DWithin(
		  Geography(ST_SetSRID(ST_POINT(p.longitude, p.latitude), 4326)),
		  ST_GeogFromText($6), $7 * 1000)`
		val := "SRID=4326; POINT(" + fmt.Sprintf("%f", data.Longitude) + " " + fmt.Sprintf("%f", data.Latitude) + ")"
		values = append(values, *userID, "%"+data.Search+"%", "%"+data.Category+"%", minAmount, maxAmount, val, data.Radius, data.Count, data.From*data.Count)
	}

	dateQuery := getDateSorting(data)
	var orderQuery string
	orderQuery += "GROUP BY p.id, uf.user_id " + getOrderQuery(data)
	resultQuery := strings.Join([]string{
		selectQuery,
		dateQuery,
		orderQuery,
		limit,
	}, " ")
	query, err := s.dbConn.Query(resultQuery, values...)

	if err != nil {
		return nil, err
	}

	defer query.Close()

	var linkStr string
	for query.Next() {
		product := &models.ProductListData{}
		var user sql.NullInt64

		err := query.Scan(
			&product.ID,
			&product.Name,
			&product.Date,
			&product.Amount,
			&linkStr,
			&user,
			&product.Tariff,
		)

		if err != nil {
			return nil, err
		}
		product.UserLiked = false
		if userID != nil && user.Valid && uint64(user.Int64) == *userID {
			product.UserLiked = true
		}

		linkStr = linkStr[1 : len(linkStr)-1]
		if linkStr != "NULL" {
			product.LinkImages = strings.Split(linkStr, ",")
		}
		products = append(products, product)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	return products, err
}
