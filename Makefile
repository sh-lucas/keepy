run:
	docker build -f server.Dockerfile -t keepy-server .
	docker compose up

deploy:
	docker build -f server.Dockerfile -t keepy-server:1.0.3 .
	docker save keepy-server:1.0.3 | ssh oracle "docker load"

hub:
	docker build -f keepy.Dockerfile -t catnipbrewer/keepy-server:latest -t catnipbrewer/keepy:1.0.3 .
	docker push catnipbrewer/keepy:1.0.3
	docker push catnipbrewer/keepy:latest