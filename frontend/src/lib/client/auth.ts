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