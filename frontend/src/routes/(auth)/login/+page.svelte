<script lang="ts">
    import { get } from '@github/webauthn-json';

    let statusMessage = $state("");
    let isLoading = $state(false);
    let isError = $state(false);

    async function loginPasskey() {
        isLoading = true;
        isError = false;
        statusMessage = "Kommunikation mit Server...";
        
        try {
            // 1. Challenge abholen
            const responseStart = await fetch('http://localhost:8080/webauthn/login/start');
            if (!responseStart.ok) throw new Error("Netzwerkfehler beim Starten des Logins.");
            const options = await responseStart.json();

            statusMessage = "Bitte Passkey bestätigen (Fingerabdruck / FaceID)...";
            
            // 2. Popup aufrufen
            const assertion = await get(options);

            statusMessage = "Verschlüsselung wird geprüft...";

            // 3. Antwort an Backend
            const responseFinish = await fetch('http://localhost:8080/webauthn/login/finish', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(assertion)
            });

            if (responseFinish.ok) {
                statusMessage = "Erfolgreich eingeloggt!";
                // HIER kommt später die Weiterleitung ins Dashboard hin:
                // window.location.href = "/files";
            } else {
                throw new Error("Login wurde vom Server abgelehnt.");
            }
        } catch (error: any) {
            isError = true;
            statusMessage = error.message || "Ein unbekannter Fehler ist aufgetreten.";
            console.error(error);
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="min-h-screen bg-slate-900 flex flex-col items-center justify-center p-4">
    
    <div class="mb-8 text-center">
        <div class="w-16 h-16 bg-blue-500 rounded-xl mx-auto flex items-center justify-center shadow-lg mb-4">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"></path></svg>
        </div>
        <h1 class="text-3xl font-bold text-white tracking-tight">My Cloud OS</h1>
    </div>

    <div class="bg-white rounded-2xl shadow-2xl p-8 max-w-sm w-full font-sans">
        
        <h2 class="text-xl font-semibold text-slate-800 mb-2 text-center">Willkommen zurück</h2>
        <p class="text-slate-500 text-sm text-center mb-8">Sicherer Zugriff nur über verifizierte Endgeräte.</p>
        
        <button 
            onclick={loginPasskey} 
            disabled={isLoading}
            class="w-full flex items-center justify-center gap-3 bg-slate-900 hover:bg-slate-800 text-white font-medium py-3 px-4 rounded-xl disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 focus:outline-none focus:ring-4 focus:ring-slate-200">
            
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 11c0 3.517-1.009 6.799-2.753 9.571m-3.44-2.04l.054-.09A13.916 13.916 0 008 11a4 4 0 118 0c0 1.017-.07 2.019-.203 3m-2.118 6.844A21.88 21.88 0 0015.171 17m3.839 1.132c.645-2.266.99-4.659.99-7.132A8 8 0 008 4.07M3 15.364c.64-1.319 1-2.8 1-4.364 0-1.457.39-2.823 1.07-4"></path></svg>
            
            {isLoading ? 'Verarbeite...' : 'Mit Passkey anmelden'}
        </button>

        {#if statusMessage}
            <div class={`mt-6 p-4 rounded-lg text-sm text-center ${isError ? 'bg-red-50 text-red-700' : 'bg-slate-50 text-slate-600'}`}>
                {statusMessage}
            </div>
        {/if}

    </div>

    <div class="mt-8 text-slate-500 text-xs text-center flex items-center gap-2">
        <svg class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
        Geschützt durch Tailscale & WebAuthn
    </div>
</div>