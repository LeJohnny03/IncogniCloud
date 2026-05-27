<script lang="ts">
    import { onMount } from 'svelte';
    import type { PageData } from './$types';

    let { data }: { data: PageData } = $props();
    let isDark = $state(false);

    onMount(() => {
        // Überprüfe beim Laden, welches Theme aktiv ist
        isDark = document.documentElement.classList.contains('dark');
    });

    function toggleTheme() {
        isDark = !isDark;
        if (isDark) {
            document.documentElement.classList.add('dark');
            localStorage.theme = 'dark';
        } else {
            document.documentElement.classList.remove('dark');
            localStorage.theme = 'light';
        }
    }

    async function logout() {
        try {
            // Sende einen Request an das Backend, damit es das HttpOnly Cookie entfernt
            await fetch('/api/logout', {
                method: 'POST',
                credentials: 'include'
            });
        } catch (e) {
            console.error("Logout Fehler", e);
        } finally {
            // Leite in jedem Fall zum Login weiter
            window.location.href = '/login';
        }
    }
</script>

<div class="min-h-screen bg-background text-foreground font-sans p-8 transition-colors duration-300">
    
    <header class="flex justify-between items-center mb-12 max-w-5xl mx-auto">
        <div class="flex items-center gap-3">
            <div class="w-10 h-10 bg-primary text-white dark:bg-accent dark:text-background rounded-lg flex items-center justify-center shadow transition-colors duration-300">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"></path>
                </svg>
            </div>
            <h1 class="text-2xl font-bold">IncogniCloud</h1>
        </div>
        
        <div class="flex items-center gap-3">
            <button 
                onclick={toggleTheme}
                class="p-2 bg-surface text-foreground border border-foregroundMuted/20 rounded-xl hover:bg-foregroundMuted/10 transition-all text-sm font-medium flex items-center gap-2 shadow-sm">
                {#if isDark}
                    ☀️ <span class="hidden sm:inline">Helles Theme</span>
                {:else}
                    🌙 <span class="hidden sm:inline">Dunkles Theme</span>
                {/if}
            </button>

            <button 
                onclick={logout}
                class="px-4 py-2 bg-surface hover:bg-red-500/10 text-foreground hover:text-red-500 border border-foregroundMuted/20 rounded-xl shadow-sm transition-all text-sm font-medium">
                Abmelden
            </button>
        </div>
    </header>

    <main class="max-w-5xl mx-auto">
        <div class="bg-surface rounded-2xl shadow-md border border-foregroundMuted/10 p-8 transition-all duration-300">
            <h2 class="text-2xl font-semibold mb-2">{data.welcomeMessage}</h2>
            <p class="text-foregroundMuted mb-8">Du bist sicher über deine Hardware authentifiziert und hast vollen Zugriff.</p>
            
            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
                <div class="p-6 bg-background hover:bg-foregroundMuted/5 rounded-xl border border-foregroundMuted/10 cursor-pointer transition-all duration-200 group">
                    <div class="w-12 h-12 bg-primary/10 text-primary dark:bg-primary-light/10 dark:text-primary-light rounded-full flex items-center justify-center mb-4 transition-colors">
                        <svg class="w-6 h-6 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
                        </svg>
                    </div>
                    <h3 class="font-semibold mb-1">Meine Dateien</h3>
                    <p class="text-sm text-foregroundMuted">Sicher verschlüsselte Cloud-Dateien verwalten.</p>
                </div>

                <div class="p-6 bg-background/50 opacity-60 rounded-xl border border-dashed border-foregroundMuted/20 flex flex-col justify-center items-center text-center">
                    <span class="text-xs font-semibold uppercase tracking-wider text-foregroundMuted bg-foregroundMuted/10 px-2 py-1 rounded-md mb-2">Coming Soon</span>
                    <p class="text-sm text-foregroundMuted">Weitere Cloud-Module folgen hier.</p>
                </div>
            </div>
        </div>
    </main>
</div>