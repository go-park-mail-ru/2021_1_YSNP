package repository

import (
	"context"
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

func(ar *AchievementRepository) GetUserAchievements(userId int, loggedUser int) ([]*models.Achievement, error) {
	tx, err := ar.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	req, err := tx.Query(`
		SELECT a.title, a.description, ua.date, a.link_pic
		FROM user_achievement ua
		JOIN achievement a ON ua.a_id = a.id
		WHERE ua.user_id = $1
		ORDER BY ua.date
	`, userId)
	if err != nil  {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
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
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}

		achievements = append(achievements, achievement)
	}

	if err := req.Err(); err != nil {
		return nil, err
	}

	if( loggedUser == userId){
		_, err = tx.Exec(
			"UPDATE users "+
				"SET new_achive = 0 "+
				"WHERE"+
				"id = $1",
			userId)

		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return achievements, err
}
