// 本地存储封装
export const storage = {
  // 设置存储
  set: (key: string, value: any): void => {
    try {
      const str = typeof value === 'string' ? value : JSON.stringify(value)
      uni.setStorageSync(key, str)
    } catch (error) {
      console.error('Storage set error:', error)
    }
  },

  // 获取存储
  get: <T>(key: string): T | null => {
    try {
      const str = uni.getStorageSync(key)
      if (!str) return null
      try {
        return JSON.parse(str) as T
      } catch {
        return str as T
      }
    } catch (error) {
      console.error('Storage get error:', error)
      return null
    }
  },

  // 移除存储
  remove: (key: string): void => {
    try {
      uni.removeStorageSync(key)
    } catch (error) {
      console.error('Storage remove error:', error)
    }
  },

  // 清空存储
  clear: (): void => {
    try {
      uni.clearStorageSync()
    } catch (error) {
      console.error('Storage clear error:', error)
    }
  }
}
