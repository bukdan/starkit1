FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://proxy.golang.org,direct
RUN go build -o /gateway ./ 

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /gateway /gateway
EXPOSE 8080
ENTRYPOINT ["/gateway"]
