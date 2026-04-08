// API 基础配置
const BASE_URL = '/api'
const TIMEOUT = 30000

// 请求拦截器
const request = (options: UniApp.RequestOptions) => {
  return new Promise<UniApp.RequestSuccessCallbackResult>((resolve, reject) => {
    // 获取 Token
    const token = uni.getStorageSync('token') || ''

    // 设置请求头
    const header = {
      'Content-Type': 'application/json',
      ...options.header
    }

    if (token) {
      header['Authorization'] = `Bearer ${token}`
    }

    uni.request({
      ...options,
      url: BASE_URL + options.url,
      header,
      timeout: TIMEOUT,
      success: (res) => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res)
        } else if (res.statusCode === 401) {
          uni.removeStorageSync('token')
          uni.removeStorageSync('userInfo')
          uni.reLaunch({ url: '/pages/auth/login' })
          reject(new Error('请先登录'))
        } else {
          const errorMsg = res.data?.message || '请求失败'
          uni.showToast({ title: errorMsg, icon: 'none' })
          reject(new Error(errorMsg))
        }
      },
      fail: (err) => {
        uni.showToast({ title: '网络错误', icon: 'none' })
        reject(err)
      }
    })
  })
}

export const get = <T>(url: string, data?: Record<string, any>) => {
  return request<T>({ url, method: 'GET', data })
}

export const post = <T>(url: string, data?: Record<string, any>) => {
  return request<T>({ url, method: 'POST', data })
}

export const put = <T>(url: string, data?: Record<string, any>) => {
  return request<T>({ url, method: 'PUT', data })
}

export const del = <T>(url: string) => {
  return request<T>({ url, method: 'DELETE' })
}

export default request
