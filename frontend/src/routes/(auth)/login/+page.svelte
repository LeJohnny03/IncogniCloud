<script lang="ts">
    import { loginWithPasskey } from '$lib/client/auth';

    let username = $state('');
    let statusMessage = $state("");
    let isLoading = $state(false);
    let isError = $state(false);

    async function handleLogin(e: Event) {
        e.preventDefault();
        isLoading = true;
        isError = false;
        statusMessage = "Bitte Passkey bestätigen...";
        
        try {
            await loginWithPasskey(username);
            statusMessage = "Erfolgreich eingeloggt!";
            window.location.href = "/";
        } catch (error: any) {
            isError = true;
            statusMessage = error.message || "Ein Fehler ist aufgetreten.";
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="min-h-screen bg-background text-foreground flex flex-col items-center justify-center p-4 transition-colors duration-300">
    
    <div class="mb-8 text-center">
        <div class="w-16 h-16 bg-primary text-white dark:text-background dark:bg-accent rounded-xl mx-auto flex items-center justify-center shadow-lg mb-4 transition-colors duration-300">
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"></path>
            </svg>
        </div>
        <h1 class="text-3xl font-bold tracking-tight">IncogniCloud</h1>
    </div>

    <div class="bg-surface rounded-2xl shadow-xl border border-foregroundMuted/10 p-8 max-w-sm w-full font-sans transition-all duration-300">
        
        <h2 class="text-xl font-semibold mb-2 text-center">Willkommen zurück</h2>
        <p class="text-foregroundMuted text-sm text-center mb-6">Sicherer Zugriff nur über verifizierte Endgeräte.</p>
        
        <form onsubmit={handleLogin} class="flex flex-col gap-4">
            <div>
                <label for="username" class="block text-sm font-medium mb-1">Benutzername</label>
                <input 
                    type="text" 
                    id="username" 
                    bind:value={username} 
                    required
                    disabled={isLoading}
                    class="w-full px-4 py-2 bg-background border border-foregroundMuted/20 text-foreground rounded-xl focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all placeholder:text-foregroundMuted/50"
                    placeholder="Dein Account"
                />
            </div>

            <button 
                type="submit"
                disabled={isLoading || !username}
                class="w-full mt-2 flex items-center justify-center gap-3 bg-primary text-white dark:bg-accent dark:text-background font-medium py-3 px-4 rounded-xl disabled:opacity-40 disabled:cursor-not-allowed hover:opacity-90 transition-all duration-200 shadow-md">
                
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 11c0 3.517-1.009 6.799-2.753 9.571m-3.44-2.04l.054-.09A13.916 13.916 0 008 11a4 4 0 118 0c0 1.017-.07 2.019-.203 3m-2.118 6.844A21.88 21.88 0 0015.171 17m3.839 1.132c.645-2.266.99-4.659.99-7.132A8 8 0 008 4.07M3 15.364c.64-1.319 1-2.8 1-4.364 0-1.457.39-2.823 1.07-4"></path>
                </svg>
                {isLoading ? 'Verarbeite...' : 'Mit Passkey anmelden'}
            </button>
        </form>

        {#if statusMessage}
            <div class={`mt-6 p-4 rounded-xl text-sm text-center border ${isError ? 'bg-red-500/10 border-red-500/30 text-red-500' : 'bg-background border-foregroundMuted/10 text-foregroundMuted'}`}>
                {statusMessage}
            </div>
        {/if}

    </div>
</div>