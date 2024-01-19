package handlers

import (
	"log"
	"net/http"

	"github.com/codebarz/go-micorservices/data"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getProducts(w http.ResponseWriter, r *http.Request) {
	d := data.GetProducts()

	err := d.ToJSON(w)

	if err != nil {
		http.Error(w, "Error getting products", http.StatusInternalServerError)
		return
	}
}
