async function webAuhtnAPI(url, username) {
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ username})
    });
    if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return await response.json();
}

async function register() {
    console.log("Register Start")
    const username = document.getElementById('registerUsername').value;
    if (!username) {
        alert('Username is required');
        return
    }

    try {
        const optionsResponse = await webAuhtnAPI('/register/begin', username)

        const publicKeyCredentialCreationOptions = {
            ...optionsResponse,
            challenge: base64URLToUint8Array(optionsResponse.challenge),
            allowCredentials: optionsResponse.allowCredentials.map(cred => ({
                ...cred,
                id: base64URLToUint8Array(cred.id)
            }))
        }

        const credential = await navigator.credentials.get({
            publicKey: publicKeyCredentialCreationOptions
        })

        const registrationResponse = await fetch('/register/finish', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                id: credential.id,
                rawId: credential.id,
                type: credential.type,
                response: {
                    attestationObject: arrayBufferToBase64URL(credential.response.attestationObject),
                    clientDataJSON: arrayBufferToBase64URL(credential.response.clientDataJSON)
                }
            })
        })
        if (registrationResponse.ok) {
            alert('Registration successful')
        } else {
            alert('Registration failed')
        }
    } catch (error) {
        console.error(error)
        alert('Registration failed')
    }
}

async function login() {
    const username = document.getElementById('loginUsername').value;
    if (!username) {
        alert('Username is required');
        return;
    }

    try {
        // ログインの開始
        const optionsResponse = await webAuthnAPI('/login/begin', username);
        
        // 公開鍵の取得オプションを取得
        const publicKeyCredentialRequestOptions = {
            ...optionsResponse,
            challenge: base64URLToUint8Array(optionsResponse.challenge),
            allowCredentials: optionsResponse.allowCredentials.map(cred => ({
                ...cred,
                id: base64URLToUint8Array(cred.id),
            })),
        };

        // ブラウザのCredentials APIを使用して認証情報を取得
        const credential = await navigator.credentials.get({
            publicKey: publicKeyCredentialRequestOptions
        });

        // 取得された認証情報をサーバーに送信
        const loginResponse = await fetch('/login/finish', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                username,
                id: credential.id,
                rawId: arrayBufferToBase64URL(credential.rawId),
                type: credential.type,
                response: {
                    authenticatorData: arrayBufferToBase64URL(credential.response.authenticatorData),
                    clientDataJSON: arrayBufferToBase64URL(credential.response.clientDataJSON),
                    signature: arrayBufferToBase64URL(credential.response.signature),
                    userHandle: credential.response.userHandle ? arrayBufferToBase64URL(credential.response.userHandle) : null,
                },
            }),
        });

        if (loginResponse.ok) {
            alert('Login successful!');
        } else {
            throw new Error('Login failed');
        }
    } catch (error) {
        console.error('Error during login:', error);
        alert('Login failed: ' + error.message);
    }
}

function base64URLToUint8Array(base64URL) {
    const padding = '='.repeat((4 - base64URL.length % 4) % 4);
    const base64 = (base64URL + padding)
        .replace(/-/g, '+')
        .replace(/_/g, '/');
    const rawData = window.atob(base64);
    const outputArray = new Uint8Array(rawData.length);
    for (let i = 0; i < rawData.length; ++i) {
        outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
}

function arrayBufferToBase64URL(buffer) {
    const base64 = btoa(String.fromCharCode(...new Uint8Array(buffer)));
    return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
}
