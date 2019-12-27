FROM golang:latest
LABEL maintaner="Valiev Ismail"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -mod=vendor -v -o ./rest-go-pg cmd/*.go
EXPOSE 3000
CMD ["./rest-go-pg"]