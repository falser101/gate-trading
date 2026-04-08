<template>
  <view class="mine-container">
    <!-- 用户头部 -->
    <view class="user-header">
      <view class="avatar-wrapper">
        <image src="/static/images/default-avatar.png" class="avatar" mode="aspectFill" />
      </view>
      <view class="user-info">
        <text class="email">{{ userEmail }}</text>
        <text class="member-since">注册时间：{{ memberSince }}</text>
      </view>
    </view>

    <!-- 功能列表 -->
    <view class="menu-list">
      <view class="menu-item" @click="goToApiKey">
        <view class="menu-left">
          <text class="menu-icon">🔑</text>
          <text class="menu-label">API Key 配置</text>
        </view>
        <view class="menu-right">
          <text :class="['status', apiKeySet ? 'set' : 'unset']">
            {{ apiKeySet ? '已配置' : '未配置' }}
          </text>
          <text class="arrow">›</text>
        </view>
      </view>

      <view class="menu-item" @click="showAbout = true">
        <view class="menu-left">
          <text class="menu-icon">ℹ️</text>
          <text class="menu-label">关于</text>
        </view>
        <view class="menu-right">
          <text class="arrow">›</text>
        </view>
      </view>
    </view>

    <!-- 退出登录按钮 -->
    <view class="logout-section">
      <u-button
        type="info"
        size="large"
        @click="handleLogout"
        class="logout-btn"
      >
        退出登录
      </u-button>
    </view>

    <!-- 关于弹窗 -->
    <u-modal
      v-model:show="showAbout"
      title="关于"
      :content="'Gate Copy Trading v1.0.0\n\n专业的 Gate.io 跟单交易平台'"
      show-cancel-button
      @confirm="showAbout = false"
      @cancel="showAbout = false"
    ></u-modal>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/store/modules/auth'
import { formatDate } from '@/utils/format'

const authStore = useAuthStore()
const showAbout = ref(false)

const userEmail = computed(() => authStore.userEmail)
const apiKeySet = computed(() => authStore.userInfo?.api_key_set || false)
const memberSince = computed(() => {
  if (authStore.userInfo?.created_at) {
    return formatDate(authStore.userInfo.created_at, 'YYYY-MM-DD')
  }
  return '-'
})

// 跳转 API Key 配置页
const goToApiKey = () => {
  uni.navigateTo({
    url: '/pages/settings/api-key'
  })
}

// 退出登录
const handleLogout = () => {
  uni.showModal({
    title: '确认退出',
    content: '确定要退出登录吗？',
    success: (res) => {
      if (res.confirm) {
        authStore.logout()
        uni.showToast({
          title: '已退出登录',
          icon: 'success'
        })
      }
    }
  })
}

onMounted(() => {
  // 刷新用户信息
  authStore.refreshUserInfo()
})
</script>

<style lang="scss" scoped>
.mine-container {
  min-height: 100vh;
  background: #F5F5F5;
}

// 用户头部
.user-header {
  background: linear-gradient(135deg, #00DC82 0%, #00C775 100%);
  padding: 60rpx 40rpx 40rpx;
  display: flex;
  align-items: center;

  .avatar-wrapper {
    width: 120rpx;
    height: 120rpx;
    border-radius: 50%;
    background: #FFFFFF;
    padding: 6rpx;
    box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);

    .avatar {
      width: 100%;
      height: 100%;
      border-radius: 50%;
    }
  }

  .user-info {
    margin-left: 30rpx;

    .email {
      display: block;
      font-size: 36rpx;
      font-weight: bold;
      color: #FFFFFF;
      margin-bottom: 12rpx;
    }

    .member-since {
      display: block;
      font-size: 26rpx;
      color: rgba(255, 255, 255, 0.8);
    }
  }
}

// 功能列表
.menu-list {
  margin: 30rpx;
  background: #FFFFFF;
  border-radius: 20rpx;
  overflow: hidden;

  .menu-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 36rpx 30rpx;
    border-bottom: 1rpx solid #F5F5F5;

    &:last-child {
      border-bottom: none;
    }

    &:active {
      background: #F8F8F8;
    }

    .menu-left {
      display: flex;
      align-items: center;

      .menu-icon {
        font-size: 40rpx;
        margin-right: 24rpx;
      }

      .menu-label {
        font-size: 30rpx;
        color: #333333;
      }
    }

    .menu-right {
      display: flex;
      align-items: center;

      .status {
        font-size: 26rpx;
        padding: 8rpx 20rpx;
        border-radius: 20rpx;
        margin-right: 16rpx;

        &.set {
          background: #F0FDF5;
          color: #00DC82;
        }

        &.unset {
          background: #FFF5F5;
          color: #FF4D4D;
        }
      }

      .arrow {
        font-size: 40rpx;
        color: #CCCCCC;
      }
    }
  }
}

// 退出登录
.logout-section {
  margin: 60rpx 30rpx;

  .logout-btn {
    height: 88rpx;
    font-size: 30rpx;
    background: #FFFFFF;
    color: #FF4D4D;
    border: 2rpx solid #FF4D4D;
  }
}
</style>
