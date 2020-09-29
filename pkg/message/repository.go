package message

type MessageRepository interface {
	CreateMessage(message *Message) (*Message, error)
	DeleteMessage(message *Message) []error
	GetMessagesByEmail(email string, limit int, pageState []byte) ([]*Message, []byte, error)
	GetByMagicNumber(magicNumber *Number) ([]*Message, error)
}
