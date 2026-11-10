// 时间与文本格式化工具

// 将 RFC3339 时间字符串格式化为 YYYY-MM-DD HH:mm
export function formatDateTime(value) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return (
    `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ` +
    `${pad(d.getHours())}:${pad(d.getMinutes())}`
  )
}

// 将 RFC3339 时间字符串格式化为 YYYY-MM-DD
export function formatDate(value) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

// 相对时间：刚刚 / x 分钟前 / x 小时前 / x 天前 / 具体日期
export function fromNow(value) {
  if (!value) return ''
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return ''
  const diff = Date.now() - d.getTime()
  const sec = Math.floor(diff / 1000)
  if (sec < 60) return '刚刚'
  const min = Math.floor(sec / 60)
  if (min < 60) return `${min} 分钟前`
  const hour = Math.floor(min / 60)
  if (hour < 24) return `${hour} 小时前`
  const day = Math.floor(hour / 24)
  if (day < 30) return `${day} 天前`
  return formatDate(value)
}

// 摘要截断：超出 length 时截断并追加省略号
export function truncate(text, length = 120) {
  if (!text) return ''
  const str = String(text)
  if (str.length <= length) return str
  return `${str.slice(0, length)}…`
}

// 取名称首字母（用于头像占位）
export function initialOf(name) {
  if (!name) return '?'
  return String(name).trim().charAt(0).toUpperCase()
}
