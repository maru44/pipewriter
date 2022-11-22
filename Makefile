.PHONY: test e2e

e2e:
	@cd e2e && \
	go run main.go && \
	cd ..

test:
	@go version && \
	go test ./...
