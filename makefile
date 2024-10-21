create-db:
	sqlite3 data.db < ./data/init.sql

build-app:
	go build -o bin/app .

run-app: build-app
	./bin/app
	