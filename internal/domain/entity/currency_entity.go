package entity

type Currency struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	IsSupported bool   `json:"is_supported"`
}

var SupportedCurrencies = map[string]Currency{
	"USD": {
		Code:        "USD",
		Name:        "United States Dollar",
		Symbol:      "$",
		IsSupported: true,
	},
	"INR": {
		Code:        "INR",
		Name:        "Indian Rupee",
		Symbol:      "₹",
		IsSupported: true,
	},
	"EUR": {
		Code:        "EUR",
		Name:        "Euro",
		Symbol:      "€",
		IsSupported: true,
	},
	"JPY": {
		Code:        "JPY",
		Name:        "Japanese Yen",
		Symbol:      "¥",
		IsSupported: true,
	},
	"GBP": {
		Code:        "GBP",
		Name:        "British Pound Sterling",
		Symbol:      "£",
		IsSupported: true,
	},
}
