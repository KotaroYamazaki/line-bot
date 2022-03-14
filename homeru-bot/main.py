import os
import base64
import hashlib
import hmac
import random

from flask import abort, jsonify

from linebot import (
    LineBotApi, WebhookParser
)
from linebot.exceptions import (
    InvalidSignatureError
)
from linebot.models import (
    MessageEvent, TextMessage, TextSendMessage
)


def main(request):
    channel_secret = os.environ.get('LINE_CHANNEL_SECRET')
    channel_access_token = os.environ.get('LINE_CHANNEL_ACCESS_TOKEN')

    line_bot_api = LineBotApi(channel_access_token)
    parser = WebhookParser(channel_secret)

    body = request.get_data(as_text=True)
    hash = hmac.new(channel_secret.encode('utf-8'),
                    body.encode('utf-8'), hashlib.sha256).digest()
    signature = base64.b64encode(hash).decode()

    if signature != request.headers['X_LINE_SIGNATURE']:
        return abort(405)

    try:
        events = parser.parse(body, signature)
    except InvalidSignatureError:
        return abort(405)

    home_words = [
        "神ですね〜",
        "今日もすごくえらい！",
        "魅力がすごいね！",
        "いい！本当にいい！",
        "良〜〜〜〜！",
        "いや〜もってるな〜！",
        "最高じゃん！",
        "愛され体質だね！",
        "ちょっとまって…良すぎ！",
        "ねえむり良すぎ",
        "う〜ん、センスがあるなあ",
        "はぁ…すごくいい…",
        "世界一えらい！",
        "さすがです！",
        "人として尊敬するよ！",
        "💘💘💘💘💘最高💘💘💘💘💘",
        "センスの塊ですか？",
        "1031(天才)！",
        "100点満点中1000000000000点だね！",
        "リスペクトだよ！",
        "尊い〜〜〜〜〜〜！！！",
    ]

    for event in events:
        if not isinstance(event, MessageEvent):
            continue
        if not isinstance(event.message, TextMessage):
            continue

        line_bot_api.reply_message(
            event.reply_token,
            TextSendMessage(
                text=home_words[random.randint(0, len(home_words)-1)])
        )
    return jsonify({'message': 'ok'})
