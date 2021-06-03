FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go env -w GO111MODULE=auto
RUN go build -o main .

CMD ["main  "]