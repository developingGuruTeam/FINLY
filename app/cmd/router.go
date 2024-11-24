package main

import (
	"cachManagerApp/app/internal/hendlers"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Post("/addTransaction", hendlers.PostTransaction) //добавление транзакции

	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Printf("Start server error: %s", err.Error())
		return
	}
}
