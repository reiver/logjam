import { html } from 'htm'

export const Icon = ({ icon, ...props }) => {
  return html`<svg data-src="assets/icons/${icon}.svg" ...${props} />`
}

export default Icon
