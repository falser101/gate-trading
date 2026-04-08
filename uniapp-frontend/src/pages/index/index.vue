<template>
  <view class="index-container">
    <!-- 筛选栏 -->
    <view class="filter-bar">
      <scroll-view scroll-x class="filter-scroll">
        <view
          v-for="item in filterOptions"
          :key="item.value"
          :class="['filter-item', currentFilter === item.value ? 'active' : '']"
          @click="selectFilter(item.value)"
        >
          {{ item.label }}
        </view>
      </scroll-view>
    </view>

    <!-- 排序栏 -->
    <view class="sort-bar">
      <text class="sort-label">排序:</text>
      <scroll-view scroll-x class="sort-scroll">
        <view
          v-for="item in sortOptions"
          :key="item.value"
          :class="['sort-item', currentSort === item.value ? 'active' : '']"
          @click="selectSort(item.value)"
        >
          {{ item.label }}
          <text v-if="currentSort === item.value" class="sort-icon">
            {{ sortDirection === 'desc' ? '↓' : '↑' }}
          </text>
        </view>
      </scroll-view>
    </view>

    <!-- 交易员列表 -->
    <view class="trader-list">
      <view
        v-for="trader in traderList"
        :key="trader.trader_id"
        class="trader-card"
        @click="goToDetail(trader.trader_id)"
      >
        <view class="trader-header">
          <image :src="trader.avatar || '/static/images/default-avatar.png'" class="avatar" mode="aspectFill" />
          <view class="trader-info">
            <text class="trader-name">{{ trader.trader_name }}</text>
            <view class="trader-tags">
              <text class="tag">{{ trader.status === 'running' ? '运行中' : '已停止' }}</text>
              <text v-if="trader.is_curated" class="tag curated">精选</text>
            </view>
          </view>
          <view class="follower-count">
            <text class="num">{{ formatLargeNumber(trader.follower_count) }}</text>
            <text class="label">粉丝</text>
          </view>
        </view>

        <view class="trader-stats">
          <view class="stat-item">
            <text class="stat-label">总收益</text>
            <text :class="['stat-value', getPnlClass(trader.total_pnl)]">
              {{ formatPnlValue(trader.total_pnl) }}
            </text>
          </view>
          <view class="stat-item">
            <text class="stat-label">收益率</text>
            <text :class="['stat-value', getPnlClass(trader.total_roi)]">
              {{ formatPercent(trader.total_roi) }}
            </text>
          </view>
          <view class="stat-item">
            <text class="stat-label">胜率</text>
            <text class="stat-value">{{ formatPercent(trader.win_rate) }}</text>
          </view>
          <view class="stat-item">
            <text class="stat-label">最大回撤</text>
            <text class="stat-value down">{{ formatPercent(trader.max_drawdown) }}</text>
          </view>
        </view>

        <view class="trader-footer" v-if="trader.style_labels">
          <text
            v-for="(label, index) in parseLabels(trader.style_labels)"
            :key="index"
            class="style-tag"
          >
            {{ label }}
          </text>
        </view>
      </view>
    </view>

    <!-- 空状态 -->
    <view v-if="!loading && traderList.length === 0" class="empty-state">
      <text class="empty-text">暂无交易员数据</text>
    </view>

    <!-- 加载状态 -->
    <view v-if="loading" class="loading-state">
      <text>加载中...</text>
    </view>

    <!-- 加载更多 -->
    <view v-if="hasMore && !loading" class="load-more" @click="loadMore">
      <text>加载更多</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTraderStore } from '@/store/modules/trader'
import { formatNumber, formatPercent, formatLargeNumber } from '@/utils/format'

const traderStore = useTraderStore()

// 筛选选项
const filterOptions = ref([
  { label: '全部', value: 'all' },
  { label: '运行中', value: 'running' },
  { label: '精选', value: 'curated' }
])
const currentFilter = ref('all')

// 排序选项
const sortOptions = ref([
  { label: '收益', value: 'follow_profit' },
  { label: 'ROI', value: 'total_roi' },
  { label: '粉丝', value: 'follower_count' },
  { label: '胜率', value: 'win_rate' }
])
const currentSort = ref('follow_profit')
const sortDirection = ref('desc')

// 分页
const page = ref(1)
const pageSize = ref(20)
const hasMore = computed(() => traderStore.traderList.length < traderStore.total)

const traderList = computed(() => traderStore.traderList)
const loading = computed(() => traderStore.loading)

// 选择筛选
const selectFilter = (value: string) => {
  currentFilter.value = value
  page.value = 1
  loadTraderList()
}

// 选择排序
const selectSort = (value: string) => {
  if (currentSort.value === value) {
    sortDirection.value = sortDirection.value === 'desc' ? 'asc' : 'desc'
  } else {
    currentSort.value = value
    sortDirection.value = 'desc'
  }
  page.value = 1
  loadTraderList()
}

