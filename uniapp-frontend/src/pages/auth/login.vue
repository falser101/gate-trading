<template>
  <view class="login-container">
    <view class="login-header">
      <text class="title">Gate Copy Trading</text>
      <text class="subtitle">专业跟单交易平台</text>
    </view>

    <view class="login-form">
      <view class="input-item">
        <input
          v-model="form.email"
          type="text"
          placeholder="请输入邮箱"
          class="input-field"
        />
      </view>

      <view class="input-item">
        <input
          v-model="form.password"
          type="password"
          placeholder="请输入密码"
          class="input-field"
        />
      </view>

      <view class="login-buttons">
        <button
          type="primary"
          :loading="loading"
          @click="handleLogin"
          class="login-btn"
        >
          登录
        </button>

        <button
          type="default"
          @click="goToRegister"
          class="register-btn"
        >
          注册账号
        </button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/store/modules/auth'

const authStore = useAuthStore()
const loading = ref(false)

const form = reactive({
  email: '',
  password: ''
})

const handleLogin = async () => {
  if (!form.email || !form.password) {
    uni.showToast({
      title: '请填写邮箱和密码',
      icon: 'none'
    })
    return
  }

  loading.value = true

  try {
    const result = await authStore.login(form.email, form.password)
    if (result.success) {
      uni.showToast({
        title: '登录成功',
        icon: 'success'
      })
      setTimeout(() => {
        uni.switchTab({
          url: '/pages/index/index'
        })
      }, 1000)
    } else {
      uni.showToast({
        title: result.message || '登录失败',
        icon: 'none'
      })
    }
  } catch (error) {
    uni.showToast({
      title: '登录失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

const goToRegister = () => {
  uni.navigateTo({
    url: '/pages/auth/register'
  })
}
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  padding: 80rpx 40rpx;
  background: linear-gradient(135deg, #00DC82 0%, #00C775 100%);
}

.login-header {
  text-align: center;
  margin-bottom: 80rpx;

  .title {
    display: block;
    font-size: 48rpx;
    font-weight: bold;
    color: #FFFFFF;
    margin-bottom: 16rpx;
  }

  .subtitle {
    display: block;
    font-size: 28rpx;
    color: rgba(255, 255, 255, 0.8);
  }
}

.login-form {
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 40rpx;
  box-shadow: 0 10rpx 40rpx rgba(0, 0, 0, 0.1);

  .input-item {
    margin-bottom: 24rpx;

    .input-field {
      padding: 24rpx;
      font-size: 32rpx;
      background: #F5F5F5;
      border-radius: 12rpx;
    }
  }
}

.login-buttons {
  margin-top: 40rpx;

  .login-btn {
    margin-bottom: 20rpx;
    background: linear-gradient(135deg, #00DC82 0%, #00C775 100%);
    border: none;
    color: #FFFFFF;
  }

  .register-btn {
    background: #FFFFFF;
    color: #00DC82;
    border: 2rpx solid #00DC82;
  }
}
</style>
