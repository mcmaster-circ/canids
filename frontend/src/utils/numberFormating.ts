// Converts digits like 1000 => 1K, 1 000 000 => 1M, etc.
export const formatNumber = (v: number): string =>
  Intl.NumberFormat('en-US', {
    notation: 'compact',
    maximumFractionDigits: 1,
  }).format(v)
