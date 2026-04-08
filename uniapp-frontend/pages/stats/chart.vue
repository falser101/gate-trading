<template>
  <view class="stats-container">
    <!-- 日期选择 -->
    <view class="date-filter">
      <scroll-view scroll-x class="filter-scroll">
        <view
          v-for="item in dateOptions"
          :key="item.value"
          :class="['filter-item', currentDateRange === item.value ? 'active' : '']"
          @click="selectDateRange(item.value)"
        >
          {{ item.label }}
        </view>
      </scroll-view>
    </view>

    <!-- 收益曲线图 -->
    <view class="chart-section">
      <view class="section-title">收益走势</view>
      <view class="chart-wrapper">
        <canvas
          type="2d"
          id="pnlChart"
          class="chart-canvas"
        ></canvas>
      </view>
    </view>

    <!-- 统计数据 -->
    <view class="stats-section">
      <view class="section-title">统计摘要</view>
      <view class="stats-grid">
        <view class="stat-item">
          <text class="stat-label">最高收益</text>
          <text :class="['stat-value', 'up']">{{ maxPnl }}</text>
        </view>
        <view class="stat-item">
          <text class="stat-label">最低收益</text>
          <text :class="['stat-value', 'down']">{{ minPnl }}</text>
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

    <!-- 每日数据列表 -->
    <view class="stats-section">
      <view class="section-title">每日明细</view>
      <view class="data-list">
        <view
          v-for="(item, index) in stats"
          :key="index"
          class="data-item"
        >
          <text class="date">{{ item.date }}</text>
          <text :class="['pnl', item.total_pnl >= 0 ? 'up' : 'down']">
            {{ formatPnl(item.total_pnl) }}
          </text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useTraderStore } from '@/store/modules/trader'

const traderStore = useTraderStore()
const traderId = ref('')

// 日期选项
const dateOptions = ref([
  { label: '7 天', value: '7' },
  { label: '15 天', value: '15' },
  { label: '30 天', value: '30' },
  { label: '90 天', value: '90' }
])
const currentDateRange = ref('30')

const stats = ref<any[]>([])

// 统计数据
const maxPnl = ref('0')
const minPnl = ref('0')
const avgPnl = ref('0')
const positiveDays = ref(0)

// 选择日期范围
const selectDateRange = (value: string) => {
  currentDateRange.value = value
  loadStats(value)
}

// 加载统计数据
const loadStats = async (days: string) => {
  const endDate = new Date()
  const startDate = new Date()
  startDate.setDate(startDate.getDate() - parseInt(days))

  const startDateStr = startDate.toISOString().split('T')[0]
  const endDateStr = endDate.toISOString().split('T')[0]

  const result = await traderStore.fetchTraderStats(traderId.value, startDateStr, endDateStr)
  if (result.success) {
    stats.value = traderStore.traderStats
    calculateStats()
    nextTick(() => {
      renderChart()
    })
  }
}

// 计算统计数据
const calculateStats = () => {
  if (stats.value.length === 0) return

  const pnls = stats.value.map(s => parseFloat(s.total_pnl) || 0)
  maxPnl.value = Math.max(...pnls).toFixed(2)
  minPnl.value = Math.min(...pnls).toFixed(2)
  avgPnl.value = (pnls.reduce((a, b) => a + b, 0) / pnls.length).toFixed(2)
  positiveDays.value = pnls.filter(p => p > 0).length
}

// 格式化盈亏
const formatPnl = (value: string) => {
  const num = parseFloat(value) || 0
  return num > 0 ? `+${num.toFixed(2)}` : num.toFixed(2)
}

// 渲染图表
const renderChart = () => {
  // 简单图表渲染逻辑
  // 实际项目中建议使用 ucharts 完整功能
  console.log('Render chart with data:', stats.value)
}

onLoad((options) => {
  if (options?.id) {
    traderId.value = options.id
    loadStats(currentDateRange.value)
  }
})
</script>

<style lang="scss" scoped>
.stats-container {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 40rpx;
}

// 日期筛选
.date-filter {
  background: #FFFFFF;
  padding: 20rpx 0;
  margin-bottom: 2rpx;

  .filter-scroll {
    white-space: nowrap;
    padding: 0 30rpx;

    .filter-item {
      display: inline-block;
      padding: 12rpx 30rpx;
      margin-right: 20rpx;
      font-size: 28rpx;
      color: #666666;
      border-radius: 30rpx;
      background: #F5F5F5;
      transition: all 0.3s;

      &.active {
        background: #00DC82;
        color: #FFFFFF;
      }
    }
  }
}

// 图表区域
.chart-section {
  background: #FFFFFF;
  margin: 20rpx;
  border-radius: 20rpx;
  padding: 30rpx;

  .section-title {
    font-size: 30rpx;
    font-weight: bold;
    color: #333333;
    margin-bottom: 24rpx;
  }

  .chart-wrapper {
    height: 400rpx;
    width: 100%;

    .chart-canvas {
      width: 100%;
      height: 100%;
    }
  }
}

// 统计摘要
.stats-section {
  background: #FFFFFF;
  margin: 20rpx;
  border-radius: 20rpx;
  padding: 30rpx;

  .section-title {
    font-size: 30rpx;
    font-weight: bold;
    color: #333333;
    margin-bottom: 24rpx;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 24rpx;

    .stat-item {
      background: #F8F8F8;
      padding: 24rpx;
      border-radius: 16rpx;
      text-align: center;

      .stat-label {
        display: block;
        font-size: 24rpx;
        color: #999999;
        margin-bottom: 12rpx;
      }

      .stat-value {
        display: block;
        font-size: 32rpx;
        font-weight: bold;

        &.up {
          color: #00DC82;
        }

        &.down {
          color: #FF4D4D;
        }
      }
    }
  }

  .data-list {
    .data-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 24rpx 0;
      border-bottom: 1rpx solid #F0F0F0;

      &:last-child {
        border-bottom: none;
      }

      .date {
        font-size: 28rpx;
        color: #666666;
      }

      .pnl {
        font-size: 30rpx;
        font-weight: bold;

        &.up {
          color: #00DC82;
        }

        &.down {
          color: #FF4D4D;
        }
      }
    }
  }
}
</style>
