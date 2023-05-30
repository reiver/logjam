import { html } from 'htm'
import clsx from 'clsx'
import { forwardRef } from 'preact/compat'

export const IconButton = forwardRef(({ children, variant, class: className, ...props }, ref) => {
  return html`<button
    ref=${ref}
    class="${clsx('transition rounded-full p-2', className, {
      'bg-red-distructive hover:bg-red-700 ': variant === 'danger',
      'bg-gray-900 bg-opacity-40 text-white hover:bg-opacity-50': variant === 'ghost',
      'dark:text-gray-800 text-gray-100 dark:bg-gray-100 hover:dark:bg-gray-300 bg-gray-900 hover:bg-gray-700': !variant,
    })}"
    ...${props}
  >
    ${children}
  </button>`
})

export default IconButton
