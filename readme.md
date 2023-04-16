# Monobank to Taxer.ua synchronization tool

This tool allows you to synchronize your Monobank transactions of entrepreneur accounts with Taxer.ua.

## How it works?

1. Getting personal information from Monobank using [Monobank API](https://api.monobank.ua/docs/) and provided token in
   order to get list of entrepreneur accounts.
2. Getting transactions from Monobank using [Monobank API](https://api.monobank.ua/docs/) and provided token in order to
   get list of transactions for each entrepreneur account. The API has has next limitations:
    * The maximum period for which you can get transactions is 2682000 seconds (31 days + 1 hour) per one request. So,
      if you need to get transactions for more than 31 days, the tool need's to make several requests.
    * Only 1 request per minute is allowed. So if tool need's to make more than 1 request to get transactions for more
      than 31 days, the tool need's to wait 1 minute between requests.
3. Check is Monobank account is already synchronized with Taxer.ua account. If not, then synchronize it (create new
   Taxer.ua account and link it with Monobank account using comment on Taxer.ua account).
4. Getting operations (it's like Monobank transactions) from Taxer.ua accounts.
5. Process transactions from Monobank and operations from Taxer.ua. If transaction is not found in operations, then
   prepare create operation in Taxer.ua.
6. Use prepared operations to create operations in Taxer.ua in batch mode (actually it is just one request to Taxer.ua).

## Remarks

* Taxer.ua has no public API, so the tool uses reverse engineering to get operations from Taxer.ua. I got an acceptance
  from Taxer.ua to use this approach.
