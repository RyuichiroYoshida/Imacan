# インフラ環境検討

## 目的

ProjectImacan の MVP を Railway にデプロイする方針と、必要な設定・未決事項を整理する。

この文書では、ローカル開発、検証環境、本番環境の分け方、Railway での構成、環境変数、運用上の注意点を整理する。AWS などの他候補は、将来の再検討用として比較情報だけを残す。

## 背景

現在の実装は次の構成になっている。

- フロントエンド: Next.js PWA
- バックエンド: Go API
- データストア: Redis 互換ストア
- 認証: Discord OAuth2 + JWT
- ローカル開発: `docker-compose.yml` または `Taskfile.yml`

Presence は Redis TTL で自然失効する設計のため、MVP では永続 DB よりも Redis 互換ストアの可用性、接続しやすさ、運用の軽さを優先する。

## 方針

- MVP のデプロイ先は Railway に決定する。
- MVP は小さく始め、フロントエンド、API、Redis の 3 要素だけを運用対象にする。
- 本番環境は HTTPS を必須にする。
- `JWT_SECRET`、`DISCORD_CLIENT_SECRET`、Redis 接続情報はホスティング基盤の secret / environment variables に置く。
- `NEXT_PUBLIC_` が付く値だけをブラウザ公開前提の設定として扱う。
- API と Redis は同じリージョン、または同じプラットフォーム内の private network に寄せる。
- 位置情報、履歴、通知など MVP 外の機能を前提にした構成にはしない。

## 決定した構成

### MVP 本番

MVP 本番は `Railway` の単一プロジェクト構成にする。

```text
[Next.js PWA]
  |
  | HTTPS
  v
[Go API]
  |
  | private network
  v
[Redis compatible service]
```

理由:

- フロントエンド、Go API、Redis を同じプロジェクト内にまとめやすい。
- 環境変数をサービス単位または共有変数として管理できる。
- private networking により、API から Redis への通信を公開せずに済む。
- MVP の小規模な検証では、インフラの分割よりも設定の少なさが効く。

注意:

- ブラウザで動く PWA は private network に直接アクセスできない。`NEXT_PUBLIC_API_BASE_URL` は公開 API URL にする。
- Go API はホスティング基盤が指定するポートに bind する必要がある。必要なら `API_ADDR` を `:8080` 固定ではなく `PORT` から組み立てる修正を検討する。
- Discord Developer Portal の redirect URI は本番 URL が確定してから HTTPS の URL を登録する。

### Railway サービス構成

Railway project には次のサービスを作る。

| サービス | 役割 | 公開 |
| --- | --- | --- |
| `frontend` | Next.js PWA | 公開する |
| `backend` | Go API | 公開する |
| `redis` | Presence 保存 | 公開しない |

`frontend` から `backend` へは Railway が発行する公開 HTTPS URL でアクセスする。`backend` から `redis` へは Railway project 内の private network で接続する。

## 代替候補

Railway に決定したため、以下は現時点では採用しない。コスト、運用、要件が変わった場合の再検討材料として残す。

### AWS を使う場合

月 5,000 円以内、かつ極力無料枠に収めるなら、代替案は `Amplify Hosting + API Gateway HTTP API + Lambda + DynamoDB TTL` になる。

```text
[Amplify Hosting / Next.js PWA]
  |
  | HTTPS
  v
[API Gateway HTTP API]
  |
  v
[Lambda / Go API]
  |
  v
[DynamoDB / Presence with TTL]
```

理由:

- Amplify Hosting、API Gateway、Lambda、DynamoDB は小規模アクセス時に無料枠へ寄せやすい。
- NAT Gateway、常時稼働 EC2、常時稼働 Redis を避けられる。
- Presence は一時データなので、DynamoDB TTL と `expiresAt` フィルタで MVP 要件を満たしやすい。
- 月 5,000 円は、2026-05-02 時点の為替目安で約 30 USD 前後のため、常時稼働リソースを増やさないことが重要になる。

ただし、現在の実装は Redis を前提にしているため、この構成では `backend/internal/store/redis/` と同じ interface を満たす DynamoDB store を追加する必要がある。

DynamoDB TTL は削除が即時ではないため、集計 API では TTL 削除だけに頼らず、`expiresAt <= now` の Presence を除外して集計する。これは現在の Redis TTL の「期限切れは集計に含めない」という要件を AWS 向けに守るための補正である。

AWS で Redis 互換を維持する場合は、次のどちらかにする。

- `Lightsail` 上で Docker Compose を動かし、API と Redis を同一インスタンスに置く。
- `ElastiCache Serverless for Valkey` を使う。

前者は実装変更が少なく、月額が読みやすい。後者は運用は軽いが、Lambda から VPC 内の ElastiCache に接続しつつ Discord API へ外向き通信も必要になるため、NAT Gateway を入れると月 5,000 円を超えやすい。無料枠重視の MVP では DynamoDB 案を優先する。

### ステージング

本番前に `staging` 環境を 1 つ用意する。

