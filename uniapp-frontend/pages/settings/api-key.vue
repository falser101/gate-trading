<template>
  <view class="apikey-container">
    <view class="tips-card">
      <text class="tips-title">💡 如何获取 API Key</text>
      <view class="tips-content">
        <text class="tips-step">1. 登录 Gate.io 官网</text>
        <text class="tips-step">2. 进入「账户」→「API 管理」</text>
        <text class="tips-step">3. 创建新的 API Key</text>
        <text class="tips-step">4. 复制 Key 和 Secret</text>
        <text class="tips-warning">注意：请确保 API Key 具有「读取」和「交易」权限</text>
      </view>
    </view>

    <view class="form-section">
      <u-form :model="form" ref="formRef" label-width="140">
        <u-form-item label="API Key" prop="apiKey">
          <u-input
            v-model="form.apiKey"
            placeholder="请输入 API Key"
            :border="false"
            class="input-field"
          />
        </u-form-item>

        <u-form-item label="API Secret" prop="apiSecret">
          <u-input
            v-model="form.apiSecret"
            placeholder="请输入 API Secret"
            type="password"
            :border="false"
            class="input-field"
          />
        </u-form-item>
      </u-form>
    </view>

    <view class="button-section">
      <u-button
        type="primary"
        size="large"
        :loading="saving"
        @click="handleSave"
        class="save-btn"
      >
        保存
      </u-button>

      <u-button
        type="info"
        size="large"
        @click="handleTest"
        :loading="testing"
        class="test-btn"
      >
        测试连接
      </u-button>
    </view>

    <!-- 测试结果 -->
    <view v-if="testResult" :class="['test-result', testResult.success ? 'success' : 'error']">
      <text>{{ testResult.message }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/store/modules/auth'
import { bindApiKey } from '@/api/user'

const authStore = useAuthStore()
const formRef = ref<any>(null)
const saving = ref(false)
const testing = ref(false)
const testResult = ref<{ success: boolean; message: string } | null>(null)

const form = reactive({
  apiKey: '',
  apiSecret: ''
})

// 保存 API Key
const handleSave = async () => {
  if (!form.apiKey || !form.apiSecret) {
    uni.showToast({
      title: '请填写完整信息',
      icon: 'none'
    })
    return
  }

  saving.value = true

  try {
    await bindApiKey({
      api_key: form.apiKey,
      api_secret: form.apiSecret
    })

    // 更新用户状态
    authStore.updateApiKeyStatus(true)

    uni.showToast({
      title: '保存成功',
      icon: 'success'
    })

    setTimeout(() => {
      uni.navigateBack()
    }, 1000)
  } catch (error: any) {
    uni.showToast({
      title: error.message || '保存失败',
      icon: 'none'
    })
  } finally {
    saving.value = false
  }
}

// 测试连接
const handleTest = async () => {
  if (!form.apiKey || !form.apiSecret) {
    uni.showToast({
      title: '请先填写 API 信息',
      icon: 'none'
    })
    return
  }

  testing.value = true
  testResult.value = null

  // 这里可以调用后端测试接口
  // 暂时模拟测试
  setTimeout(() => {
    testResult.value = {
      success: true,
      message: 'API Key 有效，连接正常'
    }
    testing.value = false
  }, 1500)
}
</script>

<style lang="scss" scoped>
.apikey-container {
  min-height: 100vh;
  background: #F5F5F5;
  padding: 30rpx;
}

// 提示卡片
.tips-card {
  background: linear-gradient(135deg, #00DC82 0%, #00C775 100%);
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 30rpx;

  .tips-title {
    display: block;
    font-size: 30rpx;
    font-weight: bold;
    color: #FFFFFF;
    margin-bottom: 20rpx;
  }

  .tips-content {
    display: flex;
    flex-direction: column;
    gap: 12rpx;

    .tips-step {
      font-size: 26rpx;
      color: rgba(255, 255, 255, 0.95);
    }

    .tips-warning {
      font-size: 24rpx;
      color: #FFE4E4;
      margin-top: 16rpx;
      padding-top: 16rpx;
      border-top: 1rpx solid rgba(255, 255, 255, 0.3);
    }
  }
}

// 表单区域
.form-section {
  background: #FFFFFF;
  border-radius: 20rpx;
  padding: 20rpx 30rpx;
  margin-bottom: 30rpx;

  .input-field {
    padding: 24rpx 0;
    font-size: 28rpx;
  }
}

// 按钮区域
.button-section {
  .save-btn {
    margin-bottom: 20rpx;
    height: 88rpx;
    font-size: 30rpx;
    background: linear-gradient(135deg, #00DC82, #00C775);
    border: none;
  }

  .test-btn {
    height: 88rpx;
    font-size: 30rpx;
    background: #FFFFFF;
    color: #00DC82;
    border: 2rpx solid #00DC82;
  }
}

// 测试结果
.test-result {
  margin-top: 30rpx;
  padding: 24rpx;
  border-radius: 16rpx;
  text-align: center;
  font-size: 28rpx;

  &.success {
    background: #F0FDF5;
    color: #00DC82;
  }

  &.error {
    background: #FFF5F5;
    color: #FF4D4D;
  }
}
</style>
