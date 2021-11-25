FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY ./course_project/api .
RUN chmod +x ./api

ENTRYPOINT ["./api", "-config", "/configs/config.toml"]

EXPOSE 8080
