package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Tega",
		Price: 10,
		SKU:   "a-c-d",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
