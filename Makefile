dep:
	@echo ">> Downloading Dependencies"
	@go mod download

build: dep
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o ./bin ./...

run-server: dep
	env $$(cat .env | xargs) go run tempo/cmd server

migrate:
	eval $$(egrep -v '^#' .env | xargs -0) go run tempo/cmd migrate

test-all: test-unit test-integration-with-infra

test-unit: dep
	@echo ">> Running Unit Test"
	@env $$(cat .env.testing | xargs) go test -tags=unit -failfast -cover -covermode=atomic ./...

test-integration: dep
	@echo ">> Running Integration Test"
	@env $$(cat .env.testing | xargs) env DB_MIGRATION_PATH=$$(pwd)/migrations go test -tags=integration -failfast -cover -covermode=atomic ./...

test-integration-with-infra: test-infra-up test-integration test-infra-down

test-infra-up:
	$(MAKE) test-infra-down
	@echo ">> Starting Test DB"
	docker run -d --rm --name test-mysql -p 3343:3306 --env-file .env.testing mysql:5.7  --sql-mode="STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION"
	docker cp $$(pwd)/deployments/docker test-mysql:/tools
	docker exec test-mysql sh -c '/tools/wait-for-mysql.sh 40'
	docker ps

test-infra-down:
	@echo ">> Shutting Down Test DB"
	@-docker kill test-mysql
