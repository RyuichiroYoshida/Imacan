"use client";

import { Suspense, useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { ApiError, exchangeDiscordCode } from "@/lib/api";
import { consumeOAuthState, getRedirectUri, saveToken } from "@/lib/auth";

export default function AuthCallbackPage() {
  return (
    <Suspense
      fallback={<CallbackShell message="Discordログインを確認しています。" />}
    >
      <AuthCallbackContent />
    </Suspense>
  );
}

function AuthCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [message, setMessage] = useState("Discordログインを確認しています。");
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    const code = searchParams.get("code");
    const state = searchParams.get("state");
    const expectedState = consumeOAuthState();

    if (!code) {
      setFailed(true);
      setMessage("認可コードが見つかりませんでした。");
      return;
    }
    if (!state || !expectedState || state !== expectedState) {
      setFailed(true);
      setMessage(
        "ログイン状態を確認できませんでした。もう一度ログインしてください。",
      );
      return;
    }

    exchangeDiscordCode(code, getRedirectUri())
      .then((token) => {
        saveToken(token.accessToken, token.expiresIn);
        router.replace("/");
      })
      .catch((error) => {
        setFailed(true);
        if (error instanceof ApiError && error.code === "NETWORK_ERROR") {
          setMessage("APIに接続できませんでした。APIが起動しているか確認してください。");
          return;
        }
        setMessage("Discordログインに失敗しました。時間をおいてもう一度試してください。");
      });
  }, [router, searchParams]);

  return (
    <main className="callback-shell">
      <section className="callback-panel" aria-live="polite">
        <h1 className="section-title">
          {failed ? "ログインできませんでした" : "ログイン中"}
        </h1>
        <p className={failed ? "message error" : "message"}>{message}</p>
        {failed ? (
          <button
            className="secondary-button"
            type="button"
            onClick={() => router.replace("/")}
          >
            ホームに戻る
          </button>
        ) : null}
      </section>
    </main>
  );
}

function CallbackShell({ message }: { message: string }) {
  return (
    <main className="callback-shell">
      <section className="callback-panel" aria-live="polite">
        <h1 className="section-title">ログイン中</h1>
        <p className="message">{message}</p>
      </section>
    </main>
  );
}
