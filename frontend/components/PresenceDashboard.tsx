"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import {
  type Activity,
  type CurrentPresence,
  type PresenceSummary,
  fetchCurrentPresence,
  fetchSummary,
  formatActivity,
  formatDateTime,
  updatePresence
} from "@/lib/api";
import { buildDiscordAuthorizeUrl, clearToken, getToken } from "@/lib/auth";

type Message = {
  kind: "info" | "error";
  text: string;
};

const initialSummary: PresenceSummary = {
  total: 0,
  class: 0,
  selfStudy: 0
};

export function PresenceDashboard() {
  const [token, setToken] = useState<string | null>(null);
  const [summary, setSummary] = useState<PresenceSummary>(initialSummary);
  const [currentPresence, setCurrentPresence] = useState<CurrentPresence | null>(null);
  const [loading, setLoading] = useState(true);
  const [updating, setUpdating] = useState<Activity | null>(null);
  const [message, setMessage] = useState<Message | null>(null);

  const signedIn = Boolean(token);
  const hasRestoredPresence = Boolean(currentPresence?.active);

  const loadSummary = useCallback(async () => {
    const nextSummary = await fetchSummary();
    setSummary(nextSummary);
  }, []);

  const loadCurrentPresence = useCallback(async (accessToken: string) => {
    const presence = await fetchCurrentPresence(accessToken);
    setCurrentPresence(presence);
  }, []);

  const refresh = useCallback(
    async (accessToken: string | null) => {
      await loadSummary();
      if (accessToken) {
        await loadCurrentPresence(accessToken);
      }
    },
    [loadCurrentPresence, loadSummary]
  );

  useEffect(() => {
    const savedToken = getToken();
    setToken(savedToken);

    refresh(savedToken)
      .catch(() => {
        setMessage({
          kind: "error",
          text: "サーバーに接続できませんでした。APIが起動しているか確認してください。"
        });
      })
      .finally(() => setLoading(false));
  }, [refresh]);

  const statusLine = useMemo(() => {
    if (!signedIn) {
      return "Discordログイン後に状態を共有できます。";
    }
    if (currentPresence?.active) {
      return `${formatActivity(currentPresence.activity)}として共有中`;
    }
    return "現在共有中の状態はありません。";
  }, [currentPresence, signedIn]);

  async function handleLogin() {
    try {
      window.location.href = buildDiscordAuthorizeUrl();
    } catch (error) {
      setMessage({
        kind: "error",
        text: error instanceof Error ? error.message : "Discordログインを開始できませんでした。"
      });
    }
  }

  function handleLogout() {
    clearToken();
    setToken(null);
    setCurrentPresence(null);
    setMessage({ kind: "info", text: "ログアウトしました。" });
  }

  async function handleUpdate(activity: Activity) {
    if (!token) {
      setMessage({ kind: "error", text: "先にDiscordでログインしてください。" });
      return;
    }

    setUpdating(activity);
    setMessage(null);
    try {
      const result = await updatePresence(token, activity);
      await refresh(token);
      if (activity === "OUT") {
        setMessage({ kind: "info", text: "帰宅に更新しました。" });
      } else {
        setMessage({
          kind: "info",
          text: `${formatActivity(result.activity)}に更新しました。有効期限は${formatDateTime(result.expiresAt)}です。`
        });
      }
    } catch {
      setMessage({ kind: "error", text: "状態を更新できませんでした。" });
    } finally {
      setUpdating(null);
    }
  }

  return (
    <main className="app-shell">
      <div className="workspace">
        <aside className="side-panel" aria-label="アカウント">
          <div className="brand">
            <div className="brand-mark" aria-hidden="true">
              I
            </div>
            <div>
              <h1 className="brand-title">Imacan</h1>
              <p className="brand-copy">今、学校に誰かいるかを手動で軽く共有する場所です。</p>
            </div>
          </div>

          <p className="status-line">{loading ? "読み込み中です。" : statusLine}</p>

          <div className="session-actions">
            {signedIn ? (
              <button className="secondary-button" type="button" onClick={handleLogout}>
                ログアウト
              </button>
            ) : (
              <button className="login-button" type="button" onClick={handleLogin}>
                Discordでログイン
              </button>
            )}
          </div>
        </aside>

        <section className="main-panel" aria-label="在席状況">
          <div className="summary-grid">
            <Metric label="在席" value={summary.total} />
            <Metric label="授業中" value={summary.class} />
            <Metric label="自習中" value={summary.selfStudy} />
          </div>

          {signedIn && hasRestoredPresence ? (
            <section className="confirm-panel" aria-label="継続確認">
              <h2 className="confirm-title">前回の状態がまだ有効です</h2>
              <p className="muted">
                {formatActivity(currentPresence?.activity)}として共有中です。
                {currentPresence?.expiresAt
                  ? ` ${formatDateTime(currentPresence.expiresAt)}まで有効です。`
                  : ""}
              </p>
              <div className="confirm-actions">
                <button
                  className="primary-button"
                  type="button"
                  onClick={() => setMessage({ kind: "info", text: "現在の状態を継続します。" })}
                >
                  継続する
                </button>
                <button
                  className="danger-button"
                  type="button"
                  disabled={updating === "OUT"}
                  onClick={() => handleUpdate("OUT")}
                >
                  帰宅にする
                </button>
              </div>
            </section>
          ) : null}

          <section className="action-panel" aria-label="状態更新">
            <h2 className="section-title">状態を更新</h2>
            <p className="muted">1タップで現在の状態を共有します。MVPでは位置情報を使いません。</p>
            <div className="status-actions">
              <button
                className="status-button"
                type="button"
                disabled={!signedIn || updating === "CLASS"}
                onClick={() => handleUpdate("CLASS")}
              >
                授業中
              </button>
              <button
                className="status-button"
                type="button"
                disabled={!signedIn || updating === "SELF_STUDY"}
                onClick={() => handleUpdate("SELF_STUDY")}
              >
                自習中
              </button>
              <button
                className="danger-button"
                type="button"
                disabled={!signedIn || updating === "OUT"}
                onClick={() => handleUpdate("OUT")}
              >
                帰宅
              </button>
            </div>
          </section>

          {message ? <p className={`message ${message.kind === "error" ? "error" : ""}`}>{message.text}</p> : null}
        </section>
      </div>
    </main>
  );
}

function Metric({ label, value }: { label: string; value: number }) {
  return (
    <div className="metric">
      <p className="metric-label">{label}</p>
      <p className="metric-value">{value}</p>
    </div>
  );
}
