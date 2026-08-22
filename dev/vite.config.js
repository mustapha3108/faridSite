import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
    plugins: [
        tailwindcss(),
    ],
    root: './glitter',
    build: {
        outDir: '../../frontend/glitter',
        emptyOutDir: false,
        cssCodeSplit: false,
        rollupOptions: {
            input: './dev.js',
            output: {
                entryFileNames: 'app.js',
                assetFileNames: 'app.css',
                format: 'iife',
                name: 'app',
            }
        }
    }
})