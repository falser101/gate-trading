<template>
  <view class="detail-container" v-if="trader">
    <!-- 交易员头部 -->
    <view class="trader-header">
      <image :src="trader.avatar || '/static/images/default-avatar.png'" class="avatar" mode="aspectFill" />
      <view class="trader-info">
        <text class="trader-name">{{ trader.trader_name }}</text>
        <view class="trader-meta">
          <text class="tag">{{ trader.status === 'running' ? '运行中' : '已停止' }}</text>
          <text v-if="trader.is_curated" class="tag curated">精选</text>
          <text class="tag">{{ trader.exchange }}</text>
        </view>
      </view>
    </view>

    <!-- 核心数据 -->
    <view class="core-stats">
      <view class="stat-card highlight">
        <text class="stat-label">总收益 (USDT)</text>
        <text :class="['stat-value', getPnlClass(trader.total_pnl)]">
          {{ formatPnlValue(trader.total_pnl) }}
        </text>
      </view>
      <view class="stat-card">
        <text class="stat-label">收益率</text>
        <text :class="['stat-value', getPnlClass(trader.total_roi)]">
          {{ formatPercent(trader.total_roi) }}
        </text>
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

    <!-- 详细数据 -->
    <view class="detail-section">
      <view class="section-title">详细数据</view>
      <view class="data-grid">
        <view class="data-item">
          <text class="data-label">跟单收益</text>
          <text :class="['data-value', getPnlClass(trader.follow_profit)]">
            {{ formatPnlValue(trader.follow_profit) }}
          </text>
        </view>
        <view class="data-item">
          <text class="data-label">最大回撤</text>
          <text class="data-value down">{{ formatPercent(trader.max_drawdown) }}</text>
        </view>
        <view class="data-item">
          <text class="data-label">持仓数</text>
          <text class="data-value">{{ trader.position_count }}</text>
        </view>
        <view class="data-item">
          <text class="data-label">平均杠杆</text>
          <text class="data-value">{{ trader.avg_leverage || '-' }}</text>
        </view>
      </view>
    </view>

    <!-- 风格标签 -->
    <view class="detail-section" v-if="trader.style_labels">
      <view class="section-title">交易风格</view>
      <view class="tags-container">
        <text
          v-for="(label, index) in parseLabels(trader.style_labels)"
          :key="index"
          class="style-tag"
        >
          {{ label }}
        </text>
      </view>
    </view>

    <!-- 统计图表按钮 -->
    <view class="detail-section">
      <view class="section-title">历史表现</view>
      <u-button
        type="primary"
        @click="goToStats"
        class="chart-btn"
      >
        查看收益曲线
      </u-button>
    </view>

    <!-- 底部操作栏 -->
    <view class="bottom-bar">
      <view class="bar-info">
        <text class="label">最后更新</text>
        <text class="value">{{ formatTime(trader.last_synced_at) }}</text>
      </view>
    </view>
  </view>

  <!-- 加载状态 -->
  <view v-else-if="loading" class="loading-state">
    <u-loading mode="circle" size="60"></u-loading>
    <text class="loading-text">加载中...</text>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useTraderStore } from '@/store/modules/trader'
import { formatNumber, formatPercent, formatLargeNumber, formatRelativeTime } from '@/utils/format'

const traderStore = useTraderStore()
const traderId = ref('')

const trader = computed(() => traderStore.selectedTrader)
const loading = computed(() => traderStore.loading)

// 格式化盈亏值
const formatPnlValue = (value: string) => {
  const num = parseFloat(value)
  if (isNaN(num)) return '0'
  return num > 0 ? `+${formatNumber(num)}` : formatNumber(num)
}

// 获取盈亏样式类
const getPnlClass = (value: string) => {
  const num = parseFloat(value)
  if (isNaN(num)) return ''
  return num >= 0 ? 'up' : 'down'
}

// 解析标签
const parseLabels = (labelsStr: string) => {
  try {
    return JSON.parse(labelsStr)
  } catch {
    return []
  }
}

// 格式化时间
const formatTime = (timeStr: string) => {
  if (!timeStr) return '-'
  return formatRelativeTime(timeStr)
}

// 跳转统计页面
const goToStats = () => {
  if (traderId.value) {
    uni.navigateTo({
      url: `/pages/stats/chart?id=${traderId.value}`
    })
  }
}

// 加载详情
const loadDetail = async (id: string) => {
  traderId.value = id
  await traderStore.fetchTraderDetail(id)
}

onLoad((options) => {
  if (options?.id) {
    loadDetail(options.id)
  }
})
</script>

<style lang="scss" scoped>
.detail-container {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 140rpx;
}

// 交易员头部
.trader-header {
  background: linear-gradient(135deg, #00DC82 0%, #00C775 100%);
  padding: 40rpx;
  display: flex;
  align-items: center;

  .avatar {
    width: 120rpx;
    height: 120rpx;
    border-radius: 50%;
    border: 4rpx solid rgba(255, 255, 255, 0.5);
  }

  .trader-info {
    margin-left: 30rpx;

    .trader-name {
      display: block;
      font-size: 40rpx;
      font-weight: bold;
      color: #FFFFFF;
      margin-bottom: 16rpx;
    }

    .trader-meta {
      display: flex;
      gap: 12rpx;

      .tag {
        font-size: 22rpx;
        padding: 6rpx 16rpx;
        border-radius: 20rpx;
        background: rgba(255, 255, 255, 0.2);
        color: #FFFFFF;

        &.curated {
          background: linear-gradient(135deg, #FFD700, #FFA500);
        }
      }
    }
  }
}

// 核心数据
.core-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
  padding: 30rpx;
  margin-top: -40rpx;

  .stat-card {
    background: #FFFFFF;
    border-radius: 20rpx;
    padding: 30rpx;
    text-align: center;
    box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);

    &.highlight {
      grid-column: span 2;
      background: linear-gradient(135deg, #00DC82, #00C775);

      .stat-label {
        color: rgba(255, 255, 255, 0.9);
      }

      .stat-value {
        color: #FFFFFF;
      }
    }

    .stat-label {
      display: block;
      font-size: 24rpx;
      color: #999999;
      margin-bottom: 12rpx;
    }

    .stat-value {
      display: block;
      font-size: 40rpx;
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

// 详细数据
.detail-section {
  margin: 30rpx;
  background: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx;

  .section-title {
    font-size: 30rpx;
    font-weight: bold;
    color: #333333;
    margin-bottom: 24rpx;
  }

  .data-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 24rpx;

    .data-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20rpx;
      background: #F8F8F8;
      border-radius: 16rpx;

      .data-label {
        font-size: 26rpx;
        color: #666666;
      }

      .data-value {
        font-size: 30rpx;
        font-weight: bold;
        color: #333333;

        &.up {
          color: #00DC82;
        }

        &.down {
          color: #FF4D4D;
        }
      }
    }
  }

  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 16rpx;

    .style-tag {
      font-size: 24rpx;
      padding: 10rpx 24rpx;
      border-radius: 30rpx;
      background: #F0FDF5;
      color: #00DC82;
    }
  }

  .chart-btn {
    width: 100%;
    height: 88rpx;
    font-size: 30rpx;
    background: linear-gradient(135deg, #00DC82, #00C775);
    border: none;
  }
}

// 底部操作栏
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 120rpx;
  background: #FFFFFF;
  border-top: 1rpx solid #EEEEEE;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 30rpx;
  padding-bottom: calc(120rpx - env(safe-area-inset-bottom));

  .bar-info {
    .label {
      display: block;
      font-size: 24rpx;
      color: #999999;
    }

    .value {
      font-size: 26rpx;
      color: #666666;
    }
  }
}

// 加载状态
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;

  .loading-text {
    margin-top: 30rpx;
    font-size: 28rpx;
    color: #999999;
  }
}
</style>
