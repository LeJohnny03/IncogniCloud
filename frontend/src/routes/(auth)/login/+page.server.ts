import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
    // Prüft, ob der Nutzer bereits eingeloggt ist
    const session = cookies.get('auth_session');
    
    // Wenn ja, direkt auf die Startseite weiterleiten
    if (session) {
        throw redirect(303, '/');
    }
    
    // Ansonsten Seite normal laden lassen
    return {};
};