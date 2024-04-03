.PHONY: build

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go;\
	sudo docker-compose build;\
	sudo docker-compose up -d
deploy:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go;\
	sudo docker-compose build;\
	docker tag chromedp1-visa:latest nuodi/chromedp-visa:latest ;\
	docker push nuodi/chromedp-visa:latest