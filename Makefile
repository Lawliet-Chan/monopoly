PROJECT=monopoly

build_chain:
	go build -v -o $(PROJECT) ./cmd/web3/main.go

run_chain:
	go run ./cmd/web3/main.go

reset:
	rm -rf yu