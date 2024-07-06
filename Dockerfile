FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .. .

RUN go build -o /voting_system

EXPOSE 50051

CMD ["/voting_system"]
