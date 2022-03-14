# 概要

とりあえず褒めてくれる BOT の

# 準備

- LINEBOT のチャンネルの作成
- GCP で Cloud functions の API の有効化
- Cloud functions の作成(Python3.7)

# デプロイ

## Cloud Functions のデプロイ

シークレットは LINEBOT MEssaging API の設定から確認する
実際には Secret Manager を用いて設定することが望ましい

```
gcloud functions deploy  homeru-bot-py \
--trigger-http \
--runtime=python37 \
--entry-point=main \
--memory=256MB \
--set-env-vars=LINE_CHANNEL_SECRET="",LINE_CHANNEL_ACCESS_TOKEN="" \
--allow-unauthenticated
```
