.PHONY: build run docker-build docker-run

build:
	go build -o bin/glasnik ./cmd/cli

run: build
	./bin/glasnik

docker-build:
	docker-compose build

docker-run: docker-build
	docker-compose up

docker-stop:
	docker-compose down

# Build individual services
docker-build-proxy:
	docker build -t glasnik-proxy -f build/proxy/Dockerfile .

docker-build-registry:
	docker build -t glasnik-registry -f build/registry/Dockerfile .

# Run individual services
docker-run-registry:
	docker run -p 8082:8082 glasnik-registry

docker-run-proxy:
	docker run -p 8083:8083 \
		-e PROXY_USERNAME=myuser \
		-e PROXY_PASSWORD=mypass \
		glasnik-proxy \
		--host 0.0.0.0 --port 8083 \
		--username myuser \
		--password mypass \
		--registry http://0.0.0.0:8082