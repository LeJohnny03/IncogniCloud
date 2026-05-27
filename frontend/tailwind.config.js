/** @type {import('tailwindcss').Config} */
export default {
    content: ['./src/**/*.{html,js,svelte,ts}'],
    darkMode: 'class', // WICHTIG: Aktiviert die manuelle Steuerung über die Klasse "dark"
    theme: {
        extend: {
            // Wir definieren semantische Namen, die sich an unsere CSS-Variablen binden
            colors: {
                background: 'var(--bg-color)',
                surface: 'var(--surface-color)',
                foreground: 'var(--text-color)',
                foregroundMuted: 'var(--text-muted)',
                primary: 'var(--primary)',
                primaryLight: 'var(--primary-light)',
                accent: 'var(--accent)',
            }
        }
    },
    plugins: []
};