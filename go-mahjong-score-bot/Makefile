include .env

deploy:
	gcloud functions deploy mahjong-score-api\
	 --runtime go116 \
	 --trigger-http \
	 --entry-point=Webhook \
	 --memory=256MB \
	 --set-env-vars=LINE_CHANNEL_SECRET=${LINE_CHANNEL_SECRET},LINE_CHANNEL_ACCESS_TOKEN=${LINE_CHANNEL_ACCESS_TOKEN}