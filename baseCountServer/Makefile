
all: clean

clean: 
	find . -name 'baseCountServer' -type f -exec rm -f {} \; 
	go clean

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt baseCountServer

install:
	go install baseCountServer

build:
	go build baseCountServer

start:
	./baseCountServer 

test-ping:
	curl localhost:3000/ping

test-gumball:
	curl localhost:3000/gumball

docker-build: 
	docker build -t vishalgadapa86/basecountserver .
	docker images

docker-run:
	docker run --name baseCountServer -td -p 3002:3002 basecountserver
	docker ps

docker-run-test:
	docker run --name baseCountServer -td -p 3002:3002 basecountserver
	docker ps

docker-network:
	docker network ls

docker-network-inspect:
	docker network inspect host

docker-shell:
	docker exec -it baseCountServer bash 

docker-clean:
	docker stop baseCountServer
	docker rm baseCountServer
	docker rmi baseCountServer

docker-ip:
	docker-machine ip

docker-ps:
	 docker ps --all --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t"

docker-ps-ports:
	 docker ps --all --format "table {{.Names}}\t{{.Ports}}\t"

