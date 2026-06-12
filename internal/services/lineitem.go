package services

import "encoding/json"

type LineItem struct {
	ItemID      int64   `json:"item_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	Taxable     bool    `json:"taxable"`
	TaxRate     string  `json:"tax_rate"`
	Discount    float64 `json:"discount"`
	Surcharge   float64 `json:"surcharge"`
}

func ParseLineItems(s string) ([]LineItem, error) {
	if s == "" || s == "[]" {
		return nil, nil
	}
	var items []LineItem
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return nil, err
	}
	return items, nil
}

func SerializeLineItems(items []LineItem) string {
	if len(items) == 0 {
		return "[]"
	}
	b, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(b)
}
