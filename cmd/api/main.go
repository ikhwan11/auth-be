package main

import (
	"log"
	"net/http"

	"github.com/ikhwan11/main-be/internal/router"
)

func main() {
	r := router.Setup()

	log.Println("🚀 Server running on :8090")

	if err := http.ListenAndServe(":8090", r); err != nil {
		log.Fatal(err)
	}
}
