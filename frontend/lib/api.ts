export type Activity = "CLASS" | "SELF_STUDY" | "OUT";

export type AuthTokenResponse = {
  accessToken: string;
  tokenType: "Bearer";
  expiresIn: number;
};

export type PresenceSummary = {
  total: number;
  class: number;
  selfStudy: number;
};

export type CurrentPresence = {
  active: boolean;
  activity?: Activity;
  updatedAt?: string;
  expiresAt?: string;
};

export type PresenceUpdateResponse = {
  activity: Activity;
  updatedAt: string;
  expiresAt?: string;
};

const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export async function exchangeDiscordCode(code: string, redirectUri: string) {
  return request<AuthTokenResponse>("/auth/discord/callback", {
    method: "POST",
    body: JSON.stringify({ code, redirectUri })
  });
}

export async function fetchSummary() {
  return request<PresenceSummary>("/presence/summary");
}

export async function fetchCurrentPresence(token: string) {
  return request<CurrentPresence>("/presence/me", {
    headers: authHeaders(token)
  });
}

export async function updatePresence(token: string, activity: Activity) {
  return request<PresenceUpdateResponse>("/presence", {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify({ activity })
  });
}

async function request<T>(path: string, init: RequestInit = {}) {
  const headers = new Headers(init.headers);
  if (init.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(`${apiBaseUrl}${path}`, {
    ...init,
    headers,
    cache: "no-store"
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed with ${response.status}`);
  }

  return (await response.json()) as T;
}

function authHeaders(token: string) {
  return {
    Authorization: `Bearer ${token}`
  };
}

export function formatActivity(activity?: Activity) {
  switch (activity) {
    case "CLASS":
      return "授業中";
    case "SELF_STUDY":
      return "自習中";
    case "OUT":
      return "帰宅";
    default:
      return "未設定";
  }
}

export function formatDateTime(value?: string) {
  if (!value) {
    return "";
  }
  return new Intl.DateTimeFormat("ja-JP", {
    month: "numeric",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  }).format(new Date(value));
}
