<template>
  <view class="detail-container" v-if="trader">
    <view class="trader-header">
      <image :src="trader.avatar || '/static/images/default.png'" class="avatar" mode="aspectFill" />
      <view class="trader-info">
        <text class="trader-name">{{ trader.trader_name }}</text>
        <view class="trader-tags">
          <text class="tag">{{ trader.status === 'running' ? '运行中' : '已停止' }}</text>
          <text v-if="trader.is_curated" class="tag curated">精选</text>
        </view>
      </view>
    </view>
    <view class="core-stats">
      <view class="stat-card highlight">
        <text class="stat-label">总收益 (USDT)</text>
        <text :class="['stat-value', trader.total_pnl >= '0' ? 'up' : 'down']">{{ formatPnl(trader.total_pnl) }}</text>
      </view>
      <view class="stat-card">
        <text class="stat-label">收益率</text>
        <text :class="['stat-value', trader.total_roi >= '0' ? 'up' : 'down']">{{ formatPercent(trader.total_roi) }}</text>
      </view>
      <view class="stat-card">
        <text class="stat-label">粉丝数</text>
        <text class="stat-value">{{ formatLargeNumber(trader.follower_count) }}</text>
      </view>
      <view class="stat-card">
        <text class="stat-label">胜率</text>
        <text class="stat-value">{{ formatPercent(trader.win_rate) }}</text>
      </view>
    </view>
    <view class="detail-section">
      <view class="section-title">详细数据</view>
      <view class="data-grid">
        <view class="data-item">
          <text class="data-label">跟单收益</text>
          <text :class="['data-value', trader.follow_profit >= '0' ? 'up' : 'down']">{{ formatPnl(trader.follow_profit) }}</text>
        </view>
        <view class="data-item">
          <text class="data-label">最大回撤</text>
          <text class="data-value down">{{ formatPercent(trader.max_drawdown) }}</text>
        </view>
        <view class="data-item">
          <text class="data-label">持仓数</text>
          <text class="data-value">{{ trader.position_count }}</text>
        </view>
      </view>
    </view>
    <view class="detail-section">
      <button type="primary" @click="goToStats" class="chart-btn">查看收益曲线</button>
    </view>
  </view>
  <view v-else class="loading">加载中...</view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useTraderStore } from '@/store/modules/trader'
import { formatPnl, formatPercent, formatLargeNumber } from '@/utils/format'

const traderStore = useTraderStore()
const traderId = ref('')
const trader = computed(() => traderStore.selectedTrader)

const goToStats = () => uni.navigateTo({ url: `/pages/stats/chart?id=${traderId.value}` })
const loadDetail = async (id: string) => {
  traderId.value = id
  await traderStore.fetchTraderDetail(id)
}

onLoad((options) => {
  if (options?.id) loadDetail(options.id)
})
</script>

<style lang="scss" scoped>
.detail-container {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 40rpx;
}
.trader-header {
  background: linear-gradient(135deg, #00DC82, #00C775);
  padding: 40rpx;
  display: flex;
  align-items: center;
  .avatar {
    width: 100rpx;
    height: 100rpx;
    border-radius: 50%;
    border: 4rpx solid rgba(255,255,255,0.5);
  }
  .trader-info {
    margin-left: 20rpx;
    .trader-name {
      display: block;
      font-size: 36rpx;
      font-weight: bold;
      color: #FFF;
      margin-bottom: 10rpx;
    }
    .tag {
      font-size: 22rpx;
      padding: 4rpx 12rpx;
      border-radius: 20rpx;
      background: rgba(255,255,255,0.2);
      color: #FFF;
      margin-right: 10rpx;
      &.curated {
        background: linear-gradient(135deg, #FFD700, #FFA500);
      }
    }
  }
}
.core-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
  padding: 30rpx;
  margin-top: -40rpx;
}
.stat-card {
  background: #FFF;
  border-radius: 16rpx;
  padding: 24rpx;
  text-align: center;
  &.highlight {
    grid-column: span 2;
    background: linear-gradient(135deg, #00DC82, #00C775);
    .stat-label { color: rgba(255,255,255,0.9); }
    .stat-value { color: #FFF; }
  }
  .stat-label {
    display: block;
    font-size: 24rpx;
    color: #999;
    margin-bottom: 10rpx;
  }
  .stat-value {
    display: block;
    font-size: 36rpx;
    font-weight: bold;
    &.up { color: #00DC82; }
    &.down { color: #FF4D4D; }
  }
}
.detail-section {
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
.data-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16rpx;
}
.data-item {
  display: flex;
  justify-content: space-between;
  padding: 16rpx;
  background: #F8F8F8;
  border-radius: 12rpx;
  .data-label {
    font-size: 24rpx;
    color: #666;
  }
  .data-value {
    font-size: 28rpx;
    font-weight: bold;
    &.up { color: #00DC82; }
    &.down { color: #FF4D4D; }
  }
}
.chart-btn {
  width: 100%;
  height: 80rpx;
  font-size: 30rpx;
  background: linear-gradient(135deg, #00DC82, #00C775);
  border: none;
  color: #FFF;
}
.loading {
  text-align: center;
  padding: 100rpx;
  color: #999;
}
</style>
