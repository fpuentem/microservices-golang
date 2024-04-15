package data

import "testing"

func TestChecksValidator(t *testing.T) {
	p := &Product{
		Name:  "nics",
		Price: -1.00,
		SKU:   "abs-asdas-daede",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
