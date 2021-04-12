
all: clean

clean: 
	find . -name 'controlPanel' -type f -exec rm -f {} \; 
	go clean

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt controlPanel

install:
	go install controlPanel

build:
	go build controlPanel

start:
	./controlPanel 

test-ping:
	curl localhost:3001/ping

test-gumball:
	curl localhost:3001/gumball

docker-build: 
	docker build -t vishalgadapa86/controlpanel .
	docker images

docker-run:
	docker run --restart always --name controlpanel -td -p 3001:3001 vishalgadapa86/controlpanel
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

