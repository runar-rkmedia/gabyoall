import { deserializeDate } from 'apiFetcher'
import { format, isToday, formatDistanceToNow } from 'date-fns'

const formatDate = (
  date: Date | string | undefined | null,
  {
    format: formatStr = 'short',
    ...rest
  }: Parameters<typeof format>[2] & { format?: 'short' | 'long' } = {}
) => {
  if (!date) {
    return null
  }
  const d = typeof date === 'string' ? deserializeDate(date) : date
  if (formatStr === 'short' && isToday(d)) {
    return formatDistanceToNow(d, {
      includeSeconds: false,
      addSuffix: true,
      locale: rest.locale,
    })
  }
  const f =
    {
      short: 'Pp',
      long: 'PPpp',
    }[formatStr] || formatStr
  return format(d, f, rest)
}

export default formatDate
