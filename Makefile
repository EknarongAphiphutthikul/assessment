test-unit:
	go clean -testcache && go test -v --tags=unit ./...

test-integrate:
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests

pq-run:
	docker compose -f docker-compose-postgres.yml up --detach

pq-stop:
	docker compose  -f docker-compose-postgres.yml stop

run: pq-run
	DATABASE_URL=postgres://root:root@localhost:5432/go-example-db?sslmode=disable PORT=2565  go run server.go

docker-build:
	docker build -t kbtg/kampus/go/assessment:latest .

docker-run: pq-run
	docker compose -f docker-compose.yml up --detach

docker-stop:
	docker compose -f docker-compose.yml stop

docker-down:
	docker compose -f docker-compose.yml down