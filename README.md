#alertmanager2tg

## Simple webhook to existing Telegram Bot for using with Prometheus AlertManager

## Usage

```
$ go get github.com/struk77/alertmanager2tg
$ cd $GOPATH/src/github.com/struk77/alertmanager2tg
$ cp tgbot.json.example tgbot.json
```

Fill the credentials in tgbots.json with your own TelegramBotAPI Key and ChatID

```
$ ./alertmanager2tg tgbot.json 8080
```