BINARY=engine
test: 
	go test -v -cover -covermode=atomic ./...

engine:
	go build -o ${BINARY} api/*.go


unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker:
	docker build -t poc-gofiber-clean-arch .

run:
	docker-compose -f docker-compose.yaml -f docker-compose.postgres.yaml up --build -d

run-mysql:
	docker-compose -f docker-compose.yaml -f docker-compose.mysql.yaml up --build -d

run-postgres:
	docker-compose -f docker-compose.yaml -f docker-compose.postgres.yaml up --build -d

stop:
	docker-compose down --remove-orphans

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: clean install unittest build docker run stop vendor lint-prepare lint