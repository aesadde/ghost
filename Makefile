get:
	go get ./...
	go get -u github.com/swaggo/swag/cmd/swag

install:
	go install ./...

test: 
	go test ./... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt

coverage: test
	go tool cover -html=.coverage.txt
