# フロントエンド PWA

## Active

- [ ] [high] Discord Developer Portal の redirect URI と `NEXT_PUBLIC_DISCORD_CLIENT_ID` を設定し、ブラウザでログインから状態更新まで通す
- [ ] [medium] PWA のインストール表示を実ブラウザで確認する
- [ ] [medium] API 未起動、Redis 未起動、JWT 期限切れ時の画面表示を実ブラウザで確認する

## Done

- [x] [done] PWA として必要な manifest、theme color、アプリアイコンを整える - 2026-05-02
- [x] [done] API 未起動、Redis 未起動、JWT 期限切れ時のフロントエンド文言を調整する - 2026-05-02
- [x] [done] フロントエンドの当面のテスト方針を、`npm.cmd run build` と localhost smoke 確認を基本にすると決める - 2026-05-02
- [x] [done] MVPでは位置情報許可を求めず、手動更新のみで進める方針を UI に反映する - 2026-05-02
- [x] [done] Next.js の最小構成を追加する - 2026-05-02
- [x] [done] Discord callback で JWT を保存する - 2026-05-02
- [x] [done] トップ画面で `/presence/me` を呼び、再訪時に継続確認を表示する - 2026-05-02
- [x] [done] 在席人数、授業中人数、自習中人数、状態更新ボタンを表示する - 2026-05-02

## Notes

- フロントエンドの入口は `frontend/app/page.tsx`。
- 主要 UI は `frontend/components/PresenceDashboard.tsx`。
- `npm.cmd run build` は 2026-05-02 に成功済み。
- MVPでは位置情報を扱わない。手動更新のみで状態を共有する。
- PWA アイコンは `frontend/public/icon.svg`、manifest は `frontend/public/manifest.webmanifest`。
- `http://127.0.0.1:3001` で dev server を起動し、トップページ、`/manifest.webmanifest`、`/icon.svg` が 200 を返すことを確認した。`3000` は既に使用中だった。
- Discord Developer Portal の設定と実ログイン確認は外部操作が必要なため Active に残す。
