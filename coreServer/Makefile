
all: clean

clean: 
	find . -name 'coreServer' -type f -exec rm -f {} \; 
	go clean

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt coreServer

install:
	go install coreServer

build:
	go build coreServer

start:
	./coreServer 

test-ping:
	curl localhost:3003/ping

test-gumball:
	curl localhost:3003/gumball

docker-build: 
	docker build -t vishalgadapa86/coreserver .
	docker images

docker-run:
	docker run --name --restart=always coreserver -td -p 3003:3003 vishalgadapa86/coreserver
	docker ps

docker-network:
	docker network ls

docker-network-inspect:
	docker network inspect host

docker-shell:
	docker exec -it goapi bash 

docker-clean:
	docker stop goapi
	docker rm goapi
	docker rmi goapi

docker-ip:
	docker-machine ip

docker-ps:
	 docker ps --all --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t"

docker-ps-ports:
	 docker ps --all --format "table {{.Names}}\t{{.Ports}}\t"

