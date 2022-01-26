
DOCKER_REDIS_NAME=redis-dev
REDIS_PORT=6379

clean:
	docker rm -f $(DOCKER_REDIS_NAME)
	docker-compose -f tools/docker-compose.yml down

dev:
	docker start $(DOCKER_REDIS_NAME) || docker run -d -p $(REDIS_PORT):$(REDIS_PORT) --name $(DOCKER_REDIS_NAME) redis:6.2-alpine 
	go run cmd/cjtim-backend-go/main.go

dev-all:
	docker-compose -f tools/docker-compose.yml build
	docker-compose -f tools/docker-compose.yml up -d --remove-orphans

# make build tag=tests
build:
	docker build -t $(tag) -f tools/Dockerfile .