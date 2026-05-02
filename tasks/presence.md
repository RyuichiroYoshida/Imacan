# Presence

## Active

- [ ] [high] 位置情報を使うか、手動在席のみで進めるかを決め、README と UI の表現を一致させる
- [ ] [medium] `CLASS`、`SELF_STUDY`、`OUT` の TTL と削除挙動を Redis 実環境で確認する
- [ ] [medium] `GET /presence/me` の期限切れ時の挙動をフロントエンドから確認する

## Done

- [x] [done] `GET /presence/me` を追加し、現在ユーザーの有効な Presence を取得できるようにする - 2026-05-02
- [x] [done] `POST /presence` と `GET /presence/summary` の基本フローを JWT 付きでテストする - 2026-05-02
- [x] [done] Redis store と memory store に単一ユーザー Presence 取得処理を追加する - 2026-05-02

## Notes

- API 契約の起点は `backend/api/typespec/main.tsp`。
- `GET /presence/me` は Presence がない、または期限切れの場合に `active: false` を返す。
- README には位置情報要件があるが、現在のフロントエンドは位置情報を要求していない。
