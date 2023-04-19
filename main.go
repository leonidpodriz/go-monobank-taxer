package main

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/leonidpodriz/go-monobank-taxer/taxer"
	"github.com/vtopc/go-monobank"
	"log"
	"os"
	"time"
)

const MonoStep = 2682000 * time.Second

type Opts struct {
	TaxerEmail string `long:"taxer-email" env:"TAXER_EMAIL" required:"true" description:"email of taxer"`
	TaxerPass  string `long:"taxer-pass" env:"TAXER_PASS" required:"true" description:"password of taxer"`
	MonoToken  string `long:"mono-token" env:"MONO_TOKEN" required:"true" description:"token of monobank"`
	From       string `long:"from" env:"FROM" required:"true" description:"from date"`
	To         string `long:"to" env:"TO" required:"true" description:"to date"`
}

type TimePeriod struct {
	From time.Time
	To   time.Time
}

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	log.Println("[INFO] Starting syncer")

	from, err := time.Parse("2006-01-02", opts.From)

	if err != nil {
		log.Fatalf("[PANIC] Error parsing from date: %s", err)
	}

	to, err := time.Parse("2006-01-02", opts.To)

	if from.After(to) {
		log.Fatalf("[PANIC] From date is after to date")
	}

	if err != nil {
		log.Fatalf("[PANIC] Error parsing to date: %s", err)
	}

	log.Println("[INFO] Creating taxer and mono clients...")
	tax, err := taxer.NewClient(opts.TaxerEmail, opts.TaxerPass)

	if err != nil {
		log.Fatalf("Error creating taxer client: %s", err)
	}
	mono := monobank.NewPersonalClient(nil).WithAuth(monobank.NewPersonalAuthorizer(opts.MonoToken))

	log.Println("[INFO] Logging in to taxer...")
	err = tax.Login()

	if err != nil {
		log.Fatalf("[INFO] Error logging in to taxer: %s", err)
	}

	log.Println("[INFO] Getting taxer user account...")
	userAccount, err := tax.GetUserAccount()

	if err != nil {
		log.Printf("[PANIC] Error getting user account: %s", err)
		return
	}

	monoInfo, err := mono.ClientInfo(ctx)

	if err != nil {
		log.Fatalf("[INFO] Error getting mono info: %s", err)
	}

	log.Printf("[INFO] Mono user: %s", monoInfo.Name)

	taxUser := userAccount.GetFirstUser()
	log.Printf("[INFO] Taxer user: %s", taxUser.Name)

	taxAccounts, err := tax.GetAllAccounts(taxUser.Id)

	if err != nil {
		log.Printf("[PANIC] Error getting taxer accounts: %s", err)
		return
	}

	var unSyncedAccounts []monobank.Account

	for _, monoAcc := range monoInfo.Accounts {
		if monoAcc.Type != "fop" {
			continue
		}

		taxAcc := CorrespondingTaxerAccount(taxAccounts, monoAcc)

		if taxAcc == nil {
			unSyncedAccounts = append(unSyncedAccounts, monoAcc)
		}
	}

	log.Printf("[INFO] Unsynced accounts: %d", len(unSyncedAccounts))

	for _, monoAcc := range unSyncedAccounts {
		err = tax.CreateAccount(taxUser.Id, taxer.Account{
			Title:    fmt.Sprintf("%s: %s", monoInfo.Name, monoAcc.AccountID),
			Num:      monoAcc.IBAN,
			Currency: CurrencyCode(monoAcc.CurrencyCode),
			Comment:  "auto-synced",
			Bank:     "monobank",
		})

		if err != nil {
			log.Printf("[ERROR] Error creating account: %s", err)
		}
	}

	taxAccounts, err = tax.GetAllAccounts(taxUser.Id)

	if err != nil {
		log.Fatalf("[ERROR] Error getting taxer accounts after sync: %s", err)
	}

	log.Println("[INFO] Syncing transactions...")

	periods := ReversedTimePeriods(from, to, MonoStep)
	var operations []taxer.UncreatedOperation

	for accIdx, monoAcc := range monoInfo.Accounts {
		var monoTransactions []monobank.Transaction
		if monoAcc.Type != "fop" {
			continue
		}

		taxAcc := CorrespondingTaxerAccount(taxAccounts, monoAcc)
		taxOps, err := tax.GetAllOperationsForPeriod(taxUser.Id, taxAcc.Id, from, to)

		if err != nil {
			log.Printf("[ERROR] Error getting taxer operations for %s: %s", monoAcc.IBAN, err)
			continue
		}

		if taxAcc == nil {
			log.Printf("[ERROR] no taxer account (after sync!) for %s", monoAcc.IBAN)
		}

		for pIdx, period := range periods {
			log.Printf("[INFO] Getting transactions for %s (%d/%d) from %s to %s (%d/%d)", monoAcc.IBAN, accIdx+1, len(monoInfo.Accounts), period.From, period.To, pIdx+1, len(periods))

			monoTrx, err := mono.Transactions(ctx, monoAcc.AccountID, period.From, period.To)

			if err != nil {
				log.Printf("[ERROR] Error getting operations for %s: %s", monoAcc.IBAN, err)
				continue
			}

			for _, t := range monoTrx {
				if t.Amount < 0 {
					continue
				}

				monoTransactions = append(monoTransactions, t)
			}

			if pIdx != len(periods)-1 || accIdx != len(monoInfo.Accounts)-1 {
				time.Sleep(time.Minute)
			}
		}

		for _, monoOp := range monoTransactions {
			isSynced := false
			for _, taxOp := range taxOps {
				if taxOp.Comment != monoOp.Description {
					continue
				}

				log.Printf("%f %f", taxOp.Contents[0].SumCurrency, float64(monoOp.Amount)/100)
				if taxOp.Contents[0].SumCurrency != float64(monoOp.Amount)/100 {
					continue
				}

				isSynced = true

			}

			if !isSynced {
				operations = append(operations, taxer.UncreatedOperation{
					Account: taxer.OperationAccount{
						Id:       taxAcc.Id,
						Title:    taxAcc.Title,
						Currency: taxAcc.Currency,
					},
					Total:     float64(monoOp.Amount) / 100,
					Comment:   monoOp.Description,
					Timestamp: int(monoOp.Time.Unix()),
				})
			}
		}

	}

	log.Printf("[INFO] Syncing %d operations...", len(operations))

	if len(operations) == 0 {
		log.Println("[INFO] Nothing to sync")
		return
	}
	err = tax.CreateOperations(taxUser.Id, operations)

	if err != nil {
		log.Printf("[ERROR] Error creating operations: %s", err)
	}
}
