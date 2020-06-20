FROM golang:1.14.1

WORKDIR /go/src/github.com/miguelsotocarlos/teleoma
COPY . .

RUN go mod download

WORKDIR cmd/api
RUN go build -v .

EXPOSE 8080

CMD ["./api"]
