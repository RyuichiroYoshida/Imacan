import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;

@service(#{ title: "ProjectImacan API" })
@doc("学校内の在席状態と活動ステータスを共有するProjectImacanのMVP APIです。")
@server("http://localhost:8080", "Local development server")
namespace ProjectImacan;

@doc("ユーザーが共有する現在の活動ステータスです。")
enum Activity {
  @doc("授業中です。授業時間相当のTTLで自動失効します。")
  CLASS,

  @doc("自習中です。当日の閉館時刻まで有効です。")
  SELF_STUDY,

  @doc("退出状態です。Presenceは即時削除されます。")
  OUT,
}

@doc("APIエラーの共通レスポンス本文です。")
@example(
  #{
    code: "INVALID_REQUEST",
    message: "activityにはCLASS、SELF_STUDY、OUTのいずれかを指定してください。",
  },
  #{ title: "バリデーションエラー" }
)
model ErrorResponseBody {
  @doc("機械的に判定できるエラーコードです。")
  code: string;

  @doc("ユーザーまたは開発者向けのエラーメッセージです。")
  message: string;
}

@doc("リクエスト内容が不正な場合のレスポンスです。")
model BadRequestResponse {
  @statusCode statusCode: 400;
  @body body: ErrorResponseBody;
}

@doc("JWTが未指定、または無効な場合のレスポンスです。")
model UnauthorizedResponse {
  @statusCode statusCode: 401;
  @body body: ErrorResponseBody;
}

@doc("サーバー内部で想定外のエラーが発生した場合のレスポンスです。")
model InternalServerErrorResponse {
  @statusCode statusCode: 500;
  @body body: ErrorResponseBody;
}

@doc("Discord OAuth2の認可コードをバックエンドへ渡すリクエストです。")
@example(
  #{
    code: "discord-oauth2-authorization-code",
    redirectUri: "http://localhost:3000/auth/callback",
  },
  #{ title: "Discord OAuth2 callback" }
)
model DiscordCallbackRequest {
  @doc("Discord OAuth2から返却された認可コードです。")
  code: string;

  @doc("フロントエンドで使用したリダイレクトURIです。環境ごとの差異を吸収するため任意にしています。")
  redirectUri?: string;
}

@doc("認証成功時に返却するJWTアクセストークンです。")
@example(
  #{
    accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.example",
    tokenType: "Bearer",
    expiresIn: 3600,
  },
  #{ title: "JWT発行成功" }
)
model AuthTokenResponseBody {
  @doc("ProjectImacan APIで利用するJWTです。")
  accessToken: string;

  @doc("Authorizationヘッダーで利用するトークン種別です。")
  tokenType: "Bearer";

  @doc("アクセストークンの有効期間です。単位は秒です。")
  expiresIn: int32;
}

@doc("Discordログイン成功時のレスポンスです。")
model AuthTokenResponse {
  @statusCode statusCode: 200;
  @body body: AuthTokenResponseBody;
}

@doc("ユーザーのPresenceを更新するリクエストです。正確な位置情報は保存しません。")
@example(
  #{ activity: Activity.SELF_STUDY, lat: 35.68, lng: 139.76 },
  #{ title: "自習中に更新" }
)
@example(#{ activity: Activity.OUT }, #{ title: "帰宅して退出" })
model PresenceUpdateRequest {
  @doc("設定する活動ステータスです。")
  activity: Activity;

  @doc("任意の緯度です。在席判定にのみ利用し、正確な値は永続保存しません。")
  lat?: float64;

  @doc("任意の経度です。在席判定にのみ利用し、正確な値は永続保存しません。")
  lng?: float64;
}

@doc("Presence更新後の状態です。")
@example(
  #{
    activity: Activity.SELF_STUDY,
    updatedAt: utcDateTime.fromISO("2026-04-26T09:30:00Z"),
    expiresAt: utcDateTime.fromISO("2026-04-26T20:30:00Z"),
  },
  #{ title: "自習中の更新結果" }
)
model PresenceUpdateResponseBody {
  @doc("更新後の活動ステータスです。")
  activity: Activity;

  @doc("Presenceを更新した日時です。")
  updatedAt: utcDateTime;

  @doc("Presenceが自動失効する日時です。OUTの場合は返却されないことがあります。")
  expiresAt?: utcDateTime;
}

@doc("Presence更新成功時のレスポンスです。")
model PresenceUpdateResponse {
  @statusCode statusCode: 200;
  @body body: PresenceUpdateResponseBody;
}

@doc("現在有効なPresenceの集計結果です。")
@example(#{ total: 5, class: 3, selfStudy: 2 }, #{ title: "在席集計" })
model PresenceSummaryResponseBody {
  @doc("現在在席している合計人数です。")
  total: int32;

  @doc("授業中として共有している人数です。")
  class: int32;

  @doc("自習中として共有している人数です。")
  selfStudy: int32;
}

@doc("Presence集計取得成功時のレスポンスです。")
model PresenceSummaryResponse {
  @statusCode statusCode: 200;
  @body body: PresenceSummaryResponseBody;
}

@doc("ログイン中ユーザーの現在有効なPresenceです。")
@example(
  #{
    active: true,
    activity: Activity.SELF_STUDY,
    updatedAt: utcDateTime.fromISO("2026-04-26T09:30:00Z"),
    expiresAt: utcDateTime.fromISO("2026-04-26T20:30:00Z"),
  },
  #{ title: "現在のPresenceがある" }
)
@example(#{ active: false }, #{ title: "現在のPresenceがない" })
model CurrentPresenceResponseBody {
  @doc("現在有効なPresenceが存在するかどうかです。")
  active: boolean;

  @doc("現在の活動ステータスです。Presenceがない場合は返却されません。")
  activity?: Activity;

  @doc("Presenceを更新した日時です。Presenceがない場合は返却されません。")
  updatedAt?: utcDateTime;

  @doc("Presenceが自動失効する日時です。Presenceがない場合は返却されません。")
  expiresAt?: utcDateTime;
}

@doc("ログイン中ユーザーのPresence取得成功時のレスポンスです。")
model CurrentPresenceResponse {
  @statusCode statusCode: 200;
  @body body: CurrentPresenceResponseBody;
}

@route("/auth")
namespace Auth {
  @post
  @route("/discord/callback")
  @summary("Discordログインcallback")
  @doc("Discord OAuth2の認可コードを検証し、ProjectImacan APIで利用するJWTを発行します。")
  @opExample(
    #{
      parameters: #{ request: #{ code: "discord-oauth2-authorization-code" } },
      returnType: #{
        statusCode: 200,
        body: #{
          accessToken: "jwt.example",
          tokenType: "Bearer",
          expiresIn: 3600,
        },
      },
    },
    #{ title: "JWTを発行する" }
  )
  op discordCallback(
    @body request: DiscordCallbackRequest,
  ): AuthTokenResponse | BadRequestResponse | InternalServerErrorResponse;
}

@route("/presence")
namespace Presence {
  @post
  @useAuth(BearerAuth)
  @summary("Presence更新")
  @doc("ログイン中のユーザーの活動ステータスを更新します。OUTの場合はPresenceを即時削除します。")
  @opExample(
    #{
      parameters: #{
        request: #{ activity: Activity.SELF_STUDY, lat: 35.68, lng: 139.76 },
      },
      returnType: #{
        statusCode: 200,
        body: #{
          activity: Activity.SELF_STUDY,
          updatedAt: utcDateTime.fromISO("2026-04-26T09:30:00Z"),
          expiresAt: utcDateTime.fromISO("2026-04-26T20:30:00Z"),
        },
      },
    },
    #{ title: "自習中に更新する" }
  )
  op updatePresence(@body request: PresenceUpdateRequest):
    | PresenceUpdateResponse
    | BadRequestResponse
    | UnauthorizedResponse
    | InternalServerErrorResponse;

  @get
  @route("/me")
  @useAuth(BearerAuth)
  @summary("ログイン中ユーザーのPresence取得")
  @doc("ログイン中ユーザーの現在有効なPresenceを返します。Presenceがない、または期限切れの場合はactive=falseを返します。")
  @opExample(
    #{
      returnType: #{
        statusCode: 200,
        body: #{
          active: true,
          activity: Activity.SELF_STUDY,
          updatedAt: utcDateTime.fromISO("2026-04-26T09:30:00Z"),
          expiresAt: utcDateTime.fromISO("2026-04-26T20:30:00Z"),
        },
      },
    },
    #{ title: "現在のPresenceを取得する" }
  )
  op getCurrentPresence(
  ): CurrentPresenceResponse | UnauthorizedResponse | InternalServerErrorResponse;

  @get
  @route("/summary")
  @summary("Presence集計取得")
  @doc("現在有効なPresenceを集計し、在席人数と活動別人数を返します。")
  @opExample(
    #{
      returnType: #{
        statusCode: 200,
        body: #{ total: 5, class: 3, selfStudy: 2 },
      },
    },
    #{ title: "現在の在席状況を取得する" }
  )
  op getPresenceSummary(
  ): PresenceSummaryResponse | InternalServerErrorResponse;
}
