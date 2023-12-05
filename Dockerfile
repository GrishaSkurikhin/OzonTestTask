FROM golang:1.21

RUN mkdir -p /app
ADD ./ /app/
WORKDIR /app

RUN go mod download && go mod verify
RUN go build -o bin ./cmd/app

EXPOSE 8081

CMD ["/app/bin"]