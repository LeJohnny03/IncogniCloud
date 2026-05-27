import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({ 
    plugins: [tailwindcss(), sveltekit()],
    server: {
        host: true, // Erlaubt Zugriff von außerhalb des Docker Containers
        port: 5173,
        proxy: {
            '/api': {
                target: 'http://localhost:8080', // Sobald Prod: 'http://backend:8080'
                changeOrigin: true,
            }
        }
    } 
});
