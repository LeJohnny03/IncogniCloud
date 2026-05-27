import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
    // Der hook garantiert, dass wir hier einen Cookie haben
    const session = cookies.get('auth_session');
    
    // Hier kannst du zukünftig Backend-Anfragen mitgeben (z.B. user profil laden)
    return {
        sessionActive: !!session,
        welcomeMessage: "Willkommen in deiner Cloud"
    };
};