# context <splatoon2_bot_project_dir>
FROM golang as build
WORKDIR /splatoon2_bot
COPY . .
RUN rm -rf /splatoon2_bot/data
VOLUME /splatoon2_bot/data
VOLUME /splatoon2_bot/migrate/sqls
ENV GO11MODULE=on
ENV GOPROXY=https://goproxy.cn
RUN go build -o ./migrate/migrate ./migrate/migrate.go
WORKDIR /splatoon2_bot/migrate
CMD ["./migrate"]

