
function getParentUntil<Elem extends Element>(
  el: Element,
  mustBe: string | ((el: Element) => boolean)
): Elem | null {
  if (!el) {
    return null
  }
  const match =
    typeof mustBe === 'string'
      ? el.tagName === mustBe.toUpperCase()
      : mustBe(el)
  if (match) {
    return el as Elem
  }
  const parent = el.parentNode
  if (!parent) {
    return null
  }
  return getParentUntil(parent as any, mustBe)
}

export default getParentUntil
