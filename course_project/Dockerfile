FROM ubuntu:latest

WORKDIR /app

COPY ./course_project/api .
RUN chmod +x ./api

ENTRYPOINT ["./api", "-config", "/configs/config.toml"]

EXPOSE 8080
