# Transaction service listing

This folder contains a restored, compilable version of the Transaction project from chapter 4.

It includes:

- contracts for deposits, withdrawals, transfers, and transaction history;
- PostgreSQL migrations for `transactions` and `transaction_entries`;
- repository operations wrapped in database transactions;
- a service layer for transaction use cases;
- an adapter that calls Account service to update balances;
- a gRPC server and application wiring.

The OCR text in the book contains broken characters, so the code here keeps the intent of the listings while fixing Go syntax, import paths, struct tags, SQL names, and money handling.
