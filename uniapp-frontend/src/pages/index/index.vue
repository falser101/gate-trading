<template>
  <view class="index-container">
    <scroll-view scroll-x class="filter-bar">
      <view v-for="item in filters" :key="item.value" :class="['filter-item', filter === item.value ? 'active' : '']" @click="setFilter(item.value)">
        {{ item.label }}
      </view>
    </scroll-view>
    <scroll-view scroll-x class="sort-bar">
      <text class="sort-label">排序:</text>
      <view v-for="item in sorts" :key="item.value" :class="['sort-item', sort === item.value ? 'active' : '']" @click="setSort(item.value)">
        {{ item.label }}
      </view>
    </scroll-view>
    <view class="trader-list">
      <view v-for="trader in traderList" :key="trader.trader_id" class="trader-card" @click="goToDetail(trader.trader_id)">
        <view class="trader-header">
          <image :src="trader.avatar || '/static/images/default.png'" class="avatar" mode="aspectFill" />
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
            <text :class="['stat-value', trader.total_pnl >= '0' ? 'up' : 'down']">{{ formatPnl(trader.total_pnl) }}</text>
          </view>
          <view class="stat-item">
            <text class="stat-label">收益率</text>
            <text :class="['stat-value', trader.total_roi >= '0' ? 'up' : 'down']">{{ formatPercent(trader.total_roi) }}</text>
          </view>
          <view class="stat-item">
            <text class="stat-label">胜率</text>
            <text class="stat-value">{{ formatPercent(trader.win_rate) }}</text>
          </view>
          <view class="stat-item">
            <text class="stat-label">回撤</text>
            <text class="stat-value down">{{ formatPercent(trader.max_drawdown) }}</text>
          </view>
        </view>
      </view>
    </view>
    <view v-if="!loading && traderList.length === 0" class="empty">暂无交易员数据</view>
    <view v-if="loading" class="loading">加载中...</view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTraderStore } from '@/store/modules/trader'
import { formatPnl, formatPercent, formatLargeNumber } from '@/utils/format'

const traderStore = useTraderStore()
const filter = ref('all')
const sort = ref('follow_profit')
const sortDir = ref('desc')

const filters = ref([
  { label: '全部', value: 'all' },
  { label: '运行中', value: 'running' }
])

const sorts = ref([
  { label: '收益', value: 'follow_profit' },
  { label: 'ROI', value: 'total_roi' },
  { label: '粉丝', value: 'follower_count' }
])

const traderList = computed(() => traderStore.traderList)
const loading = computed(() => traderStore.loading)

const setFilter = (v: string) => { filter.value = v; loadList() }
const setSort = (v: string) => {
  if (sort.value === v) sortDir.value = sortDir.value === 'desc' ? 'asc' : 'desc'
  else { sort.value = v; sortDir.value = 'desc' }
  loadList()
}

const loadList = async () => {
  const params: any = { page: 1, page_size: 20, order_by: sort.value, sort_by: sortDir.value }
  if (filter.value === 'running') params.status = 'running'
  await traderStore.fetchTraderList(params)
}

const goToDetail = (id: string) => uni.navigateTo({ url: `/pages/trader/detail?id=${id}` })

onMounted(() => loadList())
</script>

<style lang="scss" scoped>
.index-container {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 120rpx;
}
.filter-bar, .sort-bar {
  background: #FFFFFF;
  padding: 20rpx;
  white-space: nowrap;
  margin-bottom: 2rpx;
}
.filter-item, .sort-item {
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
.sort-label {
  font-size: 26rpx;
  color: #999;
  margin-right: 10rpx;
}
.trader-list {
  padding: 20rpx;
}
.trader-card {
  background: #FFF;
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}
.trader-header {
  display: flex;
  align-items: center;
  margin-bottom: 20rpx;
}
.avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  margin-right: 20rpx;
}
.trader-info {
  flex: 1;
  .trader-name {
    display: block;
    font-size: 32rpx;
    font-weight: bold;
    margin-bottom: 10rpx;
  }
  .tag {
    font-size: 22rpx;
    padding: 4rpx 12rpx;
    border-radius: 20rpx;
    background: #F5F5F5;
    margin-right: 10rpx;
    &.curated {
      background: linear-gradient(135deg, #FFD700, #FFA500);
      color: #FFF;
    }
  }
}
.follower-count {
  text-align: center;
  .num {
    display: block;
    font-size: 28rpx;
    font-weight: bold;
    color: #00DC82;
  }
  .label {
    font-size: 22rpx;
    color: #999;
  }
}
.trader-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid #F5F5F5;
}
.stat-item {
  text-align: center;
  .stat-label {
    display: block;
    font-size: 22rpx;
    color: #999;
  }
  .stat-value {
    display: block;
    font-size: 26rpx;
    font-weight: bold;
    &.up { color: #00DC82; }
    &.down { color: #FF4D4D; }
  }
}
.empty, .loading {
  text-align: center;
  padding: 60rpx;
  color: #999;
}
</style>
