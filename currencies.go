package main

var currencies = map[int]string{
	980: "UAH",
	840: "USD",
	978: "EUR",
}

func CurrencyCode(currency int) string {
	return currencies[currency]
}
