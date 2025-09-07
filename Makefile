build:
test:
docker-up:
docker-down:
APP=bbb-voting

.PHONY: build run test docker-up docker-down clean install-cobra

install-cobra:
	go get github.com/spf13/cobra@latest

build: install-cobra
	go mod tidy
	go build -o $(APP) main.go

run: build
	./$(APP)

test:
	go test ./...

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

clean:
	rm -f $(APP)

setup:
	go install github.com/golang/mock/mockgen@latest
	go get github.com/stretchr/testify/mock
