package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/codebarz/go-micorservices/data"
	"github.com/gorilla/mux"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handled GET request")

	d := data.GetProducts()

	err := d.ToJSON(w)

	if err != nil {
		http.Error(w, "Error getting products", http.StatusInternalServerError)
		return
	}
}

func (p *Product) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handled POST request")

	product := r.Context().Value(ProductContextKey{}).(*data.Product)

	p.l.Printf("JSON sent %#v", product)
	data.AddProduct(product)

}

func (p *Product) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT requests")

	pathParams := mux.Vars(r)
	id, converr := strconv.Atoi(pathParams["id"])

	p.l.Printf("The id id %v", id)

	if converr != nil {
		http.Error(w, "Error parsing id", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(ProductContextKey{}).(*data.Product)

	findErr := data.UpdateProduct(id, prod)

	if findErr == data.ErrorNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if findErr != nil {
		http.Error(w, "Error finding product", http.StatusInternalServerError)
		return
	}
}

type ProductContextKey struct{}

func (p Product) JSONValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := &data.Product{}
		err := product.FromJson(r.Body)

		if err != nil {
			errv := fmt.Sprintf("Unable to unmarshal JSON. ERROR: %v", err)
			http.Error(w, errv, http.StatusBadRequest)
			return
		}

		valerr := product.Validate()

		if valerr != nil {
			errv := fmt.Sprintf("[ERROR]: validation error. ERROR: %v", err)
			http.Error(w, errv, http.StatusBadRequest)
			return
		}

		pc := context.Background()
		ctx := context.WithValue(pc, ProductContextKey{}, product)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
