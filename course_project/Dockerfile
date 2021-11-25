FROM ubuntu:latest

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app
COPY ./course_project/api .
RUN chmod +x ./api

ENTRYPOINT ["./api", "-config", "/configs/config.toml"]

EXPOSE 8080