- `main` または `production` ブランチ: 本番
- `staging` ブランチ: 検証
- Discord OAuth2 app は、本番用と検証用を分ける
- Redis は本番と検証で分ける
- `JWT_SECRET` は環境ごとに別の値にする

検証環境では次を確認する。

- Discord ログインから JWT 発行まで通る
- `POST /presence` が JWT 付きで成功する
- `GET /presence/summary` が TTL 切れの Presence を含めない
- PWA manifest と icon が HTTPS 上で取得できる

### ローカル開発

ローカルでは既存の `docker-compose.yml` を開発環境として使う。

```powershell
docker compose up
```

または Taskfile 経由で起動する。

```powershell
task compose:up
```

詳細は [Docker Compose 開発環境](docker-compose.md) を参照する。

## 候補比較

| 候補 | 構成 | 向いている状況 | 懸念 |
| --- | --- | --- | --- |
| Railway | Next.js、Go API、Redis を同一プロジェクト | MVP を早く公開したい | 長期運用時のコスト見積もりは別途確認 |
| AWS serverless | Amplify、API Gateway、Lambda、DynamoDB TTL | 月 5,000 円以内、無料枠重視で運用したい | Redis store から DynamoDB store への実装変更が必要 |
| AWS Lightsail | 1 台に Docker Compose、Redis も同居 | 現在の実装をなるべく変えず AWS で動かしたい | OS、TLS、Redis 永続化、バックアップを自分で見る必要がある |
| Render | Web Service + Key Value | 管理画面で堅実に運用したい | フロントエンドと API の構成をどう分けるか先に決める必要がある |
| Fly.io | Go API + Upstash Redis、必要なら frontend も配置 | リージョンや低レイテンシを意識したい | 初回設定に CLI とネットワーク知識がやや必要 |
| Vercel + API 基盤 + Redis | Next.js は Vercel、API と Redis は別基盤 | PWA 配信を Next.js 最適化に寄せたい | 運用対象が複数サービスに分かれる |
| VPS + Docker Compose | 1 台に frontend、backend、Redis | 学習目的、最小コスト重視 | TLS、バックアップ、監視、障害対応を自分で持つ |

## AWS 予算案

### 予算前提

- 月額上限は 5,000 円とする。
- 2026-05-02 時点の USD/JPY はおおむね 157 円前後のため、5,000 円は約 31 USD として見る。
- AWS の表示価格は税抜き USD が多いため、日本の消費税、為替変動、データ転送料を別枠で見る。
- Billing alarm と AWS Budgets を必ず設定し、月 3,000 円、4,500 円、5,000 円で通知する。

### 無料枠重視案

| 領域 | サービス | コスト方針 |
| --- | --- | --- |
| Frontend | Amplify Hosting | 小規模なら無料枠内を狙う |
| API | API Gateway HTTP API + Lambda | リクエストが少ない間は無料枠内を狙う |
| Presence | DynamoDB TTL | 無料枠内を狙う。`expiresAt` でアプリ側も期限判定する |
| Secret | AWS Systems Manager Parameter Store または Secrets Manager | 無料枠 / 少額に収める。Secrets Manager はシークレット数で課金される点に注意 |
| Logs | CloudWatch Logs | 保持期間を 7 日から 14 日程度に制限する |

この案は最も無料枠に寄せやすいが、Redis 互換ではない。実装側では `presence.Store` の DynamoDB 実装を追加し、環境変数で Redis / DynamoDB を切り替えられるようにする。

### 実装変更少なめ案

| 領域 | サービス | コスト方針 |
| --- | --- | --- |
| App | Lightsail instance | 小さいプランから開始する |
| Frontend | 同一インスタンス上の Next.js、または Amplify Hosting | まずは単純化を優先 |
| API | 同一インスタンス上の Go API | Docker Compose を本番向けに調整する |
| Presence | 同一インスタンス上の Redis | マネージドではないため永続化と再起動手順を確認する |
| TLS | Lightsail Load Balancer、CloudFront、または Caddy / nginx | 証明書更新を自動化する |

この案は現在の Redis 実装を活かしやすい。代わりに、OS 更新、Docker 更新、Redis のデータ保存、ログ整理、死活監視を自分で持つ。

### 避けたい構成

無料枠重視の MVP では次を避ける。

- NAT Gateway を常時稼働させる構成。
- ECS Fargate、ALB、ElastiCache、RDS を最初から全部使う構成。
- Multi-AZ のマネージド Redis を MVP 初期から使う構成。
- CloudWatch Logs を無期限保持する構成。

これらは本番運用としては自然な選択肢になり得るが、月 5,000 円以内に収めるには固定費が重くなりやすい。

## 環境変数

### バックエンド

