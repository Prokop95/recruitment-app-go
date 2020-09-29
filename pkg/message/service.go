package message

import (
	"fmt"
	"github.com/pawel_prokop/recruitment-project-go/config"
	"net/smtp"
)

type MessageService interface {
	CreateMessage(m *Message) (*Message, error)
	SendMessages(magicNumber *Number) error
	GetMessagesByEmail(email string, limit int, pageState []byte) ([]*Message, []byte, error)
}

type Service struct {
	repository MessageRepository
	config     *config.Config
}

func NewService(repository MessageRepository, config *config.Config) MessageService {
	return &Service{
		repository: repository,
		config:     config,
	}
}

func (s *Service) CreateMessage(message *Message) (*Message, error) {
	message.Id = NewID()
	return s.repository.CreateMessage(message)
}

func (s *Service) SendMessages(magicNumber *Number) error {
	messages, err := s.repository.GetByMagicNumber(magicNumber)

	if err != nil {
		fmt.Println("Find by magic number Error", err)
		return err
	}

	if len(messages) > 0 {
		for i := range messages {
			s.sendEmail(messages[i])
			s.repository.DeleteMessage(messages[i])
		}
	}
	return nil
}

func (s *Service) FindByMagicNumber(magicNumber *Number) ([]*Message, error) {
	messages, err := s.repository.GetByMagicNumber(magicNumber)

	if err != nil {
		fmt.Println("Find by magic number Error", err)
		return messages, nil
	}

	return messages, nil
}

func (s *Service) DeleteMessage(message *Message) []error {
	return s.repository.DeleteMessage(message)
}

func (s *Service) GetMessagesByEmail(email string, limit int, pageState []byte) ([]*Message, []byte, error) {
	return s.repository.GetMessagesByEmail(email, limit, pageState)
}

func (s *Service) sendEmail(message *Message) {
	auth := smtp.PlainAuth("", s.config.Mail.User, s.config.Mail.Password, s.config.Mail.Smtp.Host)

	to := []string{message.Email}
	msg := []byte("To: " + message.Email +
		"Subject: " + message.Title +
		message.Content)
	err := smtp.SendMail(s.config.Mail.Smtp.Host+":"+s.config.Mail.Smtp.Port, auth, s.config.Mail.User, to, msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email send !")
	}
}
