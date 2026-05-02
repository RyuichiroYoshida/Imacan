# 認証運用メモ

## 目的

ProjectImacan の Discord OAuth2 と JWT 署名鍵を、開発環境と本番環境で安全に設定するための判断基準と確認手順をまとめる。

## 対象

- `DISCORD_CLIENT_ID`
- `DISCORD_CLIENT_SECRET`
- `DISCORD_REDIRECT_URI`
- `JWT_SECRET`
- `NEXT_PUBLIC_DISCORD_CLIENT_ID`
- `NEXT_PUBLIC_DISCORD_REDIRECT_URI`

## Discord OAuth2

Discord Developer Portal でアプリケーションを作成し、OAuth2 の redirect URI にフロントエンドの callback URL を登録する。

開発環境の例:

```text
http://localhost:3000/auth/callback
```

本番環境では、実際に公開する HTTPS の URL を登録する。

```text
https://<frontend-domain>/auth/callback
```

バックエンドには `DISCORD_CLIENT_ID`、`DISCORD_CLIENT_SECRET`、`DISCORD_REDIRECT_URI` を設定する。フロントエンドには公開してよい `NEXT_PUBLIC_DISCORD_CLIENT_ID` と `NEXT_PUBLIC_DISCORD_REDIRECT_URI` を設定する。

## JWT Secret

`JWT_SECRET` は JWT の署名に使う秘密値であり、リポジトリにコミットしない。本番環境では、ホスティング基盤の secret manager または環境変数管理機能に保存する。

生成例:

```bash
openssl rand -base64 32
```

Windows PowerShell で生成する例:

```powershell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))
```

## 確認方法

1. `.env` またはデプロイ環境に `DISCORD_CLIENT_ID`、`DISCORD_CLIENT_SECRET`、`DISCORD_REDIRECT_URI`、`JWT_SECRET` を設定する。
2. フロントエンド環境に `NEXT_PUBLIC_DISCORD_CLIENT_ID` と `NEXT_PUBLIC_DISCORD_REDIRECT_URI` を設定する。
3. Redis、バックエンド、フロントエンドを起動する。
4. ブラウザで Discord ログインを行う。
5. ログイン後、`POST /auth/discord/callback` が `accessToken` を返すことを確認する。
6. JWT 付きで `GET /presence/me` または `POST /presence` が成功することを確認する。

## 失敗時の見分け方

`POST /auth/discord/callback` は、Discord OAuth2 の失敗理由を次の `code` で返す。

- `DISCORD_OAUTH_NOT_CONFIGURED`: バックエンドの Discord OAuth2 環境変数が不足している。
- `DISCORD_AUTH_FAILED`: 認可コードが不正、期限切れ、または redirect URI が一致していない。
- `DISCORD_API_UNAVAILABLE`: Discord API への接続、token endpoint、user endpoint のいずれかが利用できない。
- `DISCORD_AUTH_ERROR`: 上記以外の認証処理エラー。

## 注意点

- `JWT_SECRET` を変更すると、既存の JWT は検証できなくなる。
- `DISCORD_CLIENT_SECRET` と `JWT_SECRET` はフロントエンドへ渡さない。
- `NEXT_PUBLIC_` が付く環境変数はブラウザに公開される前提で扱う。
- 本番環境では HTTPS を前提とし、Discord Developer Portal の redirect URI も HTTPS の URL にする。
