package main

import (
	"fmt"
	"net/http"
	"parking_app/store"
	"parking_app/internal"
)

type application struct {
	store store.Storage
}

func main() {
	app := &application{
		store: store.NewStorage(db.ConnectToDB()),
	}

	server := &http.Server{
		Addr:    ":3000",
		Handler: app.mount(),
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("there was an error while serving and listening!")
	}

	fmt.Println("")
}
