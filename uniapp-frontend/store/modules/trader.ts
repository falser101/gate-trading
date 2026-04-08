import { defineStore } from 'pinia'
import { getTraderList, getTraderDetail, getTraderStats } from '@/api/trader'
import type { Trader, TraderDailyStats, TraderListParams } from '@/types'

interface TraderState {
  traderList: Trader[]
  total: number
  currentPage: number
  loading: boolean
  selectedTrader: Trader | null
  traderStats: TraderDailyStats[]
}

export const useTraderStore = defineStore('trader', {
  state: (): TraderState => ({
    traderList: [],
    total: 0,
    currentPage: 1,
    loading: false,
    selectedTrader: null,
    traderStats: []
  }),

  getters: {
    traderCount: (state) => state.total
  },

  actions: {
    // 获取交易员列表
    async fetchTraderList(params: TraderListParams = {}) {
      this.loading = true
      try {
        const res = await getTraderList(params)
        const data = res.data as any
        this.traderList = data.data || data.data?.data || []
        this.total = data.total || data.data?.total || 0
        this.currentPage = params.page || 1
        return { success: true }
      } catch (error: any) {
        return {
          success: false,
          message: error.message || '获取交易员列表失败'
        }
      } finally {
        this.loading = false
      }
    },

    // 获取交易员详情
    async fetchTraderDetail(traderId: string) {
      try {
        const res = await getTraderDetail(traderId)
        this.selectedTrader = res.data as Trader
        return { success: true }
      } catch (error: any) {
        return {
          success: false,
          message: error.message || '获取交易员详情失败'
        }
      }
    },

    // 获取交易员历史统计
    async fetchTraderStats(traderId: string, startDate: string, endDate: string) {
      try {
        const res = await getTraderStats(traderId, startDate, endDate)
        this.traderStats = res.data as TraderDailyStats[]
        return { success: true }
      } catch (error: any) {
        return {
          success: false,
          message: error.message || '获取历史统计失败'
        }
      }
    },

    // 清空选中交易员
    clearSelectedTrader() {
      this.selectedTrader = null
      this.traderStats = []
    }
  }
})
