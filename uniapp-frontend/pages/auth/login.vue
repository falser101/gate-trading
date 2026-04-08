<template>
  <view class="login-container">
    <view class="login-header">
      <image src="/static/images/logo.png" class="logo" mode="aspectFit" />
      <text class="title">Gate Copy Trading</text>
      <text class="subtitle">专业跟单交易平台</text>
    </view>

    <view class="login-form">
      <u-form :model="form" ref="formRef" label-width="0">
        <u-form-item prop="email">
          <u-input
            v-model="form.email"
            placeholder="请输入邮箱"
            type="email"
            prefix-icon="envelope"
            :border="false"
            class="input-field"
          />
        </u-form-item>

        <u-form-item prop="password">
          <u-input
            v-model="form.password"
            placeholder="请输入密码"
            type="password"
            prefix-icon="lock"
            :border="false"
            class="input-field"
          />
        </u-form-item>
      </u-form>

      <view class="login-buttons">
        <u-button
          type="primary"
          size="large"
          :loading="loading"
          @click="handleLogin"
          class="login-btn"
        >
          登录
        </u-button>

        <u-button
          type="info"
          size="large"
          @click="goToRegister"
          class="register-btn"
        >
          注册账号
        </u-button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/store/modules/auth'

const authStore = useAuthStore()
const formRef = ref<any>(null)
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
      // 跳转到首页
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

  .logo {
    width: 120rpx;
    height: 120rpx;
    margin-bottom: 30rpx;
  }

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

  .input-field {
    padding: 24rpx 0;
    font-size: 32rpx;
  }
}

.login-buttons {
  margin-top: 60rpx;

  .login-btn {
    margin-bottom: 24rpx;
    background: linear-gradient(135deg, #00DC82 0%, #00C775 100%);
    border: none;
  }

  .register-btn {
    background: #FFFFFF;
    color: #00DC82;
    border: 2rpx solid #00DC82;
  }
}
</style>
