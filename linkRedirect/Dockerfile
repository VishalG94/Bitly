FROM golang:latest 
EXPOSE 3004
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
RUN cd /app ; go install linkRedirect
CMD ["/app/bin/linkRedirect"]
