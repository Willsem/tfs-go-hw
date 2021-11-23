FROM ubuntu:latest

WORKDIR /app
COPY ./api .
RUN chmod +x ./api

ENTRYPOINT ["./api", "-config", "/configs/config.toml"]

EXPOSE 8080
