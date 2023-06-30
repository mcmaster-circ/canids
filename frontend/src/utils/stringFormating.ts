// Converts string of money amount to custom format ("100000" => "$100,000")
export const stringToCurrency = (s: string): string =>
  new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format(+s)

// Set first symbol as uppercase, ex: name => Name
export const firstToUppercase = (s: string) =>
  s && s.length > 0 ? s.charAt(0).toUpperCase() + s.slice(1) : ''

// Format string as some_name => Some name
export const formatUnderscore = (s: string) => {
  return firstToUppercase(s).replace('_', ' ')
}
