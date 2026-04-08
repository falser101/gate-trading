<template>
  <view class="register-container">
    <view class="register-header">
      <text class="title">注册账号</text>
      <text class="subtitle">创建您的 Gate Copy Trading 账户</text>
    </view>

    <view class="register-form">
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
            placeholder="请输入密码（至少 6 位）"
            type="password"
            prefix-icon="lock"
            :border="false"
            class="input-field"
          />
        </u-form-item>

        <u-form-item prop="confirmPassword">
          <u-input
            v-model="form.confirmPassword"
            placeholder="请确认密码"
            type="password"
            prefix-icon="lock"
            :border="false"
            class="input-field"
          />
        </u-form-item>
      </u-form>

      <view class="register-buttons">
        <u-button
          type="primary"
          size="large"
          :loading="loading"
          @click="handleRegister"
          class="register-btn"
        >
          立即注册
        </u-button>

        <u-button
          type="info"
          size="large"
          @click="goToLogin"
          class="login-btn"
        >
          已有账号？去登录
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
  password: '',
  confirmPassword: ''
})

const handleRegister = async () => {
  // 验证表单
  if (!form.email || !form.password || !form.confirmPassword) {
    uni.showToast({
      title: '请填写完整信息',
      icon: 'none'
    })
    return
  }

  if (form.password.length < 6) {
    uni.showToast({
      title: '密码至少 6 位',
      icon: 'none'
    })
    return
  }

  if (form.password !== form.confirmPassword) {
    uni.showToast({
      title: '两次密码不一致',
      icon: 'none'
    })
    return
  }

  loading.value = true

  try {
    const result = await authStore.register(form.email, form.password)
    if (result.success) {
      uni.showToast({
        title: '注册成功',
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
        title: result.message || '注册失败',
        icon: 'none'
      })
    }
  } catch (error) {
    uni.showToast({
      title: '注册失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  uni.redirectTo({
    url: '/pages/auth/login'
  })
}
</script>

<style lang="scss" scoped>
.register-container {
  min-height: 100vh;
  padding: 80rpx 40rpx;
  background: linear-gradient(135deg, #00DC82 0%, #00C775 100%);
}

.register-header {
  text-align: center;
  margin-bottom: 60rpx;

  .title {
    display: block;
    font-size: 44rpx;
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

.register-form {
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 40rpx;
  box-shadow: 0 10rpx 40rpx rgba(0, 0, 0, 0.1);

  .input-field {
    padding: 24rpx 0;
    font-size: 32rpx;
  }
}

.register-buttons {
  margin-top: 60rpx;

  .register-btn {
    margin-bottom: 24rpx;
    background: linear-gradient(135deg, #00DC82 0%, #00C775 100%);
    border: none;
  }

  .login-btn {
    background: #FFFFFF;
    color: #00DC82;
    border: 2rpx solid #00DC82;
  }
}
</style>