| 変数 | 本番での扱い |
| --- | --- |
| `API_ADDR` | 基盤指定の port に合わせる。`PORT` が渡される環境では変換が必要 |
| `JWT_SECRET` | secret として設定する。開発値 `change-me` は使わない |
| `DISCORD_CLIENT_ID` | Discord Developer Portal の本番 app の値 |
| `DISCORD_CLIENT_SECRET` | secret として設定する |
| `DISCORD_REDIRECT_URI` | 本番 frontend の `https://.../auth/callback` |
| `REDIS_ADDR` | Redis の private endpoint または接続 URL |
| `REDIS_PASSWORD` | Redis 側で必要な場合だけ secret として設定 |
| `REDIS_DB` | マネージド Redis で DB 番号が制限される場合は `0` |
| `SCHOOL_CLOSE_TIME` | MVP 要件通り `20:30` |

### フロントエンド

| 変数 | 本番での扱い |
| --- | --- |
| `NEXT_PUBLIC_API_BASE_URL` | ブラウザから到達できる API の HTTPS URL |
| `NEXT_PUBLIC_DISCORD_CLIENT_ID` | Discord Developer Portal の本番 app の client ID |
| `NEXT_PUBLIC_DISCORD_REDIRECT_URI` | 本番 frontend の `https://.../auth/callback` |

## 本番化前に必要な実装

### Dockerfile

現状の `docker-compose.yml` は開発用で、`go run` と `npm run dev` を使っている。本番では次を追加する。

- Go API 用の multi-stage `Dockerfile`
- Next.js 用の build / start 手順、または Vercel / Railway / Render の framework build 設定
- `.dockerignore` の本番向け確認

### Port 設定

多くの PaaS は `PORT` 環境変数で listen port を渡す。現在の Go API は `API_ADDR` を読むため、次のどちらかに統一する。

- デプロイ環境で `API_ADDR=:$PORT` 相当を設定できる場合は設定で吸収する。
- アプリ側で `PORT` があれば `API_ADDR` より優先して `:<PORT>` を使う。

後者のほうが PaaS 移植性は高い。

### Healthcheck

既存の `/healthz` を本番の healthcheck に使う。

- API: `GET /healthz`
- フロントエンド: `/`、`/manifest.webmanifest`、`/icon.svg`
- Redis: API 起動時の接続確認、またはデプロイ基盤側の datastore status

## 運用

### ログ

MVP ではアプリログを標準出力に出し、ホスティング基盤のログビューアで確認する。

ログに出してよいもの:

- HTTP method
- path
- status code
- latency
- error code

ログに出さないもの:

- JWT
- Discord access token
- `DISCORD_CLIENT_SECRET`
- `JWT_SECRET`
- Redis password

### 監視

MVP では最初から重い監視基盤を入れない。最低限、次を確認できる状態にする。

- API が起動している
- Redis に接続できる
- 5xx が増えていない
- Discord OAuth2 設定ミスをエラーコードで判別できる

### バックアップ

Presence は TTL 付きの一時データであり、MVP ではバックアップ対象にしない。

ただし将来、プロフィール、参加履歴、通知設定などの永続データを持つ場合は、PostgreSQL などの永続 DB を追加し、Redis はキャッシュまたは一時状態に限定する。

## 決定事項

- MVP のデプロイ先は Railway に決定する。
- Railway project には `frontend`、`backend`、`redis` の 3 サービスを置く。
- Redis は Railway project 内の private network で backend から接続する。
- AWS は現時点では採用しない。無料枠重視の代替案として Amplify + API Gateway + Lambda + DynamoDB TTL を記録する。
- ローカル開発は既存の Docker Compose を継続する。
- 本番投入前に Dockerfile と port 設定の方針を決める。
- 本番と検証で Discord OAuth2 app、Redis、JWT secret を分ける。

## 未決事項

- 本番ドメインを取得するか、ホスティング基盤の標準ドメインで開始するか。
- Railway で `frontend` と `backend` を別サービスにするか、将来 1 コンテナ構成へ寄せるか。
- Go API の `PORT` 対応を実装するタイミング。
- 本番用 Dockerfile を `deployments/` に置くか、各アプリのルートに置くか。
- Railway 本番と Railway staging を同一 project 内の環境で分けるか、別 project に分けるか。

## 参考

- Railway variables: https://docs.railway.com/variables
- Railway private networking: https://docs.railway.com/private-networking
- AWS Amplify pricing: https://aws.amazon.com/amplify/pricing/
- Amazon API Gateway pricing: https://aws.amazon.com/api-gateway/pricing/
- AWS Lambda pricing: https://aws.amazon.com/lambda/pricing/
- Amazon DynamoDB pricing: https://aws.amazon.com/dynamodb/pricing/
- Amazon Lightsail pricing: https://aws.amazon.com/lightsail/pricing/
- Amazon ElastiCache pricing: https://aws.amazon.com/elasticache/pricing/
- Render web services: https://render.com/docs/web-services/
- Render environment variables: https://render.com/docs/configure-environment-variables
- Render Key Value: https://render.com/docs/redis
- Fly.io Upstash Redis: https://fly.io/docs/upstash/redis/
- Vercel Next.js: https://vercel.com/docs/concepts/next.js/overview
- Vercel environment variables: https://vercel.com/docs/environment-variables
