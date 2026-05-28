<script lang="ts">
    import { registerWithPasskey } from '$lib/client/auth'; // Wir nutzen deine Passkey-Logik (leicht abgewandelt für Registrierung später)
    import { onMount } from 'svelte';

    // Wizard State
    let currentStep = $state(1);
    let isSubmitting = $state(false);
    let error = $state("");

    // Step 1: Admin Account
    let adminUsername = $state("admin");

    // Step 2: Storage
    let selectedStorageDrive = $state("");
    let storageFolder = $state("/incognicloud_data");
    
    // Step 3: Backups
    let enableBackups = $state(false);
    let selectedBackupDrive = $state("");
    let backupType = $state("incremental"); // incremental, full, immutable

    // Step 4: Cloud Settings
    let cloudFolderName = $state("IncogniCloud");

    // Mock-Daten vom Backend (später kommt das über fetch('/api/setup/drives'))
    let availableDrives = $state([
        { path: "/mnt/hdd1", name: "Seagate IronWolf 4TB", freeSpace: "3.6 TB" },
        { path: "/mnt/nvme", name: "Samsung 980 Pro 1TB", freeSpace: "900 GB" },
        { path: "/mnt/usb_backup", name: "WD Elements USB", freeSpace: "5 TB" }
    ]);

    function nextStep() {
        error = "";
        if (currentStep === 1 && !adminUsername) {
            error = "Bitte einen Benutzernamen festlegen.";
            return;
        }
        if (currentStep === 2 && !selectedStorageDrive) {
            error = "Bitte wähle eine Festplatte für den Hauptspeicher aus.";
            return;
        }
        if (currentStep === 3 && enableBackups && !selectedBackupDrive) {
            error = "Bitte wähle ein Backup-Laufwerk aus.";
            return;
        }
        if (currentStep < 4) currentStep++;
    }

    function prevStep() {
        if (currentStep > 1) currentStep--;
        error = "";
    }

    async function finishSetup() {
        isSubmitting = true;
        error = "";
        
        const setupConfig = {
            adminUsername,
            storageDrive: selectedStorageDrive,
            storageFolder,
            enableBackups,
            backupDrive: selectedBackupDrive,
            backupType,
            cloudFolderName
        };

        try {
            // 1. Erstelle den Admin-Nutzer mit Passkey über deine API
            await registerWithPasskey(adminUsername, adminUsername);
            
            // 2. Sende die Konfiguration für die Festplatten/Cloud an das Backend
            const configResp = await fetch('/api/setup/finish', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(setupConfig)
            });

            if (!configResp.ok) throw new Error("Fehler beim Speichern der Konfiguration.");
            
            // 3. Fertig! Ab zum Login.
            window.location.href = '/login';
        } catch (e: any) {
            error = e.message || "Ein Fehler ist beim Setup aufgetreten.";
            console.error(e);
        } finally {
            isSubmitting = false;
        }
    }
</script>

<div class="min-h-screen bg-background text-foreground flex flex-col items-center justify-center p-4 sm:p-8 transition-colors duration-300">
    
    <div class="max-w-2xl w-full mb-8 text-center">
        <div class="w-16 h-16 bg-primary text-white dark:text-background dark:bg-accent rounded-xl mx-auto flex items-center justify-center shadow-lg mb-4">
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"></path></svg>
        </div>
        <h1 class="text-3xl font-bold tracking-tight mb-2">IncogniCloud Setup</h1>
        <p class="text-foregroundMuted">Richte dein System für den ersten Start ein.</p>
    </div>

    <div class="bg-surface rounded-2xl shadow-xl border border-foregroundMuted/10 p-6 sm:p-10 max-w-2xl w-full transition-all duration-300">
        
        <div class="flex gap-2 mb-8">
            {#each [1, 2, 3, 4] as step}
                <div class={`h-2 flex-1 rounded-full transition-all duration-500 ${currentStep >= step ? 'bg-primary' : 'bg-foregroundMuted/20'}`}></div>
            {/each}
        </div>

        {#if error}
            <div class="mb-6 p-4 bg-red-500/10 border border-red-500/30 text-red-500 rounded-xl text-sm">
                {error}
            </div>
        {/if}

        {#if currentStep === 1}
            <div class="animate-in fade-in slide-in-from-right-4 duration-500">
                <h2 class="text-2xl font-semibold mb-4">1. Master-Admin erstellen</h2>
                <p class="text-foregroundMuted text-sm mb-6">Dieser Account erhält alle Systemrechte und kann andere Nutzer einladen. Der Login erfolgt im nächsten Schritt via Passkey.</p>
                
                <div>
                    <label for="adminUsername" class="block text-sm font-medium mb-1">Admin Benutzername</label>
                    <input id="adminUsername" type="text" bind:value={adminUsername} class="w-full px-4 py-3 bg-background border border-foregroundMuted/20 rounded-xl focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all" placeholder="z.B. admin oder dein Name" />
                </div>
            </div>
        {/if}

        {#if currentStep === 2}
            <div class="animate-in fade-in slide-in-from-right-4 duration-500">
                <h2 class="text-2xl font-semibold mb-4">2. Cloud-Speicher wählen</h2>
                <p class="text-foregroundMuted text-sm mb-6">Wähle die physische Festplatte aus, auf der die Cloud-Dateien gespeichert werden sollen.</p>
                
                <div class="space-y-3 mb-6">
                    {#each availableDrives as drive}
                        <label class={`flex items-center p-4 border rounded-xl cursor-pointer transition-all ${selectedStorageDrive === drive.path ? 'border-primary bg-primary/5 shadow-sm' : 'border-foregroundMuted/20 hover:bg-foregroundMuted/5'}`}>
                            <input type="radio" name="storageDrive" value={drive.path} bind:group={selectedStorageDrive} class="hidden" />
                            <div class="w-5 h-5 rounded-full border-2 border-foregroundMuted/30 flex items-center justify-center mr-4">
                                {#if selectedStorageDrive === drive.path}<div class="w-2.5 h-2.5 bg-primary rounded-full"></div>{/if}
                            </div>
                            <div class="flex-1">
                                <div class="font-medium">{drive.name}</div>
                                <div class="text-xs text-foregroundMuted">{drive.path}</div>
                            </div>
                            <div class="text-sm font-semibold text-foregroundMuted">{drive.freeSpace} frei</div>
                        </label>
                    {/each}
                </div>

                {#if selectedStorageDrive}
                    <div class="animate-in fade-in slide-in-from-top-2">
                        <label for="storageFolder" class="block text-sm font-medium mb-1">Speicherordner auf dem Laufwerk</label>
                        <div class="flex items-center gap-2">
                            <span class="text-foregroundMuted">{selectedStorageDrive}</span>
                            <input id="storageFolder" type="text" bind:value={storageFolder} class="flex-1 px-4 py-2 bg-background border border-foregroundMuted/20 rounded-xl focus:ring-2 focus:ring-primary focus:border-primary outline-none" />
                        </div>
                    </div>
                {/if}
            </div>
        {/if}

        {#if currentStep === 3}
            <div class="animate-in fade-in slide-in-from-right-4 duration-500">
                <h2 class="text-2xl font-semibold mb-4">3. Backup-Strategie</h2>
                <p class="text-foregroundMuted text-sm mb-6">Möchtest du automatische Backups deiner Cloud-Daten auf eine zweite Festplatte aktivieren?</p>
                
                <label class="flex items-center gap-3 p-4 border border-foregroundMuted/20 rounded-xl cursor-pointer mb-6 hover:bg-foregroundMuted/5 transition-colors">
                    <input type="checkbox" bind:checked={enableBackups} class="w-5 h-5 rounded border-foregroundMuted/30 text-primary focus:ring-primary" />
                    <span class="font-medium">Lokale Backups aktivieren</span>
                </label>

                {#if enableBackups}
                    <div class="animate-in fade-in slide-in-from-top-4">
                        <h3 class="text-sm font-medium mb-3">Backup-Ziellaufwerk</h3>
                        <div class="space-y-3 mb-6">
                            {#each availableDrives.filter(d => d.path !== selectedStorageDrive) as drive}
                                <label class={`flex items-center p-3 border rounded-xl cursor-pointer transition-all ${selectedBackupDrive === drive.path ? 'border-primary bg-primary/5 shadow-sm' : 'border-foregroundMuted/20 hover:bg-foregroundMuted/5'}`}>
                                    <input type="radio" name="backupDrive" value={drive.path} bind:group={selectedBackupDrive} class="hidden" />
                                    <div class="w-4 h-4 rounded-full border-2 border-foregroundMuted/30 flex items-center justify-center mr-3">
                                        {#if selectedBackupDrive === drive.path}<div class="w-2 h-2 bg-primary rounded-full"></div>{/if}
                                    </div>
                                    <div class="flex-1 font-medium text-sm">{drive.name}</div>
                                </label>
                            {/each}
                        </div>

                        <h3 class="text-sm font-medium mb-3">Art des Backups</h3>
                        <select bind:value={backupType} class="w-full px-4 py-3 bg-background border border-foregroundMuted/20 rounded-xl focus:ring-2 focus:ring-primary outline-none">
                            <option value="incremental">Inkrementell (Spart Platz, schnell)</option>
                            <option value="immutable">Immutable (Schutz vor Ransomware)</option>
                            <option value="mirror">Spiegelung (1:1 Kopie)</option>
                        </select>
                    </div>
                {/if}
            </div>
        {/if}

        {#if currentStep === 4}
            <div class="animate-in fade-in slide-in-from-right-4 duration-500">
                <h2 class="text-2xl font-semibold mb-4">4. Feinschliff</h2>
                <p class="text-foregroundMuted text-sm mb-6">Lege den Namen für den zentralen Sync-Ordner fest. Dieser Name wird später auf deinen Laptops/PCs angezeigt (ähnlich wie "Dropbox" oder "OneDrive").</p>
                
                <div class="mb-8">
                    <label for="cloudFolderName" class="block text-sm font-medium mb-1">Zentraler Cloud-Ordner Name</label>
                    <input id="cloudFolderName" type="text" bind:value={cloudFolderName} class="w-full px-4 py-3 bg-background border border-foregroundMuted/20 rounded-xl focus:ring-2 focus:ring-primary focus:border-primary outline-none transition-all" placeholder="z.B. MeineCloud" />
                </div>

                <div class="bg-foregroundMuted/5 p-4 rounded-xl text-sm border border-foregroundMuted/10">
                    <h4 class="font-semibold mb-2">Zusammenfassung:</h4>
                    <ul class="space-y-1 text-foregroundMuted">
                        <li>• Admin: <span class="text-foreground font-medium">{adminUsername}</span></li>
                        <li>• Speicher: <span class="text-foreground font-medium">{selectedStorageDrive}{storageFolder}</span></li>
                        <li>• Backups: <span class="text-foreground font-medium">{enableBackups ? 'Aktiviert' : 'Deaktiviert'}</span></li>
                        <li>• Ordnername: <span class="text-foreground font-medium">{cloudFolderName}</span></li>
                    </ul>
                </div>
            </div>
        {/if}

        <div class="mt-8 flex justify-between items-center pt-6 border-t border-foregroundMuted/10">
            <button 
                onclick={prevStep}
                class={`px-6 py-2.5 rounded-xl font-medium transition-all ${currentStep === 1 ? 'invisible' : 'bg-background hover:bg-foregroundMuted/10 text-foreground border border-foregroundMuted/20'}`}>
                Zurück
            </button>

            {#if currentStep < 4}
                <button 
                    onclick={nextStep}
                    class="px-6 py-2.5 bg-primary hover:opacity-90 text-white dark:bg-accent dark:text-background rounded-xl font-medium transition-all shadow-md">
                    Weiter
                </button>
            {:else}
                <button 
                    onclick={finishSetup}
                    disabled={isSubmitting}
                    class="px-6 py-2.5 bg-primary hover:opacity-90 text-white dark:bg-accent dark:text-background rounded-xl font-medium transition-all shadow-md disabled:opacity-50 flex items-center gap-2">
                    {#if isSubmitting}
                        <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
                        Einrichten...
                    {:else}
                        Setup abschließen & Passkey scannen
                    {/if}
                </button>
            {/if}
        </div>

    </div>
</div>