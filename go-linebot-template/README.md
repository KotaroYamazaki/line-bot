# go-find-cafe-bot

キーワードから静岡県浜松市のカフェを検索する BOT
Google Map API に対して「静岡県浜松市」の「カフェ」を検索して返すシンプルな構造
結果から LINE BOT のカルーセルメッセージを用いて送信。

# deploy

```
make deploy
```

# Setting

deploy 後に発行される http-trigger の URL を LINE BOT の Messanger API の Webhook URL に設定する。
このとき、アカウント設定の応答メッセージは OFF、Webhook メッセージは ON にしておく。
