# PROTOCOL

**scrypt** is used to protect the clear text master key.

**chacha20poly1305** crypts data (_chacha20_) and authenticate (_poly1305_) that the master key could safely decrypt the cypher text.


### Data format

Encrypted fields contain base64 strings.

- First 16 chars are the salt of scrypt key
    - mostly used to quickly verifiy that field is encrypted with the current master key
    - each string could be decoded by third party knowing the clear text master key
- Second 16 chars are the nonce of chacha20poly1305
    - nonce changes for each cypher text
- Next chars are the cypher text



### Master key verification


``GET /verifkey ``
return a string with ``[16charsSalt][16charsNonce][Free Text like creation date]``

Salt is used to create the master binary key and localy verify that the clear text master key is the current shared key.
The master binary key is stored on local browser to crypt and decrypt fields.

``PUT /verifkey``
store a new master key salt.
All encrypted fields are reencoded with the new key.


### Threat model

- Shared clear master key MUST be safely distributed between allowed users.
- https is MANDATORY to avoid MITM.
- Users MUST be authenticated 
    - to avoid public use of ``/verifkey`` and (very long) brute force of scrypt key.
    - to manage read, write, deletion rights
- Attacker need to stole database AND, at least, a binary key.
- The local binary key COULD be erased on user logout or by an user action.
- On suspicious activity master key MUST be changed.
- Master key COULD be regulary changed.

