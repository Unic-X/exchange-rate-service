package entity

type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Type   string `json:"type"`
}

var SupportedCurrencies = map[string]Currency{
	"USD": {
		Code:   "USD",
		Name:   "United States Dollar",
		Symbol: "$",
		Type:   "fiat",
	},
	"INR": {
		Code:   "INR",
		Name:   "Indian Rupee",
		Symbol: "₹",
		Type:   "fiat",
	},
	"EUR": {
		Code:   "EUR",
		Name:   "Euro",
		Symbol: "€",
		Type:   "fiat",
	},
	"JPY": {
		Code:   "JPY",
		Name:   "Japanese Yen",
		Symbol: "¥",
		Type:   "fiat",
	},
	"GBP": {
		Code:   "GBP",
		Name:   "British Pound Sterling",
		Symbol: "£",
		Type:   "fiat",
	},

	"BTC": {
		Code:   "BTC",
		Name:   "Bitcoin",
		Symbol: "₿",
		Type:   "crypto",
	},
	"ETH": {
		Code:   "ETH",
		Name:   "Ethereum",
		Symbol: "Ξ",
		Type:   "crypto",
	},
	"SOL": {
		Code:   "SOL",
		Name:   "Solana",
		Symbol: "◎",
		Type:   "crypto",
	},
}
