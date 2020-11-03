# Stocks Telegram Bot

![Tests](https://github.com/dbeliakov/stocks-bot/workflows/tests/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dbeliakov/stocks-bot)](https://goreportcard.com/report/github.com/dbeliakov/stocks-bot)
[![Coverage Status](https://coveralls.io/repos/github/dbeliakov/stocks-bot/badge.svg?branch=master)](https://coveralls.io/github/dbeliakov/stocks-bot?branch=master)


Телеграм-бот для получения информации о текущем курсе акций. Позволяет узнать текущую стоимость любой акции, а также сохранить список интересующих акций и получать актуальную информацию по ним.

Доступные команды:
* `/price` - узнать текущую стоимость акции
* `/add` - добавить акцию в список отслеживаемых
* `/remove` - удалить акцию из списка отслеживаемых
* `/my` - показать текущую цену акций из списка

#### Запуск

Для работы бота необходимо получить ключ API от [finnhub.io](finnhub.io) и токен для бота у [BotFather](t.me/BotFather).

Команда для запуска: `./stocks-bot -finnhub-token <API_KEY> -tg-token <TOKEN>`
