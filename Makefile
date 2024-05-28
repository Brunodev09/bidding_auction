run: 
	docker compose up 

rebuild-run:
	docker compose up --build

test:
	go test -v ./...