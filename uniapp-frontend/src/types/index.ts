export interface User {
  id: number
  email: string
  api_key_set: boolean
  created_at: string
}

export interface Trader {
  id: number
  trader_id: string
  trader_name: string
  avatar: string
  exchange: string
  status: string
  total_pnl: string
  total_roi: string
  follow_profit: string
  win_rate: string
  follower_count: number
  position_count: number
  max_drawdown: string
  is_curated: boolean
  style_labels: string
  last_synced_at: string
}

export interface TraderDailyStats {
  id: number
  trader_id: string
  date: string
  total_pnl: string
  total_roi: string
  follower_count: number
}
