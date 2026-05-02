# Docker Compose 開発環境

## 目的

ProjectImacan のローカル開発に必要な Redis、Go API、Next.js フロントエンドを `docker compose` でまとめて起動できるようにする。

この文書では、起動方法、環境変数、確認方法、停止方法を判断できる状態にする。

## 背景

既存のローカル開発では `task dev` が Redis、バックエンド、フロントエンドを起動する。Docker Compose 環境では同じ構成をコンテナ上で再現し、ホスト側のブラウザから次の URL にアクセスする。

- フロントエンド: `http://localhost:3000`
- バックエンド API: `http://localhost:8080`
- Redis: `localhost:6379`

## 方針

- `docker-compose.yml` は開発用のコンテナ環境として扱う。
- Redis は `redis:7-alpine` を使用する。
- バックエンドは `golang:1.26-alpine` コンテナで `go run ./backend/cmd/api` を実行する。
- フロントエンドは `node:24-alpine` コンテナで `npm run dev` を実行する。
- バックエンドコンテナ内では Redis の接続先を `redis:6379` に上書きする。
- フロントエンドの `NEXT_PUBLIC_API_BASE_URL` はブラウザから到達できる `http://localhost:8080` にする。
- `.env` が未作成でもコンテナ起動自体はできるようにし、Discord OAuth2 を使う場合だけ `.env` に認証情報を設定する。

## 構成

| サービス | 役割 | 公開ポート | 主な設定 |
| --- | --- | --- | --- |
| `redis` | Presence 保存用 Redis | `6379` | named volume `redis-data` に保存 |
| `backend` | Go API | `8080` | `REDIS_ADDR=redis:6379` |
| `frontend` | Next.js 開発サーバー | `3000` | `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080` |

## 手順

`.env` が必要な場合は `.env.example` から作成する。

```powershell
Copy-Item .env.example .env
```

Discord OAuth2 ログインを確認する場合は、`.env` に次の値を設定する。

```text
DISCORD_CLIENT_ID=
DISCORD_CLIENT_SECRET=
DISCORD_REDIRECT_URI=http://localhost:3000/auth/callback
NEXT_PUBLIC_DISCORD_CLIENT_ID=
NEXT_PUBLIC_DISCORD_REDIRECT_URI=http://localhost:3000/auth/callback
```

コンテナ環境を起動する。

```powershell
docker compose up
```

バックグラウンドで起動する場合は次を使う。

```powershell
docker compose up -d
```

Taskfile 経由でも起動できる。

```powershell
task compose:up
```

ログを確認する。

```powershell
docker compose logs -f
```

停止する。

```powershell
docker compose down
```

コンテナと named volume をまとめて削除する。

```powershell
docker compose down -v
```

## 確認方法

バックエンドの healthcheck を確認する。

```powershell
Invoke-WebRequest http://localhost:8080/healthz
```

フロントエンドをブラウザで確認する。

```text
http://localhost:3000
```

`docker compose config --no-interpolate` で Compose 定義の解釈を確認できる。

```powershell
docker compose config --no-interpolate
```

## 注意点

- 初回起動時は Go module と npm package の取得が発生するため時間がかかる。
- `.env` の `REDIS_ADDR=localhost:6379` はホスト実行用の値として残してよい。Compose では backend サービス側で `redis:6379` に上書きする。
- `.env` の `REDIS_PASSWORD` は Compose の Redis では使わないため、backend サービス側で空文字に上書きする。
- frontend サービスには `DISCORD_CLIENT_SECRET` や `JWT_SECRET` を渡さない。公開してよい `NEXT_PUBLIC_` の値だけを渡す。
- Docker の設定ファイル権限に問題がある場合、`docker compose config` で `C:\Users\yohir\.docker\config.json` の access denied warning が出ることがある。Compose 定義の表示自体が成功していれば、定義の構文確認は完了している。
