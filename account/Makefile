PROTO_ROOT=proto
DOCKER_IMAGE=proto-builder

.PHONY: docker-build gen clean

docker-build:
	docker build -t $(DOCKER_IMAGE) .

gen: docker-build
	docker run --rm -v $(abspath .):/app $(DOCKER_IMAGE) \
		bash -c '\
		set -e; \
		for dir in pagination account; do \
			echo ">> Processing $$dir"; \
			mkdir -p /app/pkg/$$dir/go; \
			cd /app/$(PROTO_ROOT)/$$dir; \
			for file in *.proto; do \
				echo "Generating $$dir/$$file"; \
				protoc \
					-I . \
					-I /app/$(PROTO_ROOT) \
					-I /usr/local/include/googleapis \
					--go_out=/app/pkg/$$dir/go \
					--go_opt=paths=source_relative \
					--go-grpc_out=/app/pkg/$$dir/go \
					--go-grpc_opt=paths=source_relative \
					$$file; \
			done; \
		done'

clean:
	rm -rf pkg/account/go pkg/pagination/go
