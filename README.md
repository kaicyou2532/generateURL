# generateURL

Go製のシンプルなAPIサーバーです。`POST /api/v1/uploads` に画像ファイルを送信すると、アクセス可能なURLを返します。保存先はローカルディスクで、`/files/{filename}` からダウンロードできます。Dockerイメージも同梱しているため、ローカル/コンテナのどちらでもすぐ試せます。

## エンドポイント

| メソッド | パス | 説明 |
| --- | --- | --- |
| `GET /healthz` | 動作確認用のシンプルな応答を返します。 |
| `POST /api/v1/uploads` | `file` フィールドに画像を含む `multipart/form-data` を送信すると、`{"url":"..."}` を返します。 |
| `GET /files/{filename}` | アップロード済みのファイルを返します。 |

## ローカル実行 (Go)

```bash
# 依存関係取得 & テスト
go test ./...

# サーバー起動
PORT=8080 go run ./cmd/server
```

アップロード例:

```bash
curl -X POST http://localhost:8080/api/v1/uploads \
  -F "file=@sample.jpg"
```

## Docker で起動

```bash
# ビルド
docker build -t generateurl-api .

# 実行
mkdir -p uploads
docker run --rm -p 8080:8080 \
  -v $(pwd)/uploads:/app/uploads \
  -e BASE_URL=http://localhost:8080 \
  generateurl-api
```

## docker-compose

```bash
mkdir -p uploads
docker compose up --build
```

## 設定変数

| 変数 | デフォルト | 説明 |
| --- | --- | --- |
| `PORT` | `8080` | HTTPサーバーのポート。 |
| `UPLOAD_DIR` | `uploads` | ファイル保存先。存在しない場合は作成されます。 |
| `MAX_FILE_SIZE` | `10485760` (10MiB) | アップロード許容量 (バイト)。 |
| `BASE_URL` | リクエストのホスト/スキーム | 発行するURLの基点。外部公開時に設定します。 |

## テスト

```bash
go test ./...
```

## 今後の拡張案

- 認証/認可
- S3 や GCS など外部ストレージ対応
- ファイルメタデータ保存
