import sodium from 'libsodium-wrappers';
import E2Econfig from './E2Econfig';


function text2ua(s) {
    var ua = new Uint8Array(s.length);
    for (var i = 0; i < s.length; i++) {
        ua[i] = s.charCodeAt(i);
    }
    return ua;
}

function ua2text(ua) {
    var s = '';
    for (var i = 0; i < ua.length; i++) {
        s += String.fromCharCode(ua[i]);
    }
    return s;
}


function decrypt_b64string(keyb64, str) {
    const nonceb64 = str.substring(0,16);
    const txt = str.substring(16);
    if (keyb64 === '') {
        return '*** Encrypted field ***';
    }

    const keyb = text2ua(atob(keyb64));
    const nonce = text2ua(atob(nonceb64));
    var decoded = '';
    try {
        decoded = sodium.crypto_aead_chacha20poly1305_ietf_decrypt(
            null,
            text2ua(atob(txt)),
            nonceb64,
            nonce,
            keyb,
            'text'
        );
    }
    catch (e) {
        console.error(e);
        // Custom some error messages
        if ( e.name === 'TypeError') {
            throw new Error('Reload error');
        } else if (e.message === 'invalid usage')  { 
            throw new Error('Error : Bad Key (maybe key has changed)');
        }
        throw(e);
    }
    return decoded;
}

function encrypt_tob64string(keyb64,txt) {
    const keyb = text2ua(atob(keyb64));

    var nonce = new Uint8Array(12);
    window.crypto.getRandomValues(nonce); // New nonce for each encrypted field
    const nonceb64 = btoa(ua2text(nonce));

    const enctxt = sodium.crypto_aead_chacha20poly1305_ietf_encrypt(
        txt,
        nonceb64, // nonceb64 used as AAD
        undefined,
        nonce,
        keyb
    );
    // First 16 char are the base64 nonce
    return nonceb64 + btoa(ua2text(enctxt));
}


const Decrypt = function(json) {
    console.log('==DeCrypt==');
    console.log(json);
    const keyb64 = localStorage.getItem('keyb64')||'';
    
    E2Econfig.crypt_fields.forEach(function(f) {
        if (json.hasOwnProperty(f) === true) {
            const field =  json[f];
            if ( field !== '') {
                json[f] = decrypt_b64string(keyb64, field.substring(16)); // ignore first 16 chars salt
            }
        }
    });
    return json;
};

const Encrypt = function(json) {
    console.log('==Crypt==');
    const salt = localStorage.getItem('salt').substring(0,16); // || '';
    const keyb64 = localStorage.getItem('keyb64'); // || '';
    E2Econfig.crypt_fields.forEach(function(f) {
        if (json.hasOwnProperty(f) === true) {
            const field =  json[f];
            if ( field !== '') {
                if ( salt !== '' && keyb64 !== '') {
                    json[f] = salt + encrypt_tob64string(keyb64,field); // TODO verify max length
                    console.log(salt);
                    console.log(json[f]);
                }
            }
        }
    });
    return json;
};

export { Encrypt, Decrypt, encrypt_tob64string, decrypt_b64string, text2ua, ua2text };
