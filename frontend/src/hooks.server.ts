import { redirect, type Handle } from '@sveltejs/kit';

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
};