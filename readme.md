# 環境構築

## Nginx に使うキーを生成する
```
chomd +x ./Genkey.sh
./Genkey.sh
```

## docker コンテナを立ち上げる
```
docker compose up -d --build
```