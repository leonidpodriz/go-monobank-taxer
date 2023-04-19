package main

import "testing"

func TestCurrencyCode(t *testing.T) {
	var tests = []struct {
		currency int
		code     string
	}{
		{980, "UAH"},
		{840, "USD"},
		{978, "EUR"},
	}

	for _, test := range tests {
		if got := CurrencyCode(test.currency); got != test.code {
			t.Errorf("CurrencyCode(%d) = %s, want %s", test.currency, got, test.code)
		}
	}
}
