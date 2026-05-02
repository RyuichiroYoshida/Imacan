# 開発ワークフロー

## Active

- [ ] [medium] Task CLI の導入手順を README または `docs/` に追記する
- [ ] [medium] `task dev:redis` を、既存の停止済みコンテナがある場合にも再利用または復旧できる形にする
- [ ] [medium] README に `Taskfile.yml` を使った起動、生成、検証手順を追記する
- [ ] [low] `npm audit` の moderate 脆弱性を確認し、Next.js 依存更新の影響を判断する

## Done

- [x] [done] Railway 向けの backend / frontend Dockerfile と Config as Code を追加する - 2026-05-02
- [x] [done] Go API を Railway の `PORT` と `REDIS_URL` に対応させる - 2026-05-02
- [x] [done] Railway デプロイ手順と環境変数を `docs/infrastructure.md` に追記する - 2026-05-02
- [x] [done] Docker Compose で Redis、Go API、Next.js フロントエンドを起動できる開発環境を追加する - 2026-05-02
- [x] [done] Docker Compose 開発環境の起動、停止、ログ確認、healthcheck 手順を `docs/docker-compose.md` に文書化する - 2026-05-02
- [x] [done] `Taskfile.yml` に `compose:up`、`compose:up:detached`、`compose:down`、`compose:logs` を追加する - 2026-05-02
- [x] [done] `Taskfile.yml` を作成し、依存インストール、生成、テスト、ビルド、起動コマンドを整理する - 2026-05-02
- [x] [done] `Taskfile.yml` の `desc` を日本語に統一する - 2026-05-02
- [x] [done] `.gitignore` に `.next/` と開発サーバーログを追加する - 2026-05-02

## Notes

- Docker Compose 定義は `docker-compose.yml`、利用手順は `docs/docker-compose.md` に置いた。
- Compose の backend サービスでは `REDIS_ADDR=redis:6379`、`REDIS_PASSWORD=` を上書きする。ホスト実行用の `.env` は `REDIS_ADDR=localhost:6379` のままでよい。
- frontend サービスには secret を含む `.env` 全体を渡さず、`NEXT_PUBLIC_API_BASE_URL` と `NEXT_PUBLIC_DISCORD_*` だけを渡す。
- `docker compose config --no-interpolate` は成功した。Docker 設定ファイル `C:\Users\yohir\.docker\config.json` の access denied warning は出たが、Compose 定義の構文確認はできている。
- この環境では `task --list` 実行時に Task CLI が未インストールだった。
- `task dev` は Redis 起動後にバックエンドとフロントエンドを並列起動する想定。
- `npm install` 実行時に moderate 脆弱性が 2 件報告されたが、破壊的更新は未実行。
- Railway 用ファイルは `deployments/railway/backend.Dockerfile`、`deployments/railway/frontend.Dockerfile`、`deployments/railway/*.railway.json` に置いた。
- Railway の service settings では config file path を `/deployments/railway/backend.railway.json` と `/deployments/railway/frontend.railway.json` にする。
- Railway Redis は `REDIS_URL` を優先して接続する。ローカルと Docker Compose は従来通り `REDIS_ADDR` を使える。
