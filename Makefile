.PHONY: compile up down logs

compile:
	docker run --rm -w "/go/src/service" \
	-v "${PWD}/service:/go/src/service" \
	-e "GOPATH=/go/src/service/Godeps/_workspace:/go" \
	-e "CGO_ENABLED=0" \
	golang:alpine \
	go build -v -o service service.go

up:
	docker-compose down
	docker-compose up --build -d

down:
	docker-compose down

logs:
	docker-compose logs -f