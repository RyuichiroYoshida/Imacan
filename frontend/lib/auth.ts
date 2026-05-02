export const tokenStorageKey = "imacan.accessToken";
export const tokenExpiryStorageKey = "imacan.tokenExpiresAt";
export const oauthStateStorageKey = "imacan.oauthState";

export function saveToken(accessToken: string, expiresIn: number) {
  const expiresAt = Date.now() + expiresIn * 1000;
  localStorage.setItem(tokenStorageKey, accessToken);
  localStorage.setItem(tokenExpiryStorageKey, String(expiresAt));
}

export function getToken() {
  const token = localStorage.getItem(tokenStorageKey);
  const expiresAt = Number(localStorage.getItem(tokenExpiryStorageKey));
  if (!token || !expiresAt || Date.now() >= expiresAt) {
    clearToken();
    return null;
  }
  return token;
}

export function clearToken() {
  localStorage.removeItem(tokenStorageKey);
  localStorage.removeItem(tokenExpiryStorageKey);
}

export function createOAuthState() {
  const bytes = new Uint8Array(16);
  crypto.getRandomValues(bytes);
  const state = Array.from(bytes, (byte) => byte.toString(16).padStart(2, "0")).join("");
  localStorage.setItem(oauthStateStorageKey, state);
  return state;
}

export function consumeOAuthState() {
  const state = localStorage.getItem(oauthStateStorageKey);
  localStorage.removeItem(oauthStateStorageKey);
  return state;
}

export function buildDiscordAuthorizeUrl() {
  const clientId = process.env.NEXT_PUBLIC_DISCORD_CLIENT_ID;
  if (!clientId) {
    throw new Error("NEXT_PUBLIC_DISCORD_CLIENT_ID is not configured.");
  }

  const redirectUri =
    process.env.NEXT_PUBLIC_DISCORD_REDIRECT_URI ?? `${window.location.origin}/auth/callback`;
  const params = new URLSearchParams({
    client_id: clientId,
    redirect_uri: redirectUri,
    response_type: "code",
    scope: "identify",
    state: createOAuthState()
  });

  return `https://discord.com/api/oauth2/authorize?${params.toString()}`;
}

export function getRedirectUri() {
  return process.env.NEXT_PUBLIC_DISCORD_REDIRECT_URI ?? `${window.location.origin}/auth/callback`;
}
