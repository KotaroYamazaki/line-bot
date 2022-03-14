# 概要

LINE 上で送られたメッセージの半分を伏せ字にする LINE BOT 

![LINE_capture_668833444 134580](https://user-images.githubusercontent.com/7589567/158043206-417d1512-ea36-47f9-bdb1-e0171b6f96c7.JPG)

# 準備

- LINEBOT のチャンネルの作成
- GCP で Cloud functions の API の有効化
- Cloud functions の作成(Python3.7)

# デプロイ

## Cloud Functions のデプロイ

シークレットは LINEBOT MEssaging API の設定から確認する
実際には Secret Manager を用いて設定することが望ましい

```
gcloud functions deploy  line-bot-reply-hidden-words \
--trigger-http \
--runtime=python37 \
--entry-point=main \
--memory=256MB \
--set-env-vars=LINE_CHANNEL_SECRET={},LINE_CHANNEL_ACCESS_TOKEN={} \
--allow-unauthenticated
```
