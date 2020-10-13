FROM golang:alpine as build
WORKDIR /splatoon2_bot
COPY . .
ENV GO11MODULE=on
ENV GOPROXY=goproxy.io
RUN apk add build-base
RUN go mod download
RUN go mod vendor
RUN go build -o splatoon2_bot splatoon2_bot.go

FROM alpine
WORKDIR /splatoon2_bot
VOLUME /splatoon2_bot/data
VOLUME /splatoon2_bot/config
COPY --from=build /splatoon2_bot/splatoon2_bot ./
#COPY --from=build /splatoon2_bot/config ./config
CMD ["./splatoon2_bot"]