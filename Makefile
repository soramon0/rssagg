build:
	go build .

run:
	go build .
	./rssagg

sqlc:
	sqlc generate

migrate-up:
	goose -dir ./sql/schema/ postgres postgres://postgres:example@localhost:5432/dev_db up 

migrate-down:
	goose -dir ./sql/schema/ postgres postgres://postgres:example@localhost:5432/dev_db down 
