# 開発ワークフロー

## Active

- [ ] [medium] Task CLI の導入手順を README または `docs/` に追記する
- [ ] [medium] `task dev:redis` を、既存の停止済みコンテナがある場合にも再利用または復旧できる形にする
- [ ] [medium] README に `Taskfile.yml` を使った起動、生成、検証手順を追記する
- [ ] [low] `npm audit` の moderate 脆弱性を確認し、Next.js 依存更新の影響を判断する

## Done

- [x] [done] `Taskfile.yml` を作成し、依存インストール、生成、テスト、ビルド、起動コマンドを整理する - 2026-05-02
- [x] [done] `Taskfile.yml` の `desc` を日本語に統一する - 2026-05-02
- [x] [done] `.gitignore` に `.next/` と開発サーバーログを追加する - 2026-05-02

## Notes

- この環境では `task --list` 実行時に Task CLI が未インストールだった。
- `task dev` は Redis 起動後にバックエンドとフロントエンドを並列起動する想定。
- `npm install` 実行時に moderate 脆弱性が 2 件報告されたが、破壊的更新は未実行。
