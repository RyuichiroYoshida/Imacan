# Codex用プロンプト

あなたはシニアフルスタックエンジニアです。
このリポジトリにある `README.md` は、これから作成するサービス「ProjectImacan」のMVP要件定義書です。
READMEの内容を正として、MVPを実装してください。

## 作るサービス

学校内で「今行けば誰かいるか」「自習仲間がいるか」を確認できる、在席・活動ステータス共有サービスを作成します。
ユーザーはDiscordでログインし、自分の状態を1タップで共有できます。

## MVPで実装する範囲

- Discord OAuth2ログイン
- バックエンドでのJWT発行
- JWTによるAPI認証
- 在席状態の共有
- 活動ステータスの更新
- RedisへのPresence保存
- Redis TTLによる自動失効
- 在席人数・活動別人数の集計API
- React / Next.jsによるPWA UI

## MVP外

以下は実装しないでください。

- リアルタイム位置追跡
- 地図表示
- バックグラウンド位置取得
- 詳細プロフィール
- 通知機能
- 位置履歴の保存

## 技術スタック

- フロントエンド: React / Next.js
- バックエンド: Go + TypeSpec + oapi-codegen
- 認証: Discord OAuth2 + JWT
- データストア: Redis
- 配信形態: PWA

## バックエンド要件

TypeSpecをAPI契約の起点にしてください。
TypeSpecからOpenAPIを生成し、そのOpenAPI定義をもとに `oapi-codegen` でGoの型・サーバーインターフェースを生成してAPIサーバーを実装してください。
GoのHTTPルーティングは `oapi-codegen` が生成するインターフェースに合わせ、特定フレームワーク前提の手書き実装に寄せすぎないでください。

### API定義・生成

- TypeSpec定義は `backend/api/typespec/` に配置する
- 生成されたOpenAPI定義は `backend/api/openapi/` に配置する
- `oapi-codegen` の生成コードは `backend/internal/generated/` に配置する
- APIのリクエスト・レスポンス型は原則として生成コードを利用する
- 手書きの業務ロジックは `backend/internal/auth/`, `backend/internal/presence/`, `backend/internal/store/redis/` に分離する

### 認証

- `POST /auth/discord/callback`
- Discord OAuth2のcallbackを受け取り、ユーザーを識別する
- 認証成功時にJWTを発行する
- Discord連携に必要な値は環境変数で設定できるようにする

### Presence更新

- `POST /presence`
- `Authorization: Bearer <JWT>` を必須にする
- リクエスト例:

```json
{
  "activity": "SELF_STUDY",
  "lat": 35.68,
  "lng": 139.76
}
```

- `activity` は `CLASS`, `SELF_STUDY`, `OUT` のみ許可する
- `OUT` の場合はRedisから即時削除する
- `CLASS` の場合は90〜120分程度のTTLを設定する
- `SELF_STUDY` の場合は当日20:30までのTTLを設定する
- いずれのTTLも学校開放時間の20:30を超えないようにする
- 位置情報は正確な値を永続保存しない
- 位置情報がない場合でも手動在席を許可する

### Presence集計

- `GET /presence/summary`
- レスポンス例:

```json
{
  "total": 5,
  "class": 3,
  "self_study": 2
}
```

- 期限切れのPresenceは集計に含めない

### データモデル

READMEの以下のモデルを基準にしてください。

```go
type Presence struct {
    UserID    string
    Activity  string
    UpdatedAt time.Time
    ExpiresAt time.Time
}
```

## フロントエンド要件

Next.jsでPWAとして使えるUIを実装してください。

### 画面要件

- Discordログイン導線
- 現在の在席人数
- 授業中人数
- 自習中人数
- 「授業中」「自習中」「帰宅」の状態更新ボタン
- ステータス更新後に有効期限を表示
- 再訪時に現在の状態が残っている場合は継続確認を表示

### UX要件

- 状態変更は1タップで完了する
- 状態は更新後すぐに画面へ反映する
- 初回のみ位置情報の許可を求める
- 位置情報が取得できない場合でも手動更新できる
- 操作数を最小限にする
- 正確な位置情報や履歴をユーザーに不安に感じさせない設計にする

## セキュリティ・プライバシー要件

- JWT認証を必須にする
- 本番運用ではHTTPS前提にする
- 正確な位置情報は保存しない、または丸める
- 位置履歴は保持しない
- 位置情報利用はオプトインにする

## 実装方針

1. 既存のリポジトリ構成を確認してください。
2. 必要なディレクトリ・ファイルを作成してください。
3. TypeSpecでAPI契約を定義し、OpenAPIとGoコード生成の流れを作ってください。
4. バックエンド、フロントエンド、Redis連携をMVPとして動く形にしてください。
5. ローカル開発用の環境変数サンプルを用意してください。
6. READMEに起動方法、環境変数、API概要、API生成手順を追記してください。
7. 可能な範囲でテストまたは動作確認コマンドを追加してください。

## 受け入れ条件

- Discordログイン後にJWTを取得できる
- TypeSpecからOpenAPIを生成できる
- OpenAPIから `oapi-codegen` でGoコードを生成できる
- JWT付きでPresenceを更新できる
- `CLASS`, `SELF_STUDY`, `OUT` の状態更新ができる
- Redis TTLにより状態が自動失効する
- 集計APIで在席人数・活動別人数が取得できる
- PWA UIから状態変更と集計確認ができる
- 位置情報が取得できなくても手動で状態更新できる
- MVP外の機能を過剰に実装していない

まずリポジトリを調査し、実装計画を簡潔に提示したうえで作業を開始してください。
