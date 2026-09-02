# Go microservice examples

Restored and compilable Go microservice examples based on the book listings.

The original OCR text contained broken characters and incomplete snippets, so these examples keep the intent of the listings while fixing Go syntax, import paths, struct tags, table names, SQL, and service wiring.

## Examples

- `account` - user/account service with PostgreSQL storage, migrations, pagination, gRPC server, and local protobuf stubs.
- `auth` - authentication service with registration, login, refresh tokens, JWT validation, logout, PostgreSQL storage, and migrations.
- `gateway` - gateway service that coordinates Auth and Account services and protects private methods with a JWT interceptor.
- `transaction` - transaction service with deposits, withdrawals, transfers, history, PostgreSQL transactions, and Account-service balance updates.

Each service is self-contained and has its own `README.md`, `go.mod`, `Makefile`, `Dockerfile`, proto files, and application entry point.
