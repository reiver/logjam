import { html, render } from 'htm'
import App from './App.js'
import 'external-svg-loader'
import { signal } from '@preact/signals'

export const userInteractedWithDom = signal(false)

const interval = setInterval(() => {
  if (navigator.userActivation.isActive && navigator.userActivation.hasBeenActive) {
    userInteractedWithDom.value = true
    clearInterval(interval)
  }
}, 500)

render(html`<${App} />`, document.getElementById('root'))
