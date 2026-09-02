# Gateway service listing

This folder contains a restored, compilable version of the Gateway project from chapter 3.

It includes:

- the Gateway gRPC contract in `proto/gateway/gateway.proto`;
- local contract packages for `gateway`, `account`, `auth`, and `pagination`;
- a JWT interceptor that protects non-public Gateway methods;
- isolated adapters for calling Auth and Account services;
- Gateway business logic that coordinates registration, login, token refresh, validation, and user operations;
- a server layer that reads the current user id from request context for `me` endpoints.

The OCR text in the book contains broken characters, so the code here keeps the intent of the listings while fixing Go syntax, import paths, struct tags, field names, and JWT handling.
