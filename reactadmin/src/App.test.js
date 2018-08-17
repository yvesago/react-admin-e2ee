import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { text2ua, ua2text, encrypt_tob64string, decrypt_b64string  } from './E2Ecrypto';
import scrypt from 'scrypt-js';
var assert = require('assert');

it('renders without crashing', () => {
    const div = document.createElement('div');
    ReactDOM.render(<App />, div);
    ReactDOM.unmountComponentAtNode(div);
});

it('crypto', () => {

  //var saltb = new Uint8Array(12);
  //window.crypto.getRandomValues(saltb); // XXX window could not be used in nodejs
  //var salt = btoa(ua2text(saltb));

  var salt = "TD+gfW2LANHIFf2+";
  var keyb64 = "j1QTSt34Sc9VYoOEuo0kzS6a1oOdxyWYAObphDy0JSU=";

  // Test scrypt
  var password = "test";
  var saltb = text2ua(atob(salt));
  const N = 1024, r = 8, p = 1, dkLen = 32;
  scrypt(text2ua(password), saltb, N, r, p, dkLen, function(error, progress, key) {
     if (key) {
      assert.equal(keyb64, btoa(ua2text(key)));
     };
     });

  //var verifkey = "TD+gfW2LANHIFf2+tF6YZb1pvAUMsanHPRsT9HjZ4p1hdtlMAQ3eicYFWc4VzCUgna9oCsFVNoVMMrSvFZAMjs/cn8sNuJeN+TvmnZc58Yz9lK8mE3xzrLAKtLsRh/2TzHu8UyScjV2G3jeBRqctxkQM8CtNYupEsEs=";
  //var salt = cypher.substring(0,16);
  //var enc = encrypt_tob64string(keyb64,"some text"); // XXX need window.crypto
  //console.log(enc);
  //var dec = decrypt_b64string(keyb64,cypher.substring(16)); // XXX TypeError: _libsodiumWrappers2.default.crypto_aead_chacha20poly1305_ietf_decrypt is not a function
  //console.log(dec);


});
