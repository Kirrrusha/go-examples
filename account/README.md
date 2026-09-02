# User service listing

This folder contains a restored, compilable version of the first microservice from chapter 1.

The chapter calls the module `User`, while the code uses the `Account` service and `account` proto package. This project keeps that convention:

- `proto/account/account.proto` contains the user/account gRPC contract;
- `proto/pagination/pagination.proto` contains reusable pagination;
- `internal/model` contains business-layer user models;
- `internal/mapper` maps transport models to business models;
- `internal/repository` contains GORM storage logic and mappers;
- `internal/service` contains `AccountService`;
- `internal/server` contains the gRPC transport layer;
- `internal/app` wires config, migrations, repository, service, and transport together;
- `cmd/account` is the final lightweight entry point from listing 1.29;
- `migrations` creates the `users` table.

The OCR text in the book contains broken characters, so the code here keeps the intent of the listings while fixing Go syntax, import paths, struct tags, table names, and SQL.

The real book flow uses generated protobuf code from a separate `contracts` repository. Because this folder is self-contained, `pkg/account/go` and `pkg/pagination/go` contain lightweight local stubs matching the contract shape. They can be replaced by generated files after running `make gen`.
