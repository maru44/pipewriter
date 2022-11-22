.PHONY: test

test:
	@cd e2e && \
	go run main.go && \
	cd .. && \
	go test ./...
