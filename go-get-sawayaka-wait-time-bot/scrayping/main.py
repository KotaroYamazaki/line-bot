import os
import platform
from selenium import webdriver
from selenium.webdriver.common.by import By
import firebase_admin
from firebase_admin import firestore

firebase_admin.initialize_app()
db = firestore.client()


def init_driver():
    if platform.system() == 'Darwin':  # local
        driver = webdriver.Chrome(executable_path="./chromedriver")
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


def update_wait_time(event=None, context=None):
    driver = init_driver()
    driver.get("https://www.genkotsu-hb.com/shop/")

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

            wait_time_elm = s.find_element(By.CLASS_NAME, "boxR").find_element(
                By.CLASS_NAME, "wait_time")
            wait_time_min = wait_time_elm.find_element(
                By.CLASS_NAME, "time").find_element(
                By.CLASS_NAME, "num").text
            wait_set = wait_time_elm.find_element(
                By.CLASS_NAME, "set").find_element(
                By.CLASS_NAME, "num").text

            docs = db.collection(u'sawayaka_shops').where(
                u'name', u'==', shop_name).stream()
            for doc in docs:
                doc.reference.update(
                    {u'wait_time': wait_time_min,
                     u'wait_set': wait_set,
                     u'timestamp': firestore.SERVER_TIMESTAMP
                     })
    driver.quit()


if __name__ == '__main__':
    update_wait_time()
