package entity

type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

var SupportedCurrencies = map[string]Currency{
	"USD": {
		Code:   "USD",
		Name:   "United States Dollar",
		Symbol: "$",
	},
	"INR": {
		Code:   "INR",
		Name:   "Indian Rupee",
		Symbol: "₹",
	},
	"EUR": {
		Code:   "EUR",
		Name:   "Euro",
		Symbol: "€",
	},
	"JPY": {
		Code:   "JPY",
		Name:   "Japanese Yen",
		Symbol: "¥",
	},
	"GBP": {
		Code:   "GBP",
		Name:   "British Pound Sterling",
		Symbol: "£",
	},
}
