/*import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
    // Hole den Session-Cookie, der vom Go-Backend nach dem Login gesetzt wird
    const session = event.cookies.get('auth_session');
    
    const isLoginRoute = event.url.pathname.startsWith('/login');
    const isApiRoute = event.url.pathname.startsWith('/api'); // Falls du SvelteKit API Routen nutzt

    // Wenn keine Session existiert und man NICHT auf der Login-Seite ist -> redirect zu /login
    if (!session && !isLoginRoute && !isApiRoute) {
        throw redirect(303, '/login');
    }

    const response = await resolve(event);
    return response;
};*/

import { API_BASE } from '$lib/api/api';
import { BACKEND_INTERNAL_URL } from '$env/static/private';
import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
    const session = event.cookies.get('auth_session');
    
    const isLoginRoute = event.url.pathname.startsWith('/login');
    const isSetupRoute = event.url.pathname.startsWith('/setup');
    const isApiRoute = event.url.pathname.startsWith('/api');

    // 1. Setup-Status vom Backend abfragen (könnte man später cachen für mehr Performance)
    try {
        // Da wir im Server-Hook sind, nutzen wir fetch direkt gegen den internen Docker-Namen/Port
        // Wenn du lokal ohne Docker entwickelst, ersetze das durch http://localhost:8080
        const setupResp = await fetch(`${BACKEND_INTERNAL_URL}${API_BASE}/setup/status`);
        if (setupResp.ok) {
            const data = await setupResp.json();
            
            // System ist frisch! User MUSS auf die /setup Seite.
            if (data.needs_setup && !isSetupRoute && !isApiRoute) {
                throw redirect(303, '/setup');
            }
            
            // System ist eingerichtet, aber User versucht /setup aufzurufen -> blockieren
            if (!data.needs_setup && isSetupRoute) {
                throw redirect(303, '/login');
            }
        }
    } catch (e) {
        // Falls das Backend mal nicht erreichbar ist, ignorieren wir das hier kurz
        console.error("Konnte Setup-Status nicht prüfen:", e);
    }

    // 2. Normale Session-Prüfung (wie bisher)
    if (!session && !isLoginRoute && !isSetupRoute && !isApiRoute) {
        throw redirect(303, '/login');
    }

    return await resolve(event);
};