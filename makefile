# baseLibrary makefile
.PHONY: proto test vendor

main: clean gen proto
clean:
	@ find . -name '*pb.go' -delete
	@ find . -name '*_generated.go' -delete
generate:
	@ go generate ./...


# test
test:
	@ go test ./...
test-race:
	@ go test -race ./...


# proto
proto:
	@ make proto-clean
	@ make proto-generate
proto-clean:
	@ find ./proto -name '*_generated.go' -delete
proto-generate:
	@ spec generate	./proto/pclock
