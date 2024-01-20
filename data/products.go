package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Products) ToJSON(io io.Writer) error {
	e := json.NewEncoder(io)
	return e.Encode(p)
}

func (p *Product) FromJson(io io.Reader) error {
	e := json.NewDecoder(io)
	return e.Decode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.Id = getNextID()
	productList = append(productList, p)
}

func getNextID() int {
	lastProduct := productList[len(productList)-1]

	return lastProduct.Id + 1
}

func UpdateProduct(id int, p *Product) error {
	fmt.Println(id)

	p, i, err := p.findProduct(id)

	if err != nil {
		return err
	}

	p.Id = id
	productList[i] = p
	return nil
}

func (p *Product) findProduct(id int) (*Product, int, error) {
	for i, prod := range productList {
		if id == p.Id {
			return prod, i, nil
		}

	}

	return nil, -1, ErrorNotFound
}

var ErrorNotFound = fmt.Errorf("Product not found")

var productList = []*Product{
	&Product{
		Id:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		DeletedOn:   time.Now().UTC().String(),
	},
	&Product{
		Id:          2,
		Name:        "Esspresso",
		Description: "Short and string coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		DeletedOn:   time.Now().UTC().String(),
	},
}
