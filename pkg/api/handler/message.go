package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/pawel_prokop/recruitment-project-go/pkg/api/response"
	"github.com/pawel_prokop/recruitment-project-go/pkg/message"
	"gopkg.in/validator.v2"
	"net/http"
	"strconv"
)

func messageCreate(service message.MessageService) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		defer func() {
			err := request.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		messageError := "Created message error : "
		var mess *message.Message
		if err := json.NewDecoder(request.Body).Decode(&mess); err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}

		if err := validator.Validate(mess); err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}

		var err error

		mess, err = service.CreateMessage(mess)
		if err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}

		writer.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(writer).Encode(response.MessageResponse{Message: *mess}); err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}
	})
}

func sendMessages(service message.MessageService) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		defer func() {
			err := request.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		messageError := "Send messages error : "
		var number *message.Number
		if err := json.NewDecoder(request.Body).Decode(&number); err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}

		if err := validator.Validate(number); err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}

		if err := service.SendMessages(number); err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}

		writer.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(writer).Encode(response.SuccessResponse{Success: "Successfully send messages"}); err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}
	})
}

func getMessagesByEmail(service message.MessageService) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		defer func() {
			err := request.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		email := mux.Vars(request)["emailValue"]

		messageError := "Get messages error : "

		var messages []*message.Message
		var err error
		limit := 3
		limit, err = strconv.Atoi(request.FormValue("limit"))
		var pageState []byte
		if request.FormValue("pageState") != "" {
			pageState, err = hex.DecodeString(request.FormValue("pageState"))
		}

		if err != nil {
			fmt.Println(err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(messageError))
			return
		}
		messages, pageState, err = service.GetMessagesByEmail(email, limit, pageState)

		if err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}

		writer.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(writer).Encode(response.MessagesPagingResponse{Messages: messages, PageState: hex.EncodeToString(pageState)}); err != nil {
			fmt.Println(err.Error())
			returnErrorResponse(http.StatusInternalServerError, err, writer, messageError)
			return
		}
	})
}

func MakeMessageHandlers(router *mux.Router, neg negroni.Negroni, service message.MessageService) {
	router.Handle("/api/message", neg.With(
		negroni.Wrap(messageCreate(service)),
	)).Methods("POST", "OPTIONS").Name("messageCreate")

	router.Handle("/api/send", neg.With(
		negroni.Wrap(sendMessages(service)),
	)).Methods("POST", "OPTIONS").Name("sendMessages")

	router.Handle("/api/messages/{emailValue}", neg.With(
		negroni.Wrap(getMessagesByEmail(service)),
	)).Methods("GET", "OPTIONS").Name("getMessagesByEmail")
}

func returnErrorResponse(status int, err error, writer http.ResponseWriter, messageError string) {
	writer.WriteHeader(status)
	writer.Write([]byte(messageError + err.Error()))
}
