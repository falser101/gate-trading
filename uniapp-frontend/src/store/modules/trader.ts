import { defineStore } from 'pinia'
import { getTraderList, getTraderDetail, getTraderStats } from '@/api/trader'
import type { Trader, TraderDailyStats, TraderListParams } from '@/types'

interface TraderState {
  traderList: Trader[]
  total: number
  loading: boolean
  selectedTrader: Trader | null
  traderStats: TraderDailyStats[]
}

export const useTraderStore = defineStore('trader', {
  state: (): TraderState => ({
    traderList: [],
    total: 0,
    loading: false,
    selectedTrader: null,
    traderStats: []
  }),

  actions: {
    async fetchTraderList(params: TraderListParams = {}) {
      this.loading = true
      try {
        const res = await getTraderList(params)
        const data = res.data as any
        this.traderList = data.data || []
        this.total = data.total || 0
        return { success: true }
      } catch (error: any) {
        return { success: false, message: error.message || '获取列表失败' }
      } finally {
        this.loading = false
      }
    },

    async fetchTraderDetail(traderId: string) {
      try {
        const res = await getTraderDetail(traderId)
        this.selectedTrader = res.data as Trader
        return { success: true }
      } catch (error: any) {
        return { success: false, message: error.message || '获取详情失败' }
      }
    },

    async fetchTraderStats(traderId: string, startDate: string, endDate: string) {
      try {
        const res = await getTraderStats(traderId, startDate, endDate)
        this.traderStats = res.data as TraderDailyStats[]
        return { success: true }
      } catch (error: any) {
        return { success: false, message: error.message || '获取统计失败' }
      }
    },

    clearSelectedTrader() {
      this.selectedTrader = null
      this.traderStats = []
    }
  }
})
