# フロントエンド PWA

## Active

- [ ] [high] Discord Developer Portal の redirect URI と `NEXT_PUBLIC_DISCORD_CLIENT_ID` を設定し、ブラウザでログインから状態更新まで通す
- [ ] [medium] PWA として必要なアイコン、manifest、インストール時の表示を整える
- [ ] [medium] API 未起動、Redis 未起動、JWT 期限切れ時の画面表示を確認して文言を調整する
- [ ] [low] フロントエンドのテスト方針を決める

## Done

- [x] [done] MVPでは位置情報許可を求めず、手動更新のみで進める方針を UI に反映する - 2026-05-02
- [x] [done] Next.js の最小構成を追加する - 2026-05-02
- [x] [done] Discord callback で JWT を保存する - 2026-05-02
- [x] [done] トップ画面で `/presence/me` を呼び、再訪時に継続確認を表示する - 2026-05-02
- [x] [done] 在席人数、授業中人数、自習中人数、状態更新ボタンを表示する - 2026-05-02

## Notes

- フロントエンドの入口は `frontend/app/page.tsx`。
- 主要 UI は `frontend/components/PresenceDashboard.tsx`。
- `npm.cmd run build` は成功済み。
- MVPでは位置情報を扱わない。手動更新のみで状態を共有する。
