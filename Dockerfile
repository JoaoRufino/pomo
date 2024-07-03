FROM golang:1.16
WORKDIR /go/src/github.com/joaorufino/pomo
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/pomo/main.go
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/joaorufino/pomo/app . 
ENTRYPOINT ["./app"]  

