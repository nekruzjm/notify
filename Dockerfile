FROM alpine:3
#FROM --platform=linux/arm64 arm64v8/alpine:3

WORKDIR /app

RUN apk --no-cache add ca-certificates \
       openssl

COPY /cmd/notifications/main .
ADD configs ./configs

EXPOSE 80
CMD ["./main"]