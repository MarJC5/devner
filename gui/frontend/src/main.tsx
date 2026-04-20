import React from 'react'
import {createRoot} from 'react-dom/client'
// Wails template had a default style.css with `text-align: center` on
// html — override via App.css only, which Tailwind bootstraps.
import './App.css'
import App from './App'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <App/>
    </React.StrictMode>
)
