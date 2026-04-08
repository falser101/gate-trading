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
        // HTTP 状态码处理
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res)
        } else if (res.statusCode === 401) {
          // Token 过期，清除并跳转登录
          uni.removeStorageSync('token')
          uni.removeStorageSync('userInfo')
          uni.reLaunch({
            url: '/pages/auth/login'
          })
          reject(new Error('请先登录'))
        } else {
          // 其他错误
          const errorMsg = res.data?.message || '请求失败'
          uni.showToast({
            title: errorMsg,
            icon: 'none'
          })
          reject(new Error(errorMsg))
        }
      },
      fail: (err) => {
        uni.showToast({
          title: '网络错误',
          icon: 'none'
        })
        reject(err)
      }
    })
  })
}

// GET 请求
export const get = <T>(url: string, data?: Record<string, any>) => {
  return request<T>({
    url,
    method: 'GET',
    data
  })
}

// POST 请求
export const post = <T>(url: string, data?: Record<string, any>) => {
  return request<T>({
    url,
    method: 'POST',
    data
  })
}

// PUT 请求
export const put = <T>(url: string, data?: Record<string, any>) => {
  return request<T>({
    url,
    method: 'PUT',
    data
  })
}

// DELETE 请求
export const del = <T>(url: string) => {
  return request<T>({
    url,
    method: 'DELETE'
  })
}

export default request
