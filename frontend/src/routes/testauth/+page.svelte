<script lang="ts">
    // 'create' ist für Registrierung, 'get' ist für den Login
    import { create, get } from '@github/webauthn-json';

    const state = $state({
        statusMessage: "Bereit für den Test.",
        isLoading: false
    });

    // --- REGISTRIERUNG ---
    async function registerPasskey() {
        state.isLoading = true;
        state.statusMessage = "Hole Challenge vom Server...";
        
        try {
            // 1. Challenge vom Backend abholen
            const responseStart = await fetch('http://localhost:8080/webauthn/register/start');
            if (!responseStart.ok) throw new Error("Backend Fehler bei Start");
            const options = await responseStart.json();

            state.statusMessage = "Warte auf Fingerabdruck/FaceID...";
            
            // 2. Browser-Popup triggern (Wandelt Base64 automatisch in ArrayBuffer um!)
            const credential = await create(options);

            state.statusMessage = "Sende Passkey an Server...";

            // 3. Ergebnis an Backend senden
            const responseFinish = await fetch('http://localhost:8080/webauthn/register/finish', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(credential)
            });

            if (responseFinish.ok) {
                state.statusMessage = "✅ Registrierung erfolgreich!";
            } else {
                throw new Error("Server hat den Passkey abgelehnt.");
            }
        } catch (error: any) {
            state.statusMessage = `❌ Fehler: ${error.message}`;
            console.error(error);
        } finally {
            state.isLoading = false;
        }
    }

    // --- LOGIN ---
    async function loginPasskey() {
        state.isLoading = true;
        state.statusMessage = "Hole Login-Challenge vom Server...";
        
        try {
            // 1. Challenge vom Backend abholen
            const responseStart = await fetch('http://localhost:8080/webauthn/login/start');
            if (!responseStart.ok) throw new Error("Backend Fehler bei Login-Start");
            const options = await responseStart.json();

            state.statusMessage = "Warte auf Fingerabdruck/FaceID...";
            
            // 2. Browser-Popup triggern (WICHTIG: Hier nutzen wir get() statt create()!)
            const assertion = await get(options);

            state.statusMessage = "Prüfe Passkey auf dem Server...";

            // 3. Ergebnis an Backend senden
            const responseFinish = await fetch('http://localhost:8080/webauthn/login/finish', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(assertion)
            });

            if (responseFinish.ok) {
                state.statusMessage = "✅ Login erfolgreich! Willkommen zurück.";
            } else {
                throw new Error("Server hat den Login abgelehnt.");
            }
        } catch (error: any) {
            state.statusMessage = `❌ Fehler: ${error.message}`;
            console.error(error);
        } finally {
            state.isLoading = false;
        }
    }
</script>

<div class="min-h-screen bg-gray-100 flex flex-col items-center justify-center p-4">
    <div class="bg-white rounded-xl shadow-lg p-8 max-w-md w-full text-center font-sans">
        
        <h1 class="text-2xl font-bold mb-6 text-gray-800">Passkey Test-Labor</h1>
        
        <div class="space-y-4">
            <button 
                onclick={registerPasskey} 
                disabled={state.isLoading}
                class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded disabled:opacity-50 transition-colors">
                1. Passkey erstellen (Registrieren)
            </button>

            <button 
                onclick={loginPasskey} 
                disabled={state.isLoading}
                class="w-full bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-4 rounded disabled:opacity-50 transition-colors">
                2. Mit Passkey einloggen
            </button>
        </div>

        <div class="mt-8 p-4 bg-gray-50 rounded border border-gray-200">
            <p class="text-sm text-gray-600 font-mono wrap-break-word">{state.statusMessage}</p>
        </div>

    </div>
</div>