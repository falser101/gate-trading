<template>
  <view class="stats-container">
    <scroll-view scroll-x class="date-filter">
      <view v-for="item in dateOptions" :key="item.value" :class="['filter-item', dateRange === item.value ? 'active' : '']" @click="setDateRange(item.value)">
        {{ item.label }}
      </view>
    </scroll-view>
    <view class="stats-section">
      <view class="section-title">统计摘要</view>
      <view class="stats-grid">
        <view class="stat-item">
          <text class="stat-label">最高收益</text>
          <text class="stat-value up">{{ maxPnl }}</text>
        </view>
        <view class="stat-item">
          <text class="stat-label">最低收益</text>
          <text class="stat-value down">{{ minPnl }}</text>
        </view>
        <view class="stat-item">
          <text class="stat-label">平均收益</text>
          <text :class="['stat-value', avgPnl >= 0 ? 'up' : 'down']">{{ avgPnl }}</text>
        </view>
        <view class="stat-item">
          <text class="stat-label">正收益天数</text>
          <text class="stat-value">{{ positiveDays }}/{{ stats.length }}</text>
        </view>
      </view>
    </view>
    <view class="stats-section">
      <view class="section-title">每日明细</view>
      <view class="data-list">
        <view v-for="(item, index) in stats" :key="index" class="data-item">
          <text class="date">{{ item.date }}</text>
          <text :class="['pnl', parseFloat(item.total_pnl) >= 0 ? 'up' : 'down']">{{ formatPnl(item.total_pnl) }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useTraderStore } from '@/store/modules/trader'

const traderStore = useTraderStore()
const traderId = ref('')
const dateRange = ref('30')
const stats = ref<any[]>([])
const maxPnl = ref('0')
const minPnl = ref('0')
const avgPnl = ref('0')
const positiveDays = ref(0)

const dateOptions = ref([
  { label: '7 天', value: '7' },
  { label: '15 天', value: '15' },
  { label: '30 天', value: '30' },
  { label: '90 天', value: '90' }
])

const formatPnl = (v: string) => {
  const num = parseFloat(v) || 0
  return num > 0 ? `+${num.toFixed(2)}` : num.toFixed(2)
}

const setDateRange = (v: string) => { dateRange.value = v; loadStats() }

const loadStats = async () => {
  const endDate = new Date().toISOString().split('T')[0]
  const startDate = new Date()
  startDate.setDate(startDate.getDate() - parseInt(dateRange.value))
  const result = await traderStore.fetchTraderStats(traderId.value, startDate.toISOString().split('T')[0], endDate)
  if (result.success) {
    stats.value = traderStore.traderStats
    const pnls = stats.value.map(s => parseFloat(s.total_pnl) || 0)
    maxPnl.value = Math.max(...pnls).toFixed(2)
    minPnl.value = Math.min(...pnls).toFixed(2)
    avgPnl.value = (pnls.reduce((a, b) => a + b, 0) / pnls.length).toFixed(2)
    positiveDays.value = pnls.filter(p => p > 0).length
  }
}

onLoad((options) => {
  if (options?.id) {
    traderId.value = options.id
    loadStats()
  }
})
</script>

<style lang="scss" scoped>
.stats-container {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 40rpx;
}
.date-filter {
  background: #FFF;
  padding: 20rpx;
  white-space: nowrap;
  margin-bottom: 2rpx;
}
.filter-item {
  display: inline-block;
  padding: 12rpx 30rpx;
  margin-right: 20rpx;
  font-size: 28rpx;
  color: #666;
  border-radius: 30rpx;
  background: #F5F5F5;
  &.active {
    background: #00DC82;
    color: #FFF;
  }
}
.stats-section {
  margin: 20rpx;
  background: #FFF;
  border-radius: 16rpx;
  padding: 24rpx;
  .section-title {
    font-size: 28rpx;
    font-weight: bold;
    margin-bottom: 20rpx;
  }
}
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16rpx;
}
.stat-item {
  background: #F8F8F8;
  padding: 20rpx;
  border-radius: 12rpx;
  text-align: center;
  .stat-label {
    display: block;
    font-size: 24rpx;
    color: #999;
    margin-bottom: 10rpx;
  }
  .stat-value {
    display: block;
    font-size: 28rpx;
    font-weight: bold;
    &.up { color: #00DC82; }
    &.down { color: #FF4D4D; }
  }
}
.data-list {
  .data-item {
    display: flex;
    justify-content: space-between;
    padding: 20rpx 0;
    border-bottom: 1rpx solid #F0F0F0;
    &:last-child { border-bottom: none; }
    .date {
      font-size: 28rpx;
      color: #666;
    }
    .pnl {
      font-size: 28rpx;
      font-weight: bold;
      &.up { color: #00DC82; }
      &.down { color: #FF4D4D; }
    }
  }
}
</style>
