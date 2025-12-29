package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/Adi-ty/betteruptime-go/internal/app"
	"github.com/Adi-ty/betteruptime-go/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run server on")

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	app.Logger.Println("Started application")

	router := routes.SetUpRoutes(app)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout: time.Minute,
	}


	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Printf("Server failed to start: %v\n", err)
	}
}