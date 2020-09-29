package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/pawel_prokop/recruitment-project-go/config"
	"github.com/pawel_prokop/recruitment-project-go/pkg/api/handler"
	"github.com/pawel_prokop/recruitment-project-go/pkg/datasource/cassandra"
	"github.com/pawel_prokop/recruitment-project-go/pkg/message"
	"github.com/pawel_prokop/recruitment-project-go/pkg/middleware"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := mux.NewRouter()

	configuration := config.GetConf()
	cluster := cassandra.InitCassandra(configuration)
	messageRepo := cassandra.NewCassandraRepository(cluster)
	messageService := message.NewService(messageRepo, configuration)

	n := negroni.New(
		negroni.HandlerFunc(middleware.Cors),
		negroni.NewLogger(),
	)

	handler.MakeMessageHandlers(r, *n, messageService)

	http.Handle("/", r)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(configuration.App.Port),
		Handler:      context.ClearHandler(http.DefaultServeMux),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
