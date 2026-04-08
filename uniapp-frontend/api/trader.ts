import { get } from './request'

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

export interface TraderListParams {
  page?: number
  page_size?: number
  order_by?: string
  sort_by?: string
  cycle?: string
  status?: string
  search_key?: string
}

export interface TraderListResponse {
  data: Trader[]
  total: number
  page: number
  page_size: number
}

// 获取交易员列表
export const getTraderList = (params: TraderListParams) => {
  return get<TraderListResponse>('/copytrading/traders', params)
}

// 获取交易员详情
export const getTraderDetail = (traderId: string) => {
  return get<Trader>(`/copytrading/traders/${traderId}`)
}

// 获取交易员历史统计
export const getTraderStats = (traderId: string, startDate: string, endDate: string) => {
  return get<TraderDailyStats[]>(`/copytrading/traders/${traderId}/stats`, {
    start_date: startDate,
    end_date: endDate
  })
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
