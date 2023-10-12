import { ThemeProvider } from '@mui/material/styles'
import { signal } from '@preact/signals'
import { render } from 'preact'
import { ErrorBoundary, LocationProvider, Router, lazy } from 'preact-iso'

import { muiTheme } from 'theme'
import './global.css'

const Home = lazy(() => import('./pages'))
const HostPage = lazy(() => import('./pages/host'))
const AudiencePage = lazy(() => import('./pages/audience'))
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
    <ThemeProvider theme={muiTheme}>
      <LocationProvider>
        <main>
          <ErrorBoundary>
            <Router>
              <Home path="/" />
              <AudiencePage path="/log/:room" />
              <HostPage path="/:displayName/host" />
              <NotFound default />
            </Router>
          </ErrorBoundary>
        </main>
      </LocationProvider>
    </ThemeProvider>
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
