import { html, useRef, useEffect } from 'preact'
import { toChildArray, cloneElement } from 'https://esm.sh/preact'

export const Tooltip = ({ children, label }) => {
  const ref = useRef()
  const component = toChildArray(children)[0]
  const tippyInstance = useRef()
  useEffect(() => {
    if (ref.current.base)
      tippyInstance.current = tippy(ref.current.base, {
        content: label,
        arrow: false,
        hideOnClick: false,
      })
  }, [])

  useEffect(() => {
    tippyInstance.current.setContent(label)
  }, [label])

  if (component) return cloneElement(component, { ref })
  return null
}
