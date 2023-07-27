import { html } from 'htm'
export const Logo = () => {
  return html`<div class="flex gap-1 items-center dark:text-white text-black">
    <svg data-src="assets/images/logo.svg" />
  </div>`
}

export default Logo
