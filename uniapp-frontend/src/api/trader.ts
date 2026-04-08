import { get } from './request'

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

export interface TraderListParams {
  page?: number
  page_size?: number
  order_by?: string
  sort_by?: string
  cycle?: string
  status?: string
}

export interface TraderListResponse {
  data: Trader[]
  total: number
  page: number
  page_size: number
}

export interface TraderDailyStats {
  id: number
  trader_id: string
  date: string
  total_pnl: string
  total_roi: string
  follower_count: number
}

export const getTraderList = (params: TraderListParams) => {
  return get<TraderListResponse>('/copytrading/traders', params)
}

export const getTraderDetail = (traderId: string) => {
  return get<Trader>(`/copytrading/traders/${traderId}`)
}

export const getTraderStats = (traderId: string, startDate: string, endDate: string) => {
  return get<TraderDailyStats[]>(`/copytrading/traders/${traderId}/stats`, {
    start_date: startDate,
    end_date: endDate
  })
}
