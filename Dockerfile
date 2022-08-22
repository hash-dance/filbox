FROM  golang:1.13-alpine AS builder
WORKDIR /filbox-backend
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux  go build -mod=vendor -o app main.go


FROM alpine:latest
#RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /filbox-backend/conf/ssl ./conf/ssl
COPY --from=builder /filbox-backend/app .

EXPOSE 80

CMD ["./app"]