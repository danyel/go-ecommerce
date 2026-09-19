import { createRoot } from 'react-dom/client'
import './App.css'
import App from './App.tsx'
import { BrowserRouter } from 'react-router-dom'
import { GlobalStateProvider } from './state/global-state.tsx';
import {AuthProvider} from './state/auth-provider.tsx';

createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <GlobalStateProvider>
      <AuthProvider>
        <App />
      </AuthProvider>
    </GlobalStateProvider>
  </BrowserRouter>,
)
