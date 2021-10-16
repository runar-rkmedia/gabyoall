const encode_regex = /[\+=\/]/g
const decode_regex = /[\._\-]/g

export function decode(str: string) {
  const replaced = str.replace(decode_regex, decodeChar) // + '='.repeat(str.length % 4)

  try {
    return [atob(replaced), null] as const
  } catch (err) {
    console.error('failed during decode: ', {
      lenght: str.length,
      string: str,
      replaced,
      err,
    })
    return ['', err] as const
  }
}

function encodeChar(c) {
  switch (c) {
    case '+':
      return '.'
    case '=':
      return '-'
    case '/':
      return '_'
  }
}

function decodeChar(c) {
  switch (c) {
    case '.':
      return '+'
    case '-':
      return '='
    case '_':
      return '/'
  }
  return ''
}
function decodeJWT(s: string, pretty?: boolean) {
  if (!s) {
    return []
  }
  const parts: [string, string | null][] = s
    .split('.')
    .slice(0, 2)
    .map((p) => {
      const [b, err] = decode(p)
      if (err) {
        return [b, err.message as string]
      }
      if (pretty) {
        try {
          return [JSON.stringify(JSON.parse(b), null, 2), null]
        } catch (error) {
          return [b, err.message as string]
        }
      }
      return [b, null]
    })
  return parts
}

export default decodeJWT
