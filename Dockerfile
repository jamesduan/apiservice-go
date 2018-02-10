FROM golang:latest

RUN mkdir -p /go/src/apiservice
COPY . /go/src/apiservice
WORKDIR /go/src/apiservice
# CMD ["go", "build"]
# RUN ./apiservice
RUN go build . && ./apiservice