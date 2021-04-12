
all: clean

clean: 
	find . -name 'linkRedirect' -type f -exec rm -f {} \; 
	go clean

run:
	go run src/app/$(app).go

main:
	go run src/app/main.go

format:
	go fmt linkRedirect

install:
	go install linkRedirect

build:
	go build linkRedirect

start:
	./goapi 

test-ping:
	curl localhost:3004/ping

test-gumball:
	curl localhost:3004/gumball

docker-build: 
	docker build -t vishalgadapa86/linkredirect .
	docker images

docker-run:
	docker run --name linkredirect -td -p 3004:3004 vishalgadapa86/linkredirect
	docker ps

docker-run-test:
	docker run --restart always --name linkredirect -td -p 3004:3004 vishalgadapa86/linkredirect
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

