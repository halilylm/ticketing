.PHONY: sync_common run_docker run_auth run_orders run_payments run_tickets git_push
URL="github.com/halilylm/gommon@v1.1.22"
sync_common:
	cd orders && go get ${URL}
	cd tickets && go get ${URL}
	cd payments && go get ${URL}

run_docker:
	docker run -d -p 27017:27017 --name ticketing mongo
	docker run -d -p 4222:4222 -p 8222:8222 nats-streaming

run_auth:
	cd auth && go mod tidy && go run app/main.go

run_orders:
	cd orders && go mod tidy && go run app/main.go

run_payments:
	cd payments && go mod tidy && go run app/main.go

run_tickets:
	cd tickets && go mod tidy && go run app/main.go

git_push:
	cd auth && go mod tidy
	cd orders && go mod tidy
	cd payments && go mod tidy
	cd tickets && go mod tidy
	git add .
	git commit -m "just update"
	git push -f git@github.com:halilylm/ticketing.git main