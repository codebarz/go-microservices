package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("Handled PUT request")

		pattern := regexp.MustCompile(`/([0-9]+)`)
		g := pattern.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(w, "Here Invalid URI", http.StatusBadRequest)
			return
		}
		p.l.Println(g[0])
		if len(g[0]) != 2 {
			http.Error(w, "There Invalid URI", http.StatusBadRequest)
			return
		}

		idSring := g[0][1]
		id, err := strconv.Atoi(idSring)

		if err != nil {
			error := fmt.Sprintf(" ERROR: %v", err)
			http.Error(w, "Invalid URI"+error, http.StatusBadRequest)
			return
		}

		p.l.Printf("The ID is %v", id)

		p.updateProduct(id, w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handled GET request")

	d := data.GetProducts()

	err := d.ToJSON(w)

	if err != nil {
		http.Error(w, "Error getting products", http.StatusInternalServerError)
		return
	}
}

func (p *Product) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handled POST request")

	product := &data.Product{}
	err := product.FromJson(r.Body)

	if err != nil {
		errv := fmt.Sprintf("Unable to unmarshal JSON. ERROR: %v", err)
		http.Error(w, errv, http.StatusBadRequest)
		return
	}

	p.l.Printf("JSON sent %#v", product)
	data.AddProduct(product)

}

func (p *Product) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT requests")

	prod := &data.Product{}

	err := prod.FromJson(r.Body)

	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

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
