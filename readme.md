# Monobank to Taxer.ua syncer

`go-monobank-taxer` is a simple tool to sync transactions from Monobank (Ukrainian bank) to Taxer (a service for
tracking personal finances and taxes in Ukraine). It supports transactions made by individual entrepreneurs (FOP).

* Monobank Docs: https://api.monobank.ua/docs/

## Usage

Before running the tool, you need to set the following environment variables:

* `TAXER_EMAIL`: Your Taxer email
* `TAXER_PASS`: Your Taxer password
* `MONO_TOKEN`: Your Monobank token (get your Monobank token from https://api.monobank.ua)
* `FROM`: Date (YYYY-MM-DD) to start syncing transactions from
* `TO`: Date (YYYY-MM-DD) to sync transactions up to

Alternatively, you can provide these values as command-line flags,
e.g., `--taxer-email`, `--taxer-pass`, `--mono-token`, `--from`, and `--to`.

To run the tool, execute:

```sh
go run main.go
```

## How it works

The tool will perform the following steps:

* Log in to your Taxer account
* Fetch your Taxer and Monobank accounts
* Sync any unsynced Monobank accounts with Taxer
* Fetch transactions from Monobank for the specified date range
* Sync the fetched transactions with your Taxer account

Note that the tool will only sync transactions with a positive amount.

## Warning

Monobank API has limitations that do not allows to get transactions for period more than 31 days + 1 hour.
Also, getting transactions is not possible more frequently than once per 1 minute. It makes this process a bit slow.

## Remarks

* Taxer.ua has no public API, but I got an acceptance to use the private API that uses theirs UI.