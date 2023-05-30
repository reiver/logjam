import { html } from 'htm'
import clsx from 'clsx'

export const Button = ({ children, variant, ...props }) => {
  let variantClasses = ''
  switch (variant) {
    case 'red':
      variantClasses = 'bg-red-700 hover:bg-red-800 text-white'
      break
    default:
      variantClasses = 'bg-blue-500 hover:bg-blue-600 text-white'
  }
  return html`<button class="${clsx('transition-all px-6 py-2 rounded-[4px] font-bold text-bold-12', variantClasses, props.class)}" ...${props}>${children}</button>`
}
export default Button
