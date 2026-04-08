<template>
  <view class="mine-container">
    <view class="user-header">
      <view class="avatar-wrapper">
        <image src="/static/images/default.png" class="avatar" mode="aspectFill" />
      </view>
      <view class="user-info">
        <text class="email">{{ userEmail }}</text>
        <text class="member-since">注册时间：{{ memberSince }}</text>
      </view>
    </view>
    <view class="menu-list">
      <view class="menu-item" @click="goToApiKey">
        <view class="menu-left">
          <text class="menu-icon">🔑</text>
          <text class="menu-label">API Key 配置</text>
        </view>
        <view class="menu-right">
          <text :class="['status', apiKeySet ? 'set' : 'unset']">{{ apiKeySet ? '已配置' : '未配置' }}</text>
          <text class="arrow">›</text>
        </view>
      </view>
    </view>
    <view class="logout-section">
      <button type="default" @click="handleLogout" class="logout-btn">退出登录</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/store/modules/auth'
import { formatDate } from '@/utils/format'

const authStore = useAuthStore()
const userEmail = computed(() => authStore.userEmail)
const apiKeySet = computed(() => authStore.userInfo?.api_key_set || false)
const memberSince = computed(() => authStore.userInfo?.created_at ? formatDate(authStore.userInfo.created_at) : '-')

const goToApiKey = () => uni.navigateTo({ url: '/pages/settings/api-key' })
const handleLogout = () => {
  uni.showModal({
    title: '确认退出',
    content: '确定要退出登录吗？',
    success: (res) => res.confirm && authStore.logout()
  })
}

onMounted(() => {
  // 只在已登录时获取用户信息
  if (authStore.isLoggedIn) {
    authStore.fetchUserInfo()
  }
})
</script>

<style lang="scss" scoped>
.mine-container {
  min-height: 100vh;
  background: #F5F5F5;
}
.user-header {
  background: linear-gradient(135deg, #00DC82, #00C775);
  padding: 60rpx 40rpx;
  display: flex;
  align-items: center;
  .avatar-wrapper {
    width: 100rpx;
    height: 100rpx;
    border-radius: 50%;
    background: #FFF;
    padding: 6rpx;
    .avatar {
      width: 100%;
      height: 100%;
      border-radius: 50%;
    }
  }
  .user-info {
    margin-left: 20rpx;
    .email {
      display: block;
      font-size: 32rpx;
      font-weight: bold;
      color: #FFF;
      margin-bottom: 10rpx;
    }
    .member-since {
      display: block;
      font-size: 24rpx;
      color: rgba(255,255,255,0.8);
    }
  }
}
.menu-list {
  margin: 20rpx;
  background: #FFF;
  border-radius: 16rpx;
  overflow: hidden;
}
.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  border-bottom: 1rpx solid #F5F5F5;
  &:last-child { border-bottom: none; }
  &:active { background: #F8F8F8; }
  .menu-left {
    display: flex;
    align-items: center;
    .menu-icon {
      font-size: 40rpx;
      margin-right: 20rpx;
    }
    .menu-label {
      font-size: 28rpx;
      color: #333;
    }
  }
  .menu-right {
    display: flex;
    align-items: center;
    .status {
      font-size: 24rpx;
      padding: 6rpx 16rpx;
      border-radius: 16rpx;
      margin-right: 10rpx;
      &.set { background: #F0FDF5; color: #00DC82; }
      &.unset { background: #FFF5F5; color: #FF4D4D; }
    }
    .arrow {
      font-size: 36rpx;
      color: #CCC;
    }
  }
}
.logout-section {
  margin: 40rpx 20rpx;
  .logout-btn {
    width: 100%;
    height: 80rpx;
    font-size: 28rpx;
    background: #FFF;
    color: #FF4D4D;
    border: 2rpx solid #FF4D4D;
  }
}
</style>
