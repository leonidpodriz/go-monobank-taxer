package main

import (
	"github.com/leonidpodriz/go-monobank-taxer/taxer"
	"github.com/vtopc/go-monobank"
	"log"
	"testing"
	"time"
)

func TestReversedTimePeriods(t *testing.T) {
	from := time.Date(2022, 2, 23, 11, 0, 0, 0, time.UTC)
	to := time.Date(2022, 2, 24, 4, 0, 0, 0, time.UTC)
	step := 4 * time.Hour
	expected := []TimePeriod{
		{
			To:   time.Date(2022, 2, 24, 4, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			To:   time.Date(2022, 2, 24, 0, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 23, 20, 0, 0, 0, time.UTC),
		},
		{
			To:   time.Date(2022, 2, 23, 20, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 23, 16, 0, 0, 0, time.UTC),
		},
		{
			To:   time.Date(2022, 2, 23, 16, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 23, 12, 0, 0, 0, time.UTC),
		},
		{
			To:   time.Date(2022, 2, 23, 12, 0, 0, 0, time.UTC),
			From: time.Date(2022, 2, 23, 11, 0, 0, 0, time.UTC),
		},
	}
	periods := ReversedTimePeriods(from, to, step)
	assertReversedTimePeriods(t, expected, periods)

	// what to is after from
	from = time.Date(2022, 2, 24, 4, 0, 0, 0, time.UTC)
	to = time.Date(2022, 2, 23, 11, 0, 0, 0, time.UTC)
	step = 4 * time.Hour
	expected = []TimePeriod{}
	periods = ReversedTimePeriods(from, to, step)
	assertReversedTimePeriods(t, expected, periods)
}

func TestCorrespondingTaxerAccount(t *testing.T) {
	monoAcc := monobank.Account{
		AccountID:    "123",
		SendID:       "123",
		Balance:      0,
		CreditLimit:  0,
		CurrencyCode: 980,
		CashbackType: "",
		CardMasks:    []string{},
		Type:         "",
		IBAN:         "GB31THNQ18572134910363",
	}
	taxAccounts := []taxer.Account{
		{
			Id:       1,
			Num:      "GB31THNQ18572134910363",
			Currency: "UAH",
			Comment:  "auto-synced",
		},
		{
			Id:       2,
			Num:      "GB42VNPB86361664761395",
			Currency: "UAH",
			Comment:  "auto-synced",
		},
		{
			Id:       3,
			Num:      "GB44EQXF48981416571049",
			Currency: "UAH",
			Comment:  "auto-synced",
		},
		{
			Id:       4,
			Num:      "GB09ALPL79972788787607",
			Currency: "UAH",
			Comment:  "auto-synced",
		},
		{
			Id:       5,
			Num:      "GB85XJMY38442198542815",
			Currency: "UAH",
			Comment:  "auto-synced",
		},
	}
	acc := CorrespondingTaxerAccount(taxAccounts, monoAcc)
	if acc == nil {
		t.Fatalf("expected account, got nil")
	}

	if acc.Num != monoAcc.IBAN {
		t.Fatalf("expected %s, got %s", monoAcc.IBAN, acc.Num)
	}

	if acc.Id != 1 {
		t.Fatalf("expected %d, got %d", 1, acc.Id)
	}
}

func assertReversedTimePeriods(t *testing.T, expected, actual []TimePeriod) {
	for _, period := range actual {
		log.Printf("%v", period)
	}
	if len(expected) != len(actual) {
		t.Fatalf("expected %d periods, got %d", len(expected), len(actual))
	}

	for i, period := range expected {
		if period.From != actual[i].From || period.To != actual[i].To {
			t.Fatalf("expected %v, got %v", period, actual[i])
		}
	}
}
