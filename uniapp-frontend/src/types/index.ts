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
  cycle: string
  total_pnl: string
  total_roi: string
  follow_profit: string
  follow_roi: string
  win_rate: string
  follower_count: number
  position_count: number
  max_drawdown: string
  avg_leverage: string
  is_curated: boolean
  is_private: boolean
  style_labels: string
  last_synced_at: string
  created_at: string
  updated_at: string
}

export interface TraderDailyStats {
  id: number
  trader_id: string
  date: string
  total_pnl: string
  total_roi: string
  follow_profit: string
  follower_count: number
  created_at: string
}

export interface ApiKeyInfo {
  api_key: string
  api_secret: string
}
