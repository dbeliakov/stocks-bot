FROM golang:1.15 as build
RUN mkdir -p /sources/
COPY . /sources/stocks-bot
RUN cd /sources/stocks-bot && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM alpine:latest
MAINTAINER Dmitrii Beliakov <dmitriy.b11@gmail.com>
RUN mkdir -p /srv/app
COPY --from=build /sources/stocks-bot/stocks-bot /srv/app/stocks-bot
WORKDIR /srv/app
EXPOSE 2112
CMD /srv/app/stocks-bot -finnhub-token ${FINNHUB_TOKEN} -tg-token ${TG_TOKEN}
