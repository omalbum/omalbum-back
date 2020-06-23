FROM golang:1.14.1

WORKDIR /go/src/github.com/miguelsotocarlos/teleoma

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

WORKDIR cmd/api
RUN go build -v .

EXPOSE 8080

CMD ["./api"]
