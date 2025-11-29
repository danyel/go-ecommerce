import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './App.css'
import App from './App.tsx'
import {BrowserRouter} from "react-router-dom"
import {GlobalStateProvider} from "./state/global-state.tsx";

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <BrowserRouter>
            <GlobalStateProvider>
                <App/>
            </GlobalStateProvider>
        </BrowserRouter>
    </StrictMode>,
)
