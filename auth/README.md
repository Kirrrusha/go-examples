# Auth service listing

This folder contains a restored, compilable version of the authorization/authentication project from chapter 2.

It includes:

- the Auth gRPC contract in `proto/auth/auth.proto`;
- a `Makefile` and `Dockerfile` for generating Go code from the contract into `pkg/auth/go`;
- config fields for JWT secret and token lifetimes;
- storage models for users and refresh tokens;
- database migration for `users` and `refresh_tokens`;
- repository methods for users and refresh tokens;
- service methods for register, login, refresh, validate, and logout.

The OCR text in the book contains broken characters, so the code here keeps the intent of the listings while fixing Go syntax, import paths, struct tags, table names, and JWT/bcrypt calls.
