// 格式化数字，保留指定位数的小数
export const formatNumber = (num: string | number, decimals: number = 2): string => {
  const value = typeof num === 'string' ? parseFloat(num) : num
  if (isNaN(value)) return '0'
  return value.toFixed(decimals)
}

// 格式化百分比
export const formatPercent = (num: string | number, decimals: number = 2): string => {
  const value = typeof num === 'string' ? parseFloat(num) : num
  if (isNaN(value)) return '0%'
  return `${value.toFixed(decimals)}%`
}

// 格式化盈亏（带颜色标记）
export const formatPnl = (num: string | number): { value: string; isPositive: boolean } => {
  const value = typeof num === 'string' ? parseFloat(num) : num
  if (isNaN(value)) return { value: '0', isPositive: true }
  return {
    value: value > 0 ? `+${value.toFixed(2)}` : value.toFixed(2),
    isPositive: value >= 0
  }
}

// 格式化大数字（万、亿）
export const formatLargeNumber = (num: number): string => {
  if (num >= 100000000) {
    return `${(num / 100000000).toFixed(2)}亿`
  }
  if (num >= 10000) {
    return `${(num / 10000).toFixed(2)}万`
  }
  return num.toString()
}

// 格式化日期
export const formatDate = (dateStr: string, format: string = 'YYYY-MM-DD'): string => {
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return format
    .replace('YYYY', year.toString())
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

// 相对时间格式化
export const formatRelativeTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  const minute = 60 * 1000
  const hour = 60 * minute
  const day = 24 * hour
  const month = 30 * day
  const year = 365 * day

  if (diff < minute) {
    return '刚刚'
  } else if (diff < hour) {
    return `${Math.floor(diff / minute)}分钟前`
  } else if (diff < day) {
    return `${Math.floor(diff / hour)}小时前`
  } else if (diff < month) {
    return `${Math.floor(diff / day)}天前`
  } else if (diff < year) {
    return `${Math.floor(diff / month)}个月前`
  } else {
    return `${Math.floor(diff / year)}年前`
  }
}
