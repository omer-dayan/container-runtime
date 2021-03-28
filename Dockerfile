FROM golang:1.16

RUN mkdir /opt/runai-container-runtime
WORKDIR /opt/runai-container-runtime
ADD . / ./
RUN go get -d ./...
RUN go test ./...
RUN make