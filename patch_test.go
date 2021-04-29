package httpin

type ProductPatch struct {
	Title    String `json:"title"`
	Color    String `json:"color"`
	Quantity Int    `json:"quantity"`
}
