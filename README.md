# generateURL

Go製のシンプルなAPIサーバーです。`POST /api/v1/uploads` に画像ファイルを送信すると、アクセス可能なURLを返します。保存先はローカルディスクで、`/files/{filename}` からダウンロードできます。Dockerイメージも同梱しているため、ローカル/コンテナのどちらでもすぐ試せます。

## エンドポイント

| メソッド | パス | 説明 |
| --- | --- | --- |
| `GET /healthz` | 動作確認用のシンプルな応答を返します。 |
| `POST /api/v1/uploads` | `file` フィールドに画像を含む `multipart/form-data` を送信すると、`{"url":"..."}` を返します。 |
  -e BASE_URL=http://133.2.37.149/ig \

## ローカル実行 (Go)

```bash
# 依存関係取得 & テスト
go test ./...

# サーバー起動
PORT=8000 go run ./cmd/server
```

アップロード例:

```bash
curl -X POST http://localhost:8000/api/v1/uploads \
  -F "file=@sample.jpg"
```

## Docker で起動

```bash
# ビルド
docker build -t generateurl-api .

# 実行
mkdir -p uploads
docker run --rm -p 8000:8000 \
  -v $(pwd)/uploads:/app/uploads \
  -e BASE_URL=http://localhost:8000 \
  generateurl-api

# Linux ホストで権限エラーが出る場合
# distroless の実行ユーザー (uid/gid 65532) に書き込み権限を与えてください
sudo chown -R 65532:65532 uploads
```

## docker-compose

```bash
mkdir -p uploads
docker compose up --build

# すでに uploads ディレクトリをホストに作成している場合は書き込み権限を付与
sudo chown -R 65532:65532 uploads
```

## すぐ試す（ブラウザ）

プロジェクトルートに `test_upload.html` を追加しました。サーバーが http://localhost:8000 で動作していれば、ブラウザで `file://` 経由、あるいはシンプルな静的サーバー（`python -m http.server` など）で開いてアップロードを試せます。

## すぐ試す（スクリプト）

`scripts/upload.sh` を追加しました。使い方:

```bash
# 実行権を付与
chmod +x scripts/upload.sh

# 例: ローカルサーバーにアップロード
scripts/upload.sh http://localhost:8000 ./sample.jpg
```

## リモート環境での確認例

現在 `http://133.2.37.149/ig` にデプロイ済みの場合、ブラウザ版/スクリプト版ともにベースURL欄にそのまま入力すれば動作します。

```bash
# CLI から直接リモート環境へPOST
scripts/upload.sh http://133.2.37.149/ig ./sample.jpg
```

`test_upload.html` の入力欄にも `http://133.2.37.149/ig` をセットすると同様に試験できます。

## 設定変数

| 変数 | デフォルト | 説明 |
| --- | --- | --- |
| `PORT` | `8000` | HTTPサーバーのポート。 |
| `UPLOAD_DIR` | `uploads` | ファイル保存先。存在しない場合は作成されます。 |
| `MAX_FILE_SIZE` | `10485760` (10MiB) | アップロード許容量 (バイト)。 |
| `BASE_URL` | リクエストのホスト/スキーム | 発行するURLの基点。外部公開時に設定します。 |

## テスト

```bash
go test ./...
```