// 加载交易员列表
const loadTraderList = async () => {
  const params: any = {
    page: page.value,
    page_size: pageSize.value,
    order_by: currentSort.value,
    sort_by: sortDirection.value
  }

  if (currentFilter.value === 'running') {
    params.status = 'running'
  }

  await traderStore.fetchTraderList(params)
}

// 加载更多
const loadMore = () => {
  page.value++
  loadTraderList()
}

// 跳转详情
const goToDetail = (traderId: string) => {
  uni.navigateTo({
    url: `/pages/trader/detail?id=${traderId}`
  })
}

// 格式化盈亏值
const formatPnlValue = (value: string) => {
  const num = parseFloat(value)
  return num > 0 ? `+${formatNumber(num)}` : formatNumber(num)
}

// 获取盈亏样式类
const getPnlClass = (value: string) => {
  const num = parseFloat(value)
  return num >= 0 ? 'up' : 'down'
}

// 解析风格标签
const parseLabels = (labelsStr: string) => {
  try {
    return JSON.parse(labelsStr)
  } catch {
    return []
  }
}

onMounted(() => {
  loadTraderList()
})
</script>

<style lang="scss" scoped>
.index-container {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 120rpx;
}

.filter-bar {
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

      &.active {
        background: #00DC82;
        color: #FFFFFF;
      }
    }
  }
}

.sort-bar {
  display: flex;
  align-items: center;
  background: #FFFFFF;
  padding: 20rpx 30rpx;
  margin-bottom: 2rpx;

  .sort-label {
    font-size: 26rpx;
    color: #999999;
    margin-right: 20rpx;
  }

  .sort-scroll {
    flex: 1;
    white-space: nowrap;

    .sort-item {
      display: inline-flex;
      align-items: center;
      padding: 8rpx 20rpx;
      margin-right: 16rpx;
      font-size: 26rpx;
      color: #666666;
      border-radius: 20rpx;
      background: #F5F5F5;

      &.active {
        background: #00DC82;
        color: #FFFFFF;

        .sort-icon {
          color: #FFFFFF;
        }
      }

      .sort-icon {
        margin-left: 6rpx;
        color: #666666;
      }
    }
  }
}

.trader-list {
  padding: 20rpx;

  .trader-card {
    background: #FFFFFF;
    border-radius: 20rpx;
    padding: 30rpx;
    margin-bottom: 20rpx;
    box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.05);

    .trader-header {
      display: flex;
      align-items: center;
      margin-bottom: 30rpx;

      .avatar {
        width: 80rpx;
        height: 80rpx;
        border-radius: 50%;
        margin-right: 24rpx;
      }

      .trader-info {
        flex: 1;

        .trader-name {
          display: block;
          font-size: 32rpx;
          font-weight: bold;
          color: #333333;
          margin-bottom: 12rpx;
        }

        .trader-tags {
          display: flex;
          gap: 12rpx;

          .tag {
            font-size: 22rpx;
            padding: 6rpx 16rpx;
            border-radius: 20rpx;
            background: #F5F5F5;
            color: #666666;

            &.curated {
              background: linear-gradient(135deg, #FFD700, #FFA500);
              color: #FFFFFF;
            }
          }
        }
      }

      .follower-count {
        text-align: center;

        .num {
          display: block;
          font-size: 30rpx;
          font-weight: bold;
          color: #00DC82;
        }

        .label {
          font-size: 22rpx;
          color: #999999;
        }
      }
    }

    .trader-stats {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 20rpx;
      padding: 20rpx 0;
      border-top: 1rpx solid #F5F5F5;
      border-bottom: 1rpx solid #F5F5F5;

      .stat-item {
        text-align: center;

        .stat-label {
          display: block;
          font-size: 22rpx;
          color: #999999;
          margin-bottom: 10rpx;
        }

        .stat-value {
          display: block;
          font-size: 28rpx;
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

    .trader-footer {
      display: flex;
      flex-wrap: wrap;
      gap: 12rpx;
      margin-top: 20rpx;

      .style-tag {
        font-size: 22rpx;
        padding: 6rpx 16rpx;
        border-radius: 20rpx;
        background: #F0FDF5;
        color: #00DC82;
      }
    }
  }
}

.empty-state {
  text-align: center;
  padding: 120rpx 0;

  .empty-text {
    font-size: 28rpx;
    color: #999999;
  }
}

.loading-state {
  text-align: center;
  padding: 60rpx 0;

  text {
    font-size: 26rpx;
    color: #999999;
  }
}

.load-more {
  text-align: center;
  padding: 30rpx;
  font-size: 28rpx;
  color: #00DC82;
}
</style>
