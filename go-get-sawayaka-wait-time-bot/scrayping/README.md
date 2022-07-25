# 待ち時間取得のための scrapying

定期実行して firestore の待ち時間情報を更新します。
local で実行するためにはブラウザに合わせた ChromeDriver を公式から取得する必要があります。

# 初期化

事前準備

- env の設定
- トピックの作成

```
make init
make scheduler
```

# デプロイ

```
make deploy
```
