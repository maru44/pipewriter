.PHONY: test testfile

testfile:
	@cd e2e && \
	go run main.go && \
	cd ..

test:
	go test ./...
