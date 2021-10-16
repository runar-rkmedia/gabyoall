import { deserializeDate } from 'apiFetcher'
import { format, isToday, formatDistanceToNow } from 'date-fns'

type DateType = Date | string | undefined | null

/** Formats the date for human consumption */
const formatDate = (
  date: DateType,
  {
    format: formatStr = 'short',
    ...rest
  }: Parameters<typeof format>[2] & { format?: 'short' | 'long' } = {}
) => {
  const d = parseDate(date)
  if (!d) {
    return null
  }
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

/** parses a date, but does not validate it */
export const parseDate = (date: DateType) => {
  if (!date) {
    return null
  }
  return typeof date === 'string' ? deserializeDate(date) : date
}

/** parses and validates the that. */
export const isValidDate = (date: DateType) => {
  const d = parseDate(date)
  if (!d) {
    return null
  }
  // Our app does not deal with dates in the past...
  if (d.getFullYear() <= 1970) {
    return null
  }
  const t = d.getTime()
  if (isNaN(t)) {
    return null
  }
  return d
}

export default formatDate
