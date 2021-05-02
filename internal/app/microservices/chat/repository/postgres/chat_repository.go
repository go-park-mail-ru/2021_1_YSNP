package postgres

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/microservices/chat"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"time"
)

type ChatRepository struct {
	dbConn *sql.DB
}

func NewChatRepository (conn *sql.DB) chat.ChatRepository {
	return &ChatRepository{
		dbConn: conn,
	}
}

func (c *ChatRepository) InsertChat(chat *models.Chat, userID uint64) error {
	tx, err := c.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	query := tx.QueryRow(
		`
				INSERT INTO chats (creation_time) 
			VALUES ($1) 
			RETURNING id, creation_time; `,
		time.Now())

	err = query.Scan(&chat.ID, &chat.CreationTime)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO user_chats (user_id, partner_id, product_id, chat_id)  
                VALUES ($1, $2, $3, $4), ($2, $1, $3, $4)`,
		userID, chat.PartnerID, chat.ProductID, chat.ID)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	query = tx.QueryRow(
		`
SELECT p.name, pi.img_link, u.name, u.surname, u.avatar
				FROM product AS p
				LEFT JOIN users AS u on p.owner_id = u.id and u.id = $2
Left Join product_images pi on p.id = pi.product_id
				WHERE p.id=$1
LIMIT 1`,
		chat.ProductID, chat.PartnerID)

	err = query.Scan(&chat.ProductName, &chat.ProductAvatarLink, &chat.PartnerName, &chat.PartnerSurname, &chat.PartnerAvatarLink)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatRepository) GetChatById(chatID uint64, userID uint64) (*models.Chat, error) {
	chat := &models.Chat{}

	//добавить ссылку на аватар товара
	query := c.dbConn.QueryRow(
		`
				SELECT c.id, c.creation_time, c.last_msg_id, c.last_msg_content, c.last_msg_time, u.name, u.surname, u.avatar, p.name
				FROM chats AS c 
				JOIN user_chats AS uc ON c.id = uc.chat_id
				JOIN users AS u ON u.id = uc.partner_id
				JOIN product AS p ON p.id = uc.product_id
				WHERE uc.user_id = $1 and c.id = $2`,
				userID, chatID)

	err := query.Scan(
		&chat.ID,
		&chat.CreationTime,
		&chat.LastMsgID,
		&chat.LastMsgContent,
		&chat.LastMsgTime,
		&chat.PartnerName,
		&chat.PartnerSurname,
		&chat.PartnerAvatarLink,
		&chat.ProductName,
		)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (c *ChatRepository) GetUserChats(userID uint64) ([]*models.Chat, error) {
	query, err := c.dbConn.Query(
		`
				SELECT c.id, c.creation_time, c.last_msg_id, c.last_msg_content, c.last_msg_time, u.name, u.surname, u.avatar, p.name
				FROM chats AS c 
				JOIN user_chats AS uc ON c.id = uc.chat_id
				JOIN users AS u ON u.id = uc.partner_id
				JOIN product AS p ON p.id = uc.product_id
				WHERE uc.user_id = $1
				ORDER BY c.last_msg_time DESC, c.creation_time DESC`,
		userID)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	chats := []*models.Chat{}
	for query.Next() {
		chat := &models.Chat{}
		err := query.Scan(
			&chat.ID,
			&chat.CreationTime,
			&chat.LastMsgID,
			&chat.LastMsgContent,
			&chat.LastMsgTime,
			&chat.PartnerName,
			&chat.PartnerSurname,
			&chat.PartnerAvatarLink,
			&chat.ProductName,
		)
		if err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}
	if err := query.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (c *ChatRepository) InsertMessage(req *models.CreateMessageReq, userId uint64) (*models.Message, error) {
	msg := &models.Message{}
	query := c.dbConn.QueryRow(
		`INSERT INTO messages (content, creation_time, chat_id, user_id)
				VALUES ($1, $2, $3, $4)
				RETURNING id, content, creation_time, chat_id, user_id`,
				req.Content, time.Now(), req.ChatID, userId)

	err := query.Scan(&msg.ID, &msg.Content, &msg.CreationTime, &msg.ChatID, &msg.UserID)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (c *ChatRepository) GetLastNMessages(req *models.GetLastNMessagesReq) ([]*models.Message, error) {
	tx, err := c.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	query, err := tx.Query(
		`SELECT id, content, creation_time, chat_id, user_id
				FROM messages
				WHERE chat_id = $1
				ORDER BY id DESC 
				LIMIT $2`,
				req.ChatID, req.Count)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	defer query.Close()

	msgs := []*models.Message{}
	for query.Next() {
		msg := &models.Message{}
		err := query.Scan(&msg.ID, &msg.Content, &msg.CreationTime, &msg.ChatID, &msg.UserID)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	if err := query.Err(); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	_, err = tx.Exec(
		"UPDATE user_chats "+
			"SET new_messages = 0, "+
			"last_read_msg_id = chats.last_msg_id "+
			"FROM chats "+
			"WHERE chat_id = chats.id AND "+
			"user_id = $1 AND "+
			"chat_id = $2 ",
		req.UserID, req.ChatID)

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (c *ChatRepository) GetNMessagesBefore(req *models.GetNMessagesBeforeReq) ([]*models.Message, error) {
	query, err := c.dbConn.Query(
		`SELECT id, content, creation_time, chat_id, user_id
				FROM messages
				WHERE chat_id = $1 and id < $2
				ORDER BY id DESC 
				LIMIT $3`,
		req.ChatID, req.LastMessageID, req.Count)
	if err != nil {
		return nil, err
	}

	defer query.Close()

	msgs := []*models.Message{}
	for query.Next() {
		msg := &models.Message{}
		err := query.Scan(&msg.ID, &msg.Content, &msg.CreationTime, &msg.ChatID, &msg.UserID)
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	return msgs, nil
}


