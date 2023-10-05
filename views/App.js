import { html } from 'htm'
import Meeting from './pages/meeting.js'
import { ToastProvider } from 'components'

function App() {
  return html`
    <div class="font-['Open Sans']">
      <${Meeting} />
    </div>
    <${ToastProvider} />
  `
}

export default App
