import { signal } from '@preact/signals'
import { render } from 'preact'
import { ErrorBoundary, lazy, LocationProvider, Router } from 'preact-iso'

import './global.css'

const Home = lazy(() => import('./pages'))
const LogPage = lazy(() => import('./pages/log'))
const NotFound = lazy(() => import('./pages/_404.jsx'))

export const userInteractedWithDom = signal(false)

const events = ['mousedown', 'click', 'mouseup', 'keydown', 'keypress', 'scroll', 'touchstart', 'touchend', 'touchmove']
const onInteract = () => {
  userInteractedWithDom.value = true
  events.forEach((event) => {
    document.removeEventListener(event, onInteract)
  })
  if (interval) clearInterval(interval)
}

const interval = setInterval(() => {
  if (!navigator.userActivation) {
    events.forEach((event) => {
      document.addEventListener(event, onInteract)
    })
    clearInterval(interval)
    return
  }
  if (navigator.userActivation.isActive && navigator.userActivation.hasBeenActive) {
    onInteract()
  }
}, 500)

export function App() {
  return (
    <LocationProvider>
      <main>
        <ErrorBoundary>
          <Router>
            <Home path="/" />
            <LogPage path="/log/*" />
            <NotFound default />
          </Router>
        </ErrorBoundary>
      </main>
    </LocationProvider>
  )
}

render(<App />, document.getElementById('app'))

const documentHeight = () => {
  const doc = document.documentElement
  doc.style.setProperty('--doc-height', `${window.innerHeight}px`)
}
window.addEventListener('resize', documentHeight)
documentHeight()
setTimeout(documentHeight, 300)
