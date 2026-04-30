// AIP-122 リソース名の Format/Parse (TS 側)。バックエンドの
// internal/platform/resourcename と1:1で合わせる。

export function formatPiano(id: string): string {
  return `pianos/${id}`;
}

export function parsePiano(name: string): string {
  const parts = name.split("/");
  if (parts.length !== 2 || parts[0] !== "pianos" || !parts[1]) {
    throw new Error(`invalid piano resource name: ${name}`);
  }
  return parts[1];
}

export function formatUser(customId: string): string {
  return `users/${customId}`;
}

export function parseUser(name: string): string {
  const parts = name.split("/");
  if (parts.length !== 2 || parts[0] !== "users" || !parts[1]) {
    throw new Error(`invalid user resource name: ${name}`);
  }
  return parts[1];
}
