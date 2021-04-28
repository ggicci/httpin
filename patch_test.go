package httpin_test

import "github.com/ggicci/httpin"

type ProductPatch struct {
	Title    httpin.String `json:"title"`
	Color    httpin.String `json:"color"`
	Quantity httpin.Int    `json:"quantity"`
}
