export const storage = {
  set: (key: string, value: any): void => {
    try {
      const str = typeof value === 'string' ? value : JSON.stringify(value)
      uni.setStorageSync(key, str)
    } catch (error) {
      console.error('Storage set error:', error)
    }
  },

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

  remove: (key: string): void => {
    try {
      uni.removeStorageSync(key)
    } catch (error) {
      console.error('Storage remove error:', error)
    }
  },

  clear: (): void => {
    try {
      uni.clearStorageSync()
    } catch (error) {
      console.error('Storage clear error:', error)
    }
  }
}
