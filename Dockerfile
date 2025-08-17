
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct
RUN go mod download
COPY . .
RUN go build -o /user-service ./

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /user-service /user-service
EXPOSE 8081
ENTRYPOINT ["/user-service"]
