package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/pawel_prokop/recruitment-project-go/pkg/message"
)

type Repository struct {
	session *gocql.Session
}

func NewCassandraRepository(session *gocql.Session) message.MessageRepository {
	return &Repository{
		session: session,
	}
}

func (r *Repository) CreateMessage(message *message.Message) (*message.Message, error) {
	query := r.session.Query(`INSERT INTO recruitment.message (id, email, magic_number, title, content) VALUES (?, ?, ?, ?, ?) USING TTL ?`,
		message.Id, message.Email, message.MagicNumber, message.Title, message.Content, 300)
	queryByMagicNumber := r.session.Query(`INSERT INTO recruitment.message_by_magic_number (id, email, magic_number, title, content) VALUES (?, ?, ?, ?, ?) USING TTL ?`,
		message.Id, message.Email, message.MagicNumber, message.Title, message.Content, 300)
	if err := query.Exec(); err != nil {
		return message, err
	}
	if err := queryByMagicNumber.Exec(); err != nil {
		return message, err
	}
	return message, nil
}

func (r *Repository) DeleteMessage(message *message.Message) []error {
	var errs []error
	query := r.session.Query(`DELETE FROM recruitment.message WHERE email = ? AND id = ?`, message)
	queryByMagicNumber := r.session.Query(`DELETE FROM recruitment.message_by_magic_number WHERE magic_number = ? AND id = ?`, message)

	if err := query.Exec(); err != nil {
		errs = append(errs, err)
		fmt.Println("Delete Error", err)
	}

	if err := queryByMagicNumber.Exec(); err != nil {
		errs = append(errs, err)
		fmt.Println("Delete Error", err)
	}

	if errs != nil {
		fmt.Println("Deleting Errors ! : ", errs)
		return errs
	} else {
		fmt.Println("Success Delete! ID : ", message.Id)
		return nil
	}
}

func (r *Repository) GetByMagicNumber(magicNumber *message.Number) ([]*message.Message, error) {
	var messages []*message.Message
	query := r.session.Query(`SELECT * FROM recruitment.message_by_magic_number WHERE magic_number = ?`, magicNumber.MagicNumber)

	iterMap := map[string]interface{}{}

	iter := query.Iter()
	for iter.MapScan(iterMap) {
		messages = append(messages, &message.Message{
			Id:          iterMap["id"].(gocql.UUID),
			Email:       iterMap["email"].(string),
			MagicNumber: iterMap["magic_number"].(int),
			Title:       iterMap["title"].(string),
			Content:     iterMap["content"].(string),
		})

		iterMap = map[string]interface{}{}
	}
	return messages, nil
}

func (r *Repository) GetMessagesByEmail(email string, limit int, pageState []byte) ([]*message.Message, []byte, error) {
	var messages []*message.Message

	query := r.session.Query(`SELECT * FROM recruitment.message WHERE email = ?`, email)

	rows := query.Iter().NumRows()
	query.PageSize(limit)
	query.PageState(pageState)

	iter := query.Iter()
	iterMap := map[string]interface{}{}

	for iter.MapScan(iterMap) {
		messages = append(messages, &message.Message{
			Id:          iterMap["id"].(gocql.UUID),
			Email:       iterMap["email"].(string),
			MagicNumber: iterMap["magic_number"].(int),
			Title:       iterMap["title"].(string),
			Content:     iterMap["content"].(string),
		})
		iterMap = map[string]interface{}{}
	}
	pageState = iter.PageState()

	if err := iter.Close(); err != nil {
		return nil, nil, err
	}
	if rows == limit {
		pageState = []byte("")
	}

	return messages, pageState, nil
}
