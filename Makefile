#!make
POSTGRESQL_URL := postgres://postgres:postgres@localhost:5432/master?sslmode=disable

tools:
	cd tools && go mod tidy && go mod vendor && go mod verify && go generate -tags tools
.PHONY: tools

vendor:
	go mod tidy && go mod vendor && go mod verify
.PHONY: vendor

gen_contracts:
	solc0611 --abi ./contracts/DepositContract/DepositContract.sol -o ./contracts/DepositContract/abi --overwrite && \
 	solc0611 --bin ./contracts/DepositContract/DepositContract.sol -o ./contracts/DepositContract/bin --overwrite && \
 		./bin/abigen --bin=./contracts/DepositContract/bin/DepositContract.bin \
     	             --abi=./contracts/DepositContract/abi/DepositContract.abi \
     	             --pkg=deposit --out=internal/generated/contracts/deposit.go
.PHONY: gen_contracts

abigen:

.PHONY: abigen

### Migrations
.PHONY: rollback
rollback:
	bin/migrate -database ${POSTGRESQL_URL} -path db/migrations down

.PHONY: migrate
migrate:
	bin/migrate -database ${POSTGRESQL_URL} -path db/migrations up

.PHONY: up
up:
	UID_GID="$(id -u):$(id -g)" docker-compose -f docker-compose.yml up -d

.PHONY: up-with-data
up-with-data:
	docker-compose -f docker-compose-pg-with-data.yml up -d

.PHONY: down
down:
	UID_GID="$(id -u):$(id -g)" docker-compose -f docker-compose.yml down

.PHONY: run
run:
	make up && make migrate