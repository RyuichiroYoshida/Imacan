# Presence

## Active

- [ ] [medium] `GET /presence/me` の期限切れ時の挙動をフロントエンドから確認する

## Done

- [x] [done] MVPでは位置情報を使わず手動在席のみで進める方針に決め、README と UI の表現を一致させる - 2026-05-02
- [x] [done] `CLASS`、`SELF_STUDY`、`OUT` の TTL と削除挙動を単体テストと Redis 実環境で確認する - 2026-05-02
- [x] [done] `GET /presence/me` を追加し、現在ユーザーの有効な Presence を取得できるようにする - 2026-05-02
- [x] [done] `POST /presence` と `GET /presence/summary` の基本フローを JWT 付きでテストする - 2026-05-02
- [x] [done] Redis store と memory store に単一ユーザー Presence 取得処理を追加する - 2026-05-02

## Notes

- MVP の在席判定は手動更新のみとする。位置情報による学校半径内判定は MVP 外の候補として README に残す。
- UI は位置情報許可を求めず、`frontend/components/PresenceDashboard.tsx` で「MVPでは位置情報を使いません」と表示する。
- TTL と削除挙動は `backend/internal/presence/service_test.go` と `backend/internal/store/redis/presence_store_test.go` で確認する。
- Redis 統合テストは `IMACAN_REDIS_INTEGRATION=1` を設定した場合だけ実行する。
- 2026-05-02 に PowerShell で `$env:IMACAN_REDIS_INTEGRATION='1'; $env:REDIS_ADDR='localhost:6379'; go test ./backend/internal/store/redis -run TestPresenceStoreWithRedisAppliesTTLAndDeletesOut -v` が成功した。
- `GET /presence/me` の期限切れ API 挙動は `backend/internal/server/router_test.go` で確認済み。フロントエンドのブラウザ操作での期限切れ表示確認は未実施。
- API 契約の起点は `backend/api/typespec/main.tsp`。
- `GET /presence/me` は Presence がない、または期限切れの場合に `active: false` を返す。
