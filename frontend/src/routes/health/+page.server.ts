// frontend/src/routes/health/+page.server.ts

import { API_BASE } from '$lib/api/api';

// Wir definieren das TypeScript-Interface exakt so wie das Go-Struct!
// Das ist unser Ersatz für OpenAPI/gRPC im klassischen REST.

interface HealthResponse {
    status: string;
    message: string;
}

export const load = async ({ fetch }) => {
    // Rufe das Go-Backend auf
    const response = await fetch(`${API_BASE}/health`);

    console.log("Response vom Backend:", response);
    
    if (!response.ok) {
        throw new Error("Backend nicht erreichbar");
    }

    // Wandle die Antwort in unser Interface um
    const data: HealthResponse = await response.json();

    // Alles, was du hier zurückgibst, steht der Svelte-Komponente zur Verfügung
    return {
        backendData: data
    };
};