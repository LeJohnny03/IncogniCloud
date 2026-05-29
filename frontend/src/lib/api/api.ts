import { env } from '$env/dynamic/public';

// Im Dev-Modus (Vite) nutzen wir '/api' (wegen dem Proxy).
// Im Prod-Modus (Docker) nutzen wir die direkte Backend-URL. 
// So landen Anfragen nicht mehr bei SvelteKit (was den 404 Fehler verursacht).
export const API_BASE = import.meta.env.DEV 
    ? '/api' 
    : (env.PUBLIC_API_URL || 'http://localhost:8080/api');