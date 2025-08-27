run:
	docker build -f server.Dockerfile -t keepy-server .
	docker compose up

deploy:
	docker build -f server.Dockerfile -t keepy-server:1.0.0 .
	docker save keepy-server:1.0.0 | ssh oracle "docker load"

hub:
	docker build -f keepy.Dockerfile -t catnipbrewer/keepy-server:latest -t catnipbrewer/keepy-server:1.0.2 .
	docker push catnipbrewer/keepy-server:1.0.2
	docker push catnipbrewer/keepy-server:latest