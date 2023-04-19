package main

import (
	"github.com/leonidpodriz/go-monobank-taxer/taxer"
	"github.com/vtopc/go-monobank"
	"time"
)

type TimePeriod struct {
	From time.Time
	To   time.Time
}

func ReversedTimePeriods(from, to time.Time, step time.Duration) []TimePeriod {
	var periods []TimePeriod
	var cursor = to

	if to.Before(from) {
		return periods
	}

	for {
		period := TimePeriod{
			From: cursor.Add(-step),
			To:   cursor,
		}

		if period.From.Before(from) {
			period.From = from
		}

		periods = append(periods, period)

		if period.From.Equal(from) {
			break
		}

		cursor = period.From
	}

	return periods
}

func CorrespondingTaxerAccount(taxAccounts []taxer.Account, monoAcc monobank.Account) *taxer.Account {
	for _, taxAcc := range taxAccounts {
		if taxAcc.Comment != "auto-synced" {
			continue
		}

		if taxAcc.Num != monoAcc.IBAN {
			continue
		}

		monoCurrency := currencies[monoAcc.CurrencyCode]

		if taxAcc.Currency != monoCurrency {
			continue
		}

		return &taxAcc
	}

	return nil
}
