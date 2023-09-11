package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/fpuentem/microservices-golang/ep5/data"
	"github.com/gorilla/mux"
)

type KeyProduct struct{}

type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// // getProductID returns the product ID from the URL
// func getProductID(r *http.Request) (int, error) {
// 	// parse the product id from the url
// 	vars := mux.Vars()

// 	// Convert the id into an integer and return
// 	return strconv.Atoi(vars["id"])
// }

// getProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products request")

	// fetch the products from the store
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		p.l.Println("[ERROR] serializing products", err)
		http.Error(rw, "Unable to marshal products", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products request")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Product", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		p.l.Println("[ERROR] Unable to find product id in URL", r.URL.Path, err)
		http.Error(rw, "Missing product id, url should be formatted /products/[id] for PUT requests", http.StatusBadRequest)
		return
	}
}

// type KeyProduct struct

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserialization error", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
