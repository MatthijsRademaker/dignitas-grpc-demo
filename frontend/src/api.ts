// Shapes returned by the BFF (bff/Dtos.cs). The browser never sees protobuf.

export interface Lot {
  id: string
  title: string
  description: string
  emoji: string
  startingPrice: number
  minIncrement: number
}

export interface Bid {
  bidder: string
  amount: number
  placedAt: string
  /** How old the bid was when the server sent this: how late we are seeing it. */
  ageMs: number
}

/** What the whole room costs the auction server, over the last few seconds. */
export interface ServerLoad {
  pollsPerSecond: number
  emptyPollRatio: number
  pushesPerSecond: number
}

export type LotStatus = 'open' | 'sold' | 'unsold'

export interface Auction {
  lot: Lot
  status: LotStatus
  highestBid: Bid | null
  recentBids: Bid[]
  bidCount: number
  remainingMs: number
  watchers: number
  load: ServerLoad
}

export type AuctionEventKind = 'snapshot' | 'bid_placed' | 'lot_opened' | 'lot_closed' | 'watchers_changed' | 'load_changed'

export interface AuctionEvent {
  kind: AuctionEventKind
  auction: Auction
}

export interface Problem {
  status: number
  title: string
  detail?: string
}

export type Result<T> = { ok: true, value: T } | { ok: false, problem: Problem }

export function bidKey (bid: Bid): string {
  return bid.placedAt + bid.bidder
}

export function minimumBid (auction: Auction): number {
  return auction.highestBid ? auction.highestBid.amount + auction.lot.minIncrement : auction.lot.startingPrice
}

export async function getInfo (): Promise<{ auctionHost: string }> {
  const response = await fetch('/api/info')
  if (!response.ok) throw new Error(`GET /api/info failed with ${response.status}`)
  return response.json()
}

export async function getAuction (): Promise<Result<Auction>> {
  return toResult(await fetch('/api/auction'))
}

export async function placeBid (lotId: string, bidder: string, amount: number): Promise<Result<Auction>> {
  return toResult(await fetch('/api/bids', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ lotId, bidder, amount }),
  }))
}

async function toResult<T> (response: Response): Promise<Result<T>> {
  if (response.ok) return { ok: true, value: await response.json() }
  const isProblem = response.headers.get('content-type')?.includes('json')
  const problem: Problem = isProblem
    ? await response.json()
    : { status: response.status, title: response.statusText || 'Request failed' }
  return { ok: false, problem }
}
