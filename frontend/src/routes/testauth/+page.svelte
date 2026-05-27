<script lang="ts">
    import {
        base64urlDecode,
        base64urlEncode
    } from '$lib/utils/base64url';
    // simple reactive variable for username
    let regUsername = $state('');
    let regDisplayName = $state('');

    let authUsername = $state('');

    let result = $state('');

    const API_BASE = 'https://localhost:8080/api';

    async function register() {
        try {
            const beginResp = await fetch(`${API_BASE}/register/begin`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                    body: JSON.stringify({
                        username: regUsername,
                        display_name: regDisplayName
                    })
            });

            const options = await beginResp.json();

            options.publicKey.challenge = base64urlDecode(
                options.publicKey.challenge
            );

            options.publicKey.user.id = base64urlDecode(
                options.publicKey.user.id
            );

            const credential = (await navigator.credentials.create(
                options
            )) as PublicKeyCredential;

            const response =
                credential.response as AuthenticatorAttestationResponse;

            const attestationResponse = {
                id: credential.id,
                rawId: base64urlEncode(credential.rawId),
                type: credential.type,
                response: {
                    attestationObject: base64urlEncode(
                        response.attestationObject
                    ),
                    clientDataJSON: base64urlEncode(
                        response.clientDataJSON
                    )
                }
            };

            const finishResp = await fetch(`${API_BASE}/register/finish`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify(attestationResponse)
            });

            if (!finishResp.ok) {
                throw new Error('Registration failed');
            }

            result = 'Registration successful!';
        } catch (error) {
            result = `Registration failed: ${(error as Error).message}`;
            console.error(error);
        }
    }

    async function authenticate() {
        try {
            const beginResp = await fetch(`${API_BASE}/login/begin`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify({
                    username: authUsername
                })
            });

            const options = await beginResp.json();

            options.publicKey.challenge = base64urlDecode(
                options.publicKey.challenge
            );

            options.publicKey.allowCredentials =
                options.publicKey.allowCredentials.map(
                    (cred: PublicKeyCredentialDescriptor) => ({
                        ...cred,
                        id: base64urlDecode(cred.id as unknown as string)
                    })
                );

            const assertion = (await navigator.credentials.get(
                options
            )) as PublicKeyCredential;

            const response =
                assertion.response as AuthenticatorAssertionResponse;

            const assertionResponse = {
                id: assertion.id,
                rawId: base64urlEncode(assertion.rawId),
                type: assertion.type,
                response: {
                    authenticatorData: base64urlEncode(
                        response.authenticatorData
                    ),
                    clientDataJSON: base64urlEncode(
                        response.clientDataJSON
                    ),
                    signature: base64urlEncode(response.signature),
                    userHandle: response.userHandle
                        ? base64urlEncode(response.userHandle)
                        : null
                }
            };

            const finishResp = await fetch(`${API_BASE}/login/finish`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify(assertionResponse)
            });

            if (!finishResp.ok) {
                throw new Error('Authentication failed');
            }

            result = 'Login successful!';
        } catch (error) {
            result = `Login failed: ${(error as Error).message}`;
            console.error(error);
        }
    }
</script>

<svelte:head>
    <title>Passkey Authentication</title>
</svelte:head>

<div class="container">
    <h1>Passkey Authentication</h1>

    <section>
        <h2>Register</h2>

        <input
            bind:value={regUsername}
            type="text"
            placeholder="Username"
        />

        <input
            bind:value={regDisplayName}
            type="text"
            placeholder="Display Name"
        />

        <button onclick={register}>
            Register Passkey
        </button>
    </section>

    <section>
        <h2>Login</h2>

        <input
            bind:value={authUsername}
            type="text"
            placeholder="Username"
        />

        <button onclick={authenticate}>
            Login with Passkey
        </button>
    </section>

    <p>{result}</p>
</div>

<style>
    .container {
        max-width: 500px;
        margin: 2rem auto;
        display: flex;
        flex-direction: column;
        gap: 2rem;
        font-family: sans-serif;
    }

    section {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        padding: 1rem;
        border: 1px solid #ccc;
        border-radius: 8px;
    }

    input,
    button {
        padding: 0.75rem;
        font-size: 1rem;
    }

    button {
        cursor: pointer;
    }
</style>