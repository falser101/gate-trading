<template>
  <view class="register-container">
    <view class="register-header">
      <text class="title">注册账号</text>
      <text class="subtitle">创建您的 Gate Copy Trading 账户</text>
    </view>
    <view class="register-form">
      <input v-model="form.email" type="text" placeholder="请输入邮箱" class="input-field" />
      <input v-model="form.password" type="password" placeholder="请输入密码（至少 6 位）" class="input-field" />
      <input v-model="form.confirmPassword" type="password" placeholder="请确认密码" class="input-field" />
      <button type="primary" :loading="loading" @click="handleRegister" class="register-btn">立即注册</button>
      <button type="default" @click="goToLogin" class="login-btn">已有账号？去登录</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/store/modules/auth'

const authStore = useAuthStore()
const loading = ref(false)
const form = reactive({ email: '', password: '', confirmPassword: '' })

const handleRegister = async () => {
  if (!form.email || !form.password || !form.confirmPassword) {
    uni.showToast({ title: '请填写完整信息', icon: 'none' })
    return
  }
  if (form.password.length < 6) {
    uni.showToast({ title: '密码至少 6 位', icon: 'none' })
    return
  }
  if (form.password !== form.confirmPassword) {
    uni.showToast({ title: '两次密码不一致', icon: 'none' })
    return
  }
  loading.value = true
  const result = await authStore.register(form.email, form.password)
  if (result.success) {
    uni.showToast({ title: '注册成功', icon: 'success' })
    setTimeout(() => uni.switchTab({ url: '/pages/index/index' }), 1000)
  } else {
    uni.showToast({ title: result.message, icon: 'none' })
  }
  loading.value = false
}

const goToLogin = () => uni.redirectTo({ url: '/pages/auth/login' })
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
  .input-field {
    margin-bottom: 24rpx;
    padding: 24rpx;
    font-size: 32rpx;
    background: #F5F5F5;
    border-radius: 12rpx;
  }
  .register-btn {
    margin-bottom: 20rpx;
    background: linear-gradient(135deg, #00DC82, #00C775);
    border: none;
    color: #FFFFFF;
  }
  .login-btn {
    background: #FFFFFF;
    color: #00DC82;
    border: 2rpx solid #00DC82;
  }
}
</style>
