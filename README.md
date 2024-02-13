# bbtmvbot

[![License](https://img.shields.io/github/license/erkexzcx/bbtmvbot)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/erkexzcx/bbtmvbot)](https://goreportcard.com/report/github.com/erkexzcx/bbtmvbot)

This bot scans the most popular flat rent portals for latest posts in Vilnius, which will be sent to subscribed users using Telegram app.

Hardware requirements are so low that you can even run this completelly fine on a lowest-end SBC. On RPI0W, RAM usage is only about 8mb and CPU load is only few percent, so you can run this on any _potato_ you want :)

# Usage

Feel free to use my hosted instance in cloud: http://t.me/bbtmvbot

Otherwise see below steps on how to host it yourself.

1. Setup Telegram bot

Using [BotFather](https://t.me/BotFather) create your own Telegram bot.

Also using BotFather use command `/setcommands` and update your bot commands:
```
info - Information about BBTMV Bot
enable - Enable notifications
disable - Disable notifications
config - Configure bot settings
```
Once you set-up bot, you should have your bot's Telegram **API key**.

2. Docker

TODO. As of now I use this command, no port forward for prometheus metrics:

```
docker run --name bbtmvbot -v $(pwd)/config.yml:/config.yml -v $(pwd)/data:/data ghcr.io/erkexzcx/bbtmvbot:latest
```
