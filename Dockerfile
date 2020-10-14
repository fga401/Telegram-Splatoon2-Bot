FROM golang:alpine as build
WORKDIR /splatoon2_bot
COPY . .
ENV GO11MODULE=on
ENV GOPROXY=https://goproxy.cn
RUN apk add build-base
RUN go build -o splatoon2_bot splatoon2_bot.go

FROM alpine
WORKDIR /splatoon2_bot
VOLUME /splatoon2_bot/data
VOLUME /splatoon2_bot/config
COPY --from=build /splatoon2_bot/splatoon2_bot ./
CMD ["./splatoon2_bot"]