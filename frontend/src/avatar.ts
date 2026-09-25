// Initials avatars. The colour is a stable function of the name, so "Alice" looks the same on
// every screen in the room. All eight take white text at ≥ 5:1.
const AVATAR_COLOURS = ['#733FCB', '#0f766e', '#1d4ed8', '#be123c', '#92400e', '#15803d', '#475569', '#a21caf']

export function avatarColour (name: string): string {
  // FNV-1a: tiny, and spreads similar names (bot-ada, bot-ken) across the set.
  let hash = 0x811C9DC5
  for (const char of name) hash = Math.imul(hash ^ char.codePointAt(0)!, 0x01000193)
  return AVATAR_COLOURS[(hash >>> 0) % AVATAR_COLOURS.length]!
}

export function initials (name: string): string {
  const words = name.trim().split(/[\s_\-.]+/).filter(Boolean)
  // First and last word: "Ada Lovelace" → AL, "bot-ken" → BK, "Bob" → B.
  const picked = words.length > 1 ? [words[0]!, words.at(-1)!] : [words[0] ?? '?']
  return picked.map(word => [...word][0]).join('').toUpperCase()
}
