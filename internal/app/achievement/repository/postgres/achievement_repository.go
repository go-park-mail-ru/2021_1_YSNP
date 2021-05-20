package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/achievement"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
)

type AchievementRepository struct {
	dbConn *sql.DB
}

func NewAchievementRepository(conn *sql.DB) achievement.AchievementRepository {
	return &AchievementRepository{
		dbConn: conn,
	}
}

func(ar *AchievementRepository) GetUserAchievements(userId int) ([]*models.Achievement, error) {
	req, err := ar.dbConn.Query(`
		SELECT a.title, a.description, ua.date, a.link_pic
		FROM user_achievement ua
		JOIN achievement a ON ua.a_id = a.id
		WHERE ua.user_id = $1
		ORDER BY ua.date
	`, userId)
	if err != nil  {
		return nil, err
	}

	defer req.Close()

	achievements := make([]*models.Achievement, 0)
	for req.Next() {
		achievement := &models.Achievement{}

		err := req.Scan(
			&achievement.Titie,
			&achievement.Description,
			&achievement.Date, 
			&achievement.LinkPic,
		)

		if err != nil {
			return nil, err
		}

		achievements = append(achievements, achievement)
	}

	if err := req.Err(); err != nil {
		return nil, err
	}
	return achievements, err
}
