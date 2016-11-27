.PHONY: compile publish up down logs

compile:
	docker run --rm -w "/go/src/${name}" \
		-v "${PWD}/${name}:/go/src/${name}" \
		-e "GOPATH=/go/src/${name}/Godeps/_workspace:/go" \
		-e "CGO_ENABLED=0" \
		golang:alpine \
		go build -v -o ${name} ${name}.go

publish: compile
	docker-compose build

helm: publish
	helm install --debug ${name}

run: compile
	docker-compose down
	docker-compose up --build -d

clean:
	docker-compose down

logs:
	docker-compose logs -f