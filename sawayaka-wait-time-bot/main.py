from flask import abort, jsonify

import os
from selenium import webdriver
from selenium.webdriver.common.by import By
import platform
import hashlib
import hmac
import base64

from linebot import (
    LineBotApi, WebhookParser
)
from linebot.exceptions import (
    InvalidSignatureError
)
from linebot.models import (
    MessageEvent, TextMessage, TextSendMessage
)


def init_driver():
    if platform.system() == 'Darwin':
        driver = webdriver.Chrome(executable_path="./bin/local_chromedriver")
    else:
        chrome_options = webdriver.ChromeOptions()
        chrome_options.add_argument('--headless')
        chrome_options.add_argument('--disable-gpu')
        chrome_options.add_argument('--window-size=1280x1696')
        chrome_options.add_argument('--no-sandbox')
        chrome_options.add_argument('--hide-scrollbars')
        chrome_options.add_argument('--enable-logging')
        chrome_options.add_argument('--log-level=0')
        chrome_options.add_argument('--v=99')
        chrome_options.add_argument('--single-process')
        chrome_options.add_argument('--ignore-certificate-errors')
        chrome_options.add_argument(
            'user-agent=Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/32.0.1667.0 Safari/537.36')
        chrome_options.binary_location = os.getcwd() + "/bin/headless-chromium"
        driver = webdriver.Chrome(
            os.getcwd() + "/bin/chromedriver", chrome_options=chrome_options)
    return driver


def get_wait_time(keyword):
    driver = init_driver()
    driver.get("https://www.genkotsu-hb.com/shop/")

    msg = []
    is_find = False

    shop_list = driver.find_element(By.ID, "shop_list")
    shop_boxes = shop_list.find_elements(By.CLASS_NAME, "bg_shadow")
    for s in shop_boxes:
        if (s.get_attribute("id") == "shop_navi"):
            continue

        shops = s.find_element(By.CLASS_NAME, "areaBox").find_elements(
            By.CLASS_NAME, "shop_info")
        for s in shops:
            shop_name = s.find_element(By.CLASS_NAME, "shop_box").find_element(
                By.CLASS_NAME, "shop_name").find_element(By.CLASS_NAME, "name").text[:-1]  # exclude "店"

            if not (keyword in shop_name):
                continue

            is_find = True
            wait_time_elm = s.find_element(By.CLASS_NAME, "boxR").find_element(
                By.CLASS_NAME, "wait_time")
            wait_time_min = wait_time_elm.find_element(
                By.CLASS_NAME, "time").find_element(
                By.CLASS_NAME, "num").text
            wait_set = wait_time_elm.find_element(
                By.CLASS_NAME, "set").find_element(
                By.CLASS_NAME, "num").text

            if wait_time_min == '-':
                msg.append(f"{shop_name}店は現在営業していない可能性があります")
            else:
                msg.append(
                    f"{shop_name}店のただいまの待ち時間: 約{wait_time_min}分 {wait_set}組待ち")

    if is_find == False:
        msg.append(f"検索キーワード:'{keyword}'の店舗は見つかりませんでした")

    driver.quit()

    return "\n".join(msg)


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

    for event in events:
        if not isinstance(event, MessageEvent):
            continue
        if not isinstance(event.message, TextMessage):
            continue

        line_bot_api.reply_message(
            event.reply_token,
            TextSendMessage(
                text=get_wait_time(event.message.text)
            ))
