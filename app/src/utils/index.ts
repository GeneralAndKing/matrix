export const formatNumber = (num: number): string => {
  if (num >= 1e8) {
    return (num / 1e8).toFixed(2) + 'b'
  } else if (num >= 1e4) {
    return (num / 1e4).toFixed(2) + 'w'
  } else if (num >= 1e3) {
    return (num / 1e3).toFixed(2) + 'k'
  } else {
    return num.toString()
  }
}

export const formatChineseNumber = (num: number): string => {
  if (num >= 1e8) {
    return (num / 1e8).toFixed(2) + '亿'
  } else if (num >= 1e4) {
    return (num / 1e4).toFixed(2) + '万'
  } else if (num >= 1e3) {
    return (num / 1e3).toFixed(2) + '千'
  } else {
    return num.toString()
  }
}
