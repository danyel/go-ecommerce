import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(
    {
        plugins: [react()],
        server: {
            proxy: {
                // Proxy /api to your backend
                '^/api': {
                    target: 'http://localhost:8080', // your backend server
                    changeOrigin: true,
                    rewrite: (path) => {
                        console.log(path);
                        return path; },
                    configure: (proxy, options) => {
                        proxy.on('proxyReq', (_, req, _) => {
                            console.log(`[PROXY] ${req.method} ${req.url} -> ${options.target}${req.url}`);
                        });
                    }
                }
            }
        }
    }
)