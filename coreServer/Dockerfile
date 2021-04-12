FROM golang:latest 
EXPOSE 3003
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
RUN cd /app ; go install coreServer
CMD ["/app/bin/coreServer"]
