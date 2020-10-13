# context <splatoon2_bot_project_dir>
FROM golang as build
WORKDIR /splatoon2_bot
VOLUME /splatoon2/data
COPY . .
ENV GO11MODULE=on
ENV GOPROXY=goproxy.io
RUN go mod download
RUN go mod vendor
RUN go build -o ./migrate/migrate ./migrate/migrate.go
WORKDIR /splatoon2_bot/migrate
CMD ["./migrate"]

