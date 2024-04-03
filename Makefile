.PHONY: build

build-amd:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go;\
	sudo DOCKERFILE=Dockerfile.amd TAG=amd docker-compose build;\
	sudo DOCKERFILE=Dockerfile.amd TAG=amd docker-compose up -d

build-arm:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o main main.go;\
	sudo DOCKERFILE=Dockerfile.arm TAG=arm docker-compose build;\
	sudo DOCKERFILE=Dockerfile.arm TAG=arm docker-compose up -d

# deploy:
# 	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go;\
# 	sudo docker-compose build;\
# 	docker tag chromedp1-visa:latest nuodi/chromedp-visa:latest ;\
# 	docker push nuodi/chromedp-visa:latest

