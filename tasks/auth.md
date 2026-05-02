# 認証

## Active

- [ ] [high] Discord Developer Portal の本番設定を確認し、`DISCORD_CLIENT_ID`、`DISCORD_CLIENT_SECRET`、`DISCORD_REDIRECT_URI` を実環境の値で動作確認する
- [ ] [medium] Discord OAuth2 失敗時のレスポンスを、設定不備、認可コード不正、Discord API 障害で区別して扱う
- [ ] [medium] JWT の `JWT_SECRET` を本番環境で安全に設定する手順を文書化する

## Done

- [x] [done] Discord OAuth2 の認可コードを token endpoint で交換し、`/users/@me` から Discord user ID を取得して JWT を発行する - 2026-05-02
- [x] [done] フロントエンドの `/auth/callback` で認可コードをバックエンドへ渡し、JWT を `localStorage` に保存する - 2026-05-02

## Notes

- バックエンド実装は `backend/internal/auth/token.go` と `backend/internal/server/handler.go`。
- フロントエンド実装は `frontend/app/auth/callback/page.tsx` と `frontend/lib/auth.ts`。
- 現時点では実 Discord アプリのクライアント情報が未設定のため、実環境 OAuth2 の手動確認は未完了。
