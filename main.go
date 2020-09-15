package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Message struct
type Message struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdat"`
}

// APIResponse price feed
type APIResponse struct {
	Data struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"data"`
}

// Service struct
type Service struct {
	client http.Client
}

func main() {
	app := chi.NewRouter()

	// middleware
	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(middleware.URLFormat)

	// public routes
	app.Get("/price", priceFn)
	app.Get("/", handlerFn)

	// start server
	// defer http.ListenAndServe(":3001", app)
	fmt.Println("[rei]: Server running...")
	err := http.ListenAndServe(":3001", app)
	if err != nil {
		log.Fatal("Serving: ", err)
	}
}

func handlerFn(res http.ResponseWriter, req *http.Request) {
	// create a new message
	msg := Message{
		ID:        "f80b342c-f90c-4804-9df1-faeb244ab9b8",
		Message:   "Kurama api",
		CreatedAt: time.Now(),
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(msg)
}

// Returns a list of partners
func priceFn(res http.ResponseWriter, req *http.Request) {
	var d APIResponse
	url := "https://api.coinbase.com/v2/prices/spot?currency=USD"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewDecoder(response.Body).Decode(&d)
	json.NewEncoder(res).Encode(d)
}

func (api Service) getPrice() (res APIResponse, err error) {
	url := "https://api.coinbase.com/v2/prices/spot?currency=USD"
	response, err := api.client.Get(url)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(response.Body).Decode(&res)
	return
}
