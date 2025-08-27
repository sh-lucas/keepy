run:
	docker build -f server.Dockerfile -t keepy-server .
	docker compose up