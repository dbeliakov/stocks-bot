FROM alpine:latest
MAINTAINER Dmitrii Beliakov <dmitriy.b11@gmail.com>

RUN mkdir -p /srv/app
COPY stocks-bot /srv/app/stocks-bot

WORKDIR /src/app
CMD /srv/app/stocks-bot -finnhub-token ${FINNHUB_TOKEN} -tg-token ${TG_TOKEN}
