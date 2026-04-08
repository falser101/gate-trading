export const formatNumber = (num: string | number, decimals = 2): string => {
  const value = typeof num === 'string' ? parseFloat(num) : num
  return isNaN(value) ? '0' : value.toFixed(decimals)
}

export const formatPercent = (num: string | number, decimals = 2): string => {
  const value = typeof num === 'string' ? parseFloat(num) : num
  return isNaN(value) ? '0%' : `${value.toFixed(decimals)}%`
}

export const formatPnl = (num: string | number): string => {
  const value = typeof num === 'string' ? parseFloat(num) : num
  if (isNaN(value)) return '0'
  return value > 0 ? `+${value.toFixed(2)}` : value.toFixed(2)
}

export const formatLargeNumber = (num: number): string => {
  if (num >= 100000000) return `${(num / 100000000).toFixed(2)}亿`
  if (num >= 10000) return `${(num / 10000).toFixed(2)}万`
  return num.toString()
}

export const formatRelativeTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  const diff = Date.now() - date.getTime()
  const minute = 60 * 1000
  const hour = 60 * minute
  const day = 24 * hour
  const month = 30 * day

  if (diff < minute) return '刚刚'
  if (diff < hour) return `${Math.floor(diff / minute)}分钟前`
  if (diff < day) return `${Math.floor(diff / hour)}小时前`
  if (diff < month) return `${Math.floor(diff / day)}天前`
  return `${Math.floor(diff / month)}个月前`
}
