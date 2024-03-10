package app

import (
	"log"
	"net/http"
	"shop/internal/pkg/item/delivery"
	"shop/internal/pkg/item/usecase"
)

func App() {
	items := usecase.NewInmemory()
	itemHandler := delivery.New(items)
	http.Handle("/items/", itemHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
