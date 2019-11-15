FROM golang:1.13 AS stage
ENV CGO_ENABLED 0
WORKDIR /rest-go-pg
COPY . .
RUN go build -mod=vendor -v -o ./application cmd/*.go

FROM alpine:3.7
WORKDIR /rest-go-pg
COPY --from=stage /rest-go-pg/application application
ENTRYPOINT [ "/rest-go-pg/application" ]