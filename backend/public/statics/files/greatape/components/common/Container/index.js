import { html } from 'htm'

export const Container = ({ children, ...props }) => {
  return html`<div class=" mx-auto w-full px-4 lg:px-10" ...${props}>${children}</div>`
}

export default Container
