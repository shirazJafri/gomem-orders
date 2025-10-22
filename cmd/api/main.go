package main

import (
	"log"
	"net/http"

	"github.com/shirazJafri/gomem-orders/handlers"
	"github.com/shirazJafri/gomem-orders/models"
)

func setupRoutes(store *models.Store) {
	orderHandler := handlers.NewOrderHandler(store)

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		orderHandler.List(w, r)
	})
}

func main() {
	store := models.NewStore()
	setupRoutes(store)

	log.Println("Server starting on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
