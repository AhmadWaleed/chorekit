FROM golang:1.16-alpine AS mod
RUN apk add -U git
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM golang:1.16-alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=mod $GOCACHE $GOCACHE
COPY --from=mod $GOPATH/pkg/mod $GOPATH/pkg/mod
COPY --from=mod /src/main .
COPY --from=mod /src/.env . 

EXPOSE 8080

CMD ["./main"]