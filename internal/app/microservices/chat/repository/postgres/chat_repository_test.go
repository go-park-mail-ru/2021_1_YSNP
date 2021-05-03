package postgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_1_YSNP/internal/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var chatTest = &models.Chat{
	ID:                1,
	CreationTime:     	time.Now(),
	LastMsgID:         2,
	LastMsgContent:    "",
	LastMsgTime:       time.Time{},
	PartnerID:         3,
	PartnerName:       "a",
	PartnerSurname:    "s",
	PartnerAvatarLink: "d",
	ProductID:         4,
	ProductName:       "f",
	ProductAmount:     1000,
	ProductAvatarLink: "g",
	LastReadMsgId:     5,
	NewMessages:       9,
}

var msgTest = &models.Message{
	ID:           15,
	Content:      "dsds",
	CreationTime: time.Now(),
	ChatID:       1,
	UserID:       10,
}

var userID uint64 = 10;

func TestChatRepository_InsertChat(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	chatRepo := NewChatRepository(db)

	mock.ExpectBegin()
	answer := sqlmock.NewRows([]string{"id", "creation_time"}).AddRow(chatTest.ID, chatTest.CreationTime)
	mock.ExpectQuery(`INSERT INTO chats`).WithArgs(sqlmock.AnyArg()).WillReturnRows(answer)
	mock.ExpectExec(`INSERT INTO user_chats`).WithArgs(userID, chatTest.PartnerID, chatTest.ProductID, chatTest.ID).WillReturnResult(sqlmock.NewResult(int64(chatTest.ID), 1))
	rows := sqlmock.NewRows([]string{"p.name", "p.amount", "pi.img_link", "u.name", "u.surname", "u.avatar"})
	rows.AddRow(
		&chatTest.ProductName,
		&chatTest.ProductAmount,
		&chatTest.ProductAvatarLink,
		&chatTest.PartnerName,
		&chatTest.PartnerSurname,
		&chatTest.PartnerAvatarLink)
	mock.ExpectQuery(`SELECT`).WithArgs(chatTest.ProductID, chatTest.PartnerID).WillReturnRows(rows)
	mock.ExpectCommit()

	err = chatRepo.InsertChat(chatTest, userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//begin err
	mock.ExpectBegin().WillReturnError(sqlmock.ErrCancelled)

	err = chatRepo.InsertChat(chatTest, userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//inser chats err
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO chats`).WithArgs(sqlmock.AnyArg()).WillReturnError(sqlmock.ErrCancelled)

	err = chatRepo.InsertChat(chatTest, userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//insert user_chats err
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO chats`).WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id", "creation_time"}).AddRow(chatTest.ID, chatTest.CreationTime))
	mock.ExpectExec(`INSERT INTO user_chats`).WithArgs(userID, chatTest.PartnerID, chatTest.ProductID, chatTest.ID).WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))


	err = chatRepo.InsertChat(chatTest, userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//select err
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO chats`).WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id", "creation_time"}).AddRow(chatTest.ID, chatTest.CreationTime))
	mock.ExpectExec(`INSERT INTO user_chats`).WithArgs(userID, chatTest.PartnerID, chatTest.ProductID, chatTest.ID).WillReturnResult(sqlmock.NewResult(int64(chatTest.ID), 1))
	mock.ExpectQuery(`SELECT`).WithArgs(chatTest.ProductID, chatTest.PartnerID).WillReturnError(sqlmock.ErrCancelled)

	err = chatRepo.InsertChat(chatTest, userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_CheckChatExist(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	chatRepo := NewChatRepository(db)


	rows := sqlmock.NewRows([]string{"c.id", "c.creation_time", "c.last_msg_id", "c.last_msg_content", "c.last_msg_time", "u.id", "u.name", "u.surname", "u.avatar", "p.id", "p.name", "p.amount", "uc.new_messages", "pi.img_link"})
	rows.AddRow(
		&chatTest.ID,
		&chatTest.CreationTime,
		&chatTest.LastMsgID,
		&chatTest.LastMsgContent,
		&chatTest.LastMsgTime,
		&chatTest.PartnerID,
		&chatTest.PartnerName,
		&chatTest.PartnerSurname,
		&chatTest.PartnerAvatarLink,
		&chatTest.ProductID,
		&chatTest.ProductName,
		&chatTest.ProductAmount,
		&chatTest.NewMessages,
		&chatTest.ProductAvatarLink,)
	mock.ExpectQuery(`SELECT`).WithArgs(userID, chatTest.PartnerID, chatTest.ProductID).WillReturnRows(rows)

	err = chatRepo.CheckChatExist(chatTest, userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//error
	mock.ExpectQuery(`SELECT`).WithArgs(userID, chatTest.PartnerID, chatTest.ProductID).WillReturnError(sqlmock.ErrCancelled)

	err = chatRepo.CheckChatExist(chatTest, userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetChatById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	chatRepo := NewChatRepository(db)

	rows := sqlmock.NewRows([]string{"c.id", "c.creation_time", "c.last_msg_id", "c.last_msg_content", "c.last_msg_time", "u.id", "u.name", "u.surname", "u.avatar", "p.id", "p.name", "p.amount", "uc.new_messages", "pi.img_link"})
	rows.AddRow(
		&chatTest.ID,
		&chatTest.CreationTime,
		&chatTest.LastMsgID,
		&chatTest.LastMsgContent,
		&chatTest.LastMsgTime,
		&chatTest.PartnerID,
		&chatTest.PartnerName,
		&chatTest.PartnerSurname,
		&chatTest.PartnerAvatarLink,
		&chatTest.ProductID,
		&chatTest.ProductName,
		&chatTest.ProductAmount,
		&chatTest.NewMessages,
		&chatTest.ProductAvatarLink,)
	mock.ExpectQuery(`SELECT`).WithArgs(userID, chatTest.ID).WillReturnRows(rows)

	_, err = chatRepo.GetChatById(chatTest.ID, userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//error
	mock.ExpectQuery(`SELECT`).WithArgs(userID, chatTest.ID).WillReturnError(sqlmock.ErrCancelled)

	_, err = chatRepo.GetChatById(chatTest.ID, userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetUserChats(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	chatRepo := NewChatRepository(db)

	rows := sqlmock.NewRows([]string{"c.id", "c.creation_time", "c.last_msg_id", "c.last_msg_content", "c.last_msg_time", "u.id", "u.name", "u.surname", "u.avatar", "p.id", "p.name", "p.amount", "uc.new_messages", "pi.img_link"})
	rows.AddRow(
		&chatTest.ID,
		&chatTest.CreationTime,
		&chatTest.LastMsgID,
		&chatTest.LastMsgContent,
		&chatTest.LastMsgTime,
		&chatTest.PartnerID,
		&chatTest.PartnerName,
		&chatTest.PartnerSurname,
		&chatTest.PartnerAvatarLink,
		&chatTest.ProductID,
		&chatTest.ProductName,
		&chatTest.ProductAmount,
		&chatTest.NewMessages,
		&chatTest.ProductAvatarLink,)
	mock.ExpectQuery(`WITH ORDERED AS`).WithArgs(userID).WillReturnRows(rows)

	_, err = chatRepo.GetUserChats(userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//error
	mock.ExpectQuery(`WITH ORDERED AS`).WithArgs(userID).WillReturnError(sql.ErrNoRows)

	_, err = chatRepo.GetUserChats(userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_InsertMessage(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	chatRepo := NewChatRepository(db)

	mock.ExpectQuery(`INSERT INTO messages`).WithArgs(msgTest.Content, sqlmock.AnyArg(), chatTest.ID, userID).WillReturnRows(sqlmock.NewRows([]string{"id", "content", "creation_time", "chat_id", "user_id"}).AddRow(msgTest.ID, msgTest.Content, msgTest.CreationTime, msgTest.ChatID, msgTest.UserID))

	_, err = chatRepo.InsertMessage(&models.CreateMessageReq{
		ChatID:  chatTest.ID,
		Content: msgTest.Content,
	}, userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//error
	mock.ExpectQuery(`INSERT INTO messages`).WithArgs(msgTest.Content, sqlmock.AnyArg(), chatTest.ID, userID).WillReturnError(sql.ErrConnDone)

	_, err = chatRepo.InsertMessage(&models.CreateMessageReq{
		ChatID:  chatTest.ID,
		Content: msgTest.Content,
	}, userID)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetLastNMessages(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	chatRepo := NewChatRepository(db)
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "content", "creation_time", "chat_id", "user_id"})
	rows.AddRow(
		&msgTest.ID, &msgTest.Content, &msgTest.CreationTime, &msgTest.ChatID, &msgTest.UserID)
	mock.ExpectQuery(`SELECT`).WithArgs(msgTest.ChatID, 100).WillReturnRows(rows)
	mock.ExpectExec(`UPDATE user_chats`).WithArgs(userID, msgTest.ChatID).WillReturnResult(sqlmock.NewResult(int64(chatTest.ID), 1))
	mock.ExpectCommit()

	_, err = chatRepo.GetLastNMessages(&models.GetLastNMessagesReq{
		UserID: userID,
		ChatID: chatTest.ID,
		Count:  100,
	})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//select error
	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT`).WithArgs(msgTest.ChatID, 100).WillReturnError(sql.ErrConnDone)

	_, err = chatRepo.GetLastNMessages(&models.GetLastNMessagesReq{
		UserID: userID,
		ChatID: chatTest.ID,
		Count:  100,
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//update error
	mock.ExpectBegin()
	rows.AddRow(
		&msgTest.ID, &msgTest.Content, &msgTest.CreationTime, &msgTest.ChatID, &msgTest.UserID)
	mock.ExpectQuery(`SELECT`).WithArgs(msgTest.ChatID, 100).WillReturnRows(rows)
	mock.ExpectExec(`UPDATE user_chats`).WithArgs(userID, msgTest.ChatID).WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

	_, err = chatRepo.GetLastNMessages(&models.GetLastNMessagesReq{
		UserID: userID,
		ChatID: chatTest.ID,
		Count:  100,
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetNMessagesBefore(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	chatRepo := NewChatRepository(db)

	rows := sqlmock.NewRows([]string{"id", "content", "creation_time", "chat_id", "user_id"})
	rows.AddRow(
		&msgTest.ID, &msgTest.Content, &msgTest.CreationTime, &msgTest.ChatID, &msgTest.UserID)
	mock.ExpectQuery(`SELECT`).WithArgs(chatTest.ID, 100, 1000).WillReturnRows(rows)

	_, err = chatRepo.GetNMessagesBefore(&models.GetNMessagesBeforeReq{
		ChatID:        chatTest.ID,
		Count:         1000,
		LastMessageID: 100,
	})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	//error
	mock.ExpectQuery(`SELECT`).WithArgs(chatTest.ID, 100, 1000).WillReturnError(sql.ErrConnDone)

	_, err = chatRepo.GetNMessagesBefore(&models.GetNMessagesBeforeReq{
		ChatID:        chatTest.ID,
		Count:         1000,
		LastMessageID: 100,
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
