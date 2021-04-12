
all: clean

clean: 
	find . -name 'trendServer' -type f -exec rm -f {} \; 
	go clean

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt trendServer

install:
	go install trendServer

build:
	go build trendServer

start:
	./trendServer 

test-ping:
	curl localhost:3000/ping

test-gumball:
	curl localhost:3000/gumball

docker-build: 
	docker build -t vishalgadapa86/trendserver .
	docker images

docker-run:
	docker run --name goapi -td -p 3005:3005 vishalgadapa86/trendserver
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

