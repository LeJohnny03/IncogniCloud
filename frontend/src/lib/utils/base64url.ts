export function base64urlDecode(base64url: string): ArrayBuffer {
    const base64 = base64url.replace(/-/g, '+').replace(/_/g, '/');

    const padLen = (4 - (base64.length % 4)) % 4;
    const padded = base64 + '='.repeat(padLen);

    const binary = atob(padded);

    const bytes = new Uint8Array(binary.length);

    for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }

    return bytes.buffer;
}

export function base64urlEncode(buffer: ArrayBuffer): string {
    const bytes = new Uint8Array(buffer);

    let binary = '';

    for (let i = 0; i < bytes.length; i++) {
        binary += String.fromCharCode(bytes[i]);
    }

    return btoa(binary)
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=/g, '');
}