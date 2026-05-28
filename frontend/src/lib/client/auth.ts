import { API_BASE } from '$lib/api/api';
import { base64urlDecode, base64urlEncode } from '$lib/utils/base64url';

export async function loginWithPasskey(username: string): Promise<void> {
    const beginResp = await fetch(`${API_BASE}/login/begin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ username })
    });

    if (!beginResp.ok) throw new Error("Benutzername ungültig oder Netzwerkfehler.");
    
    const options = await beginResp.json();
    options.publicKey.challenge = base64urlDecode(options.publicKey.challenge);
    options.publicKey.allowCredentials = options.publicKey.allowCredentials.map(
        (cred: any) => ({ ...cred, id: base64urlDecode(cred.id) })
    );

    const assertion = (await navigator.credentials.get(options)) as PublicKeyCredential;
    const response = assertion.response as AuthenticatorAssertionResponse;

    const assertionResponse = {
        id: assertion.id,
        rawId: base64urlEncode(assertion.rawId),
        type: assertion.type,
        response: {
            authenticatorData: base64urlEncode(response.authenticatorData),
            clientDataJSON: base64urlEncode(response.clientDataJSON),
            signature: base64urlEncode(response.signature),
            userHandle: response.userHandle ? base64urlEncode(response.userHandle) : null
        }
    };

    const finishResp = await fetch(`${API_BASE}/login/finish`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(assertionResponse)
    });

    if (!finishResp.ok) {
        throw new Error("Login wurde vom Server abgelehnt.");
    }
}

export async function registerWithPasskey(username: string, displayName: string): Promise<void> {
    // 1. Challenge für die Registrierung abholen
    const beginResp = await fetch(`${API_BASE}/register/begin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ username, display_name: displayName })
    });

    if (!beginResp.ok) throw new Error("Fehler beim Start der Registrierung.");
    
    const options = await beginResp.json();
    
    // Base64url Decoding für WebAuthn API (Create)
    options.publicKey.challenge = base64urlDecode(options.publicKey.challenge);
    options.publicKey.user.id = base64urlDecode(options.publicKey.user.id);
    if (options.publicKey.excludeCredentials) {
        options.publicKey.excludeCredentials = options.publicKey.excludeCredentials.map(
            (cred: any) => ({ ...cred, id: base64urlDecode(cred.id) })
        );
    }

    // 2. Browser WebAuthn API aufrufen, um Passkey zu ERSTELLEN
    const credential = (await navigator.credentials.create(options)) as PublicKeyCredential;
    const response = credential.response as AuthenticatorAttestationResponse;

    // 3. Daten für das Backend verpacken
    const attestationResponse = {
        id: credential.id,
        rawId: base64urlEncode(credential.rawId),
        type: credential.type,
        response: {
            clientDataJSON: base64urlEncode(response.clientDataJSON),
            attestationObject: base64urlEncode(response.attestationObject),
            // Transports sind wichtig für moderne Passkeys (z.B. Handy über Bluetooth)
            //transports: credential.response.getTransports ? credential.response.getTransports() : []
        }
    };

    // 4. Passkey an das Backend senden und Registrierung abschließen
    const finishResp = await fetch(`${API_BASE}/register/finish`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(attestationResponse)
    });

    if (!finishResp.ok) {
        const errData = await finishResp.json();
        throw new Error(errData.error || "Registrierung wurde vom Server abgelehnt.");
    }
}