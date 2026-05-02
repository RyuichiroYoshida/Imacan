# 認証

## Active

- [ ] [high] Discord Developer Portal の本番設定を確認し、`DISCORD_CLIENT_ID`、`DISCORD_CLIENT_SECRET`、`DISCORD_REDIRECT_URI` を実環境の値で動作確認する blocked: Discord Developer Portal の登録内容確認とブラウザでの実ログイン操作が必要

## Done

- [x] [done] Discord OAuth2 失敗時のレスポンスを、設定不備、認可コード不正、Discord API 障害で区別して扱う - 2026-05-02
- [x] [done] JWT の `JWT_SECRET` を本番環境で安全に設定する手順を文書化する - 2026-05-02
- [x] [done] `.env` に認証関連の環境変数が設定されていることを、値を表示せずに確認する - 2026-05-02
- [x] [done] Discord OAuth2 の認可コードを token endpoint で交換し、`/users/@me` から Discord user ID を取得して JWT を発行する - 2026-05-02
- [x] [done] フロントエンドの `/auth/callback` で認可コードをバックエンドへ渡し、JWT を `localStorage` に保存する - 2026-05-02

## Notes

- バックエンド実装は `backend/internal/auth/token.go` と `backend/internal/server/handler.go`。
- フロントエンド実装は `frontend/app/auth/callback/page.tsx` と `frontend/lib/auth.ts`。
- 認証運用手順は `docs/auth-operations.md`。
- `.env` には `JWT_SECRET`、`DISCORD_CLIENT_ID`、`DISCORD_CLIENT_SECRET`、`DISCORD_REDIRECT_URI` が設定済み。ただし値はタスク管理に記録しない。
- 実 Discord アプリの redirect URI 登録内容と、ブラウザからの実ログイン成功は未確認。
