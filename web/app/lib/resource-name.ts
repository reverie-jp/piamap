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

export function formatPianoPost(pianoId: string, postId: string): string {
  return `pianos/${pianoId}/posts/${postId}`;
}

export function parsePianoPost(name: string): { pianoId: string; postId: string } {
  const parts = name.split("/");
  if (
    parts.length !== 4 ||
    parts[0] !== "pianos" ||
    parts[2] !== "posts" ||
    !parts[1] ||
    !parts[3]
  ) {
    throw new Error(`invalid piano post resource name: ${name}`);
  }
  return { pianoId: parts[1], postId: parts[3] };
}

export function formatPianoPostComment(pianoId: string, postId: string, commentId: string): string {
  return `pianos/${pianoId}/posts/${postId}/comments/${commentId}`;
}

export function parsePianoPostComment(
  name: string,
): { pianoId: string; postId: string; commentId: string } {
  const parts = name.split("/");
  if (
    parts.length !== 6 ||
    parts[0] !== "pianos" ||
    parts[2] !== "posts" ||
    parts[4] !== "comments" ||
    !parts[1] ||
    !parts[3] ||
    !parts[5]
  ) {
    throw new Error(`invalid piano post comment resource name: ${name}`);
  }
  return { pianoId: parts[1], postId: parts[3], commentId: parts[5] };
}
