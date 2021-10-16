// https://stackoverflow.com/a/58436959
export type Join<K, P> = K extends string | number
  ? P extends string | number
    ? `${K}${'' extends P ? '' : '.'}${P}`
    : never
  : never

// This can cause the types to become deep fast if there are too many, so beware of that.
type Prev = [never, 0, 1, 2, 3, 4, 5, ...0[]]

/** Returns the Paths for an object https://stackoverflow.com/a/58436959 */
export type Paths<T, D extends number = 10> = [D] extends [never]
  ? never
  : T extends object
  ? {
      [K in keyof T]-?: K extends string | number
        ? `${K}` | Join<K, Paths<T[K], Prev[D]>>
        : never
    }[keyof T]
  : ''
/** https://stackoverflow.com/a/58436959 */
export type Leaves<T, D extends number = 10> = [D] extends [never]
  ? never
  : T extends object
  ? { [K in keyof T]-?: Join<K, Leaves<T[K], Prev[D]>> }[keyof T]
  : ''
