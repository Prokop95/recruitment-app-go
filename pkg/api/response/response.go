package response

import (
	"github.com/pawel_prokop/recruitment-project-go/pkg/message"
)

type SuccessResponse struct {
	Success string `json:"status"`
}

type MessagesPagingResponse struct {
	Messages  []*message.Message `json:"messages"`
	PageState string             `json:"pageState"`
}

type MessageResponse struct {
	Message message.Message `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
