import formatterYaml from 'yaml'
import formatterToml from 'toml-js'

export function convertStringToCodeFormat(code: string | {}, format: string): readonly [string, null] | readonly [null, string] {
  if (!code) {
    return ['', null] as const
  }
  if (!format) {
    return [null, 'missing format'] as const
  }
  try {
    let obj
    if (typeof code === 'string') {
      const [c, _, err] = parseStringCode(code)
      if (err) {
        return [null, err]
      }
      obj = c
    } else {
      obj = code
    }
    switch (format) {
      case 'yaml':
      case 'yml':
        return [formatterYaml.stringify(obj, { sortMapEntries: true }), null] as const
      case 'toml':
        return [formatterToml.dump(obj), null] as const
      case 'json':
        return [JSON.stringify(obj, null, 2), null] as const
      default:
        console.error('Unsupported format', format)
        break
    }
  } catch (error) {
    console.error('failed to convert to format', { code, format, error })
    return [null, error.message] as const
  }
  return [null, 'unhandled']
}

const parseStringCode = (s: string): readonly [code: {}, kind: 'yaml' | 'toml' | 'json', error: null] | readonly [code: null, kind: null, error: string] => {
  try {

    if (s.startsWith('{') || s.startsWith('[')) {
      // smells like JSON
      return [JSON.parse(s), 'json', null]
    }
    if (s.includes(' = ')) {
      return [formatterToml.parse(s), 'toml', null]
    }
    return [formatterYaml.parse(s), 'yaml', null]
  } catch (err) {
    console.error('failed to parse code', err, s)
    return [null, null, err.message]
  }
}
