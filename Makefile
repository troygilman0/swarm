test:
	go test ./... -count=1 --race -v

dev:
	go run ./examples/simple/.

serv:
	go run ./examples/server/.
