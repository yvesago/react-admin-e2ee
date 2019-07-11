import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import CardActions from '@material-ui/core/CardActions';
import { GET_LIST, UPDATE, GET_ONE, showNotification as showNotificationAction } from 'react-admin';
import { ViewTitle } from 'react-admin';
import { push as pushAction } from 'react-router-redux';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';

import dataProvider from './dataProvider';

import scrypt from 'scrypt-js';

import { text2ua, ua2text, encrypt_tob64string, decrypt_b64string  } from './E2Ecrypto';
import E2Econfig from './E2Econfig';

const formStyle = { padding: '0 1em 3em 1em' };


class VaultReencrypt extends Component {
    constructor() {
        super();
        this.handleSubmit = this.handleSubmit.bind(this);
        const oldkey = localStorage.getItem('keyb64');
        const vault = localStorage.getItem('vault');
        const vaultid = localStorage.getItem('vaultid');
        this.VaultName = vault;
        this.VaultId = vaultid;
        this.OldKey = oldkey;
        this.OldSalt = (oldkey) ? localStorage.getItem('salt').substring(0,16) : '';

        dataProvider(GET_ONE,'vault', {id: vaultid})
            .then((data) => {
                const currentsalt = data.data.verifykey.substring(0,16);
                const saltstring = localStorage.getItem('salt');
                const oldsalt = (saltstring) ? saltstring.substring(0,16) : '';
                if (currentsalt !== oldsalt) {
			        localStorage.removeItem('keyb64');
        			localStorage.removeItem('salt');
                    this.OldKey = '';
                    this.OldSalt = '';
                }
            });

        var total = {};
        E2Econfig.crypt_ressources.forEach(function(r) {
            dataProvider(GET_LIST,r, {pagination: [0,1],sort: 'id'}) 
                .then((data) => {
                    total[r] = data.total;
                });
        });
        this.Totals = total;
    }

    handleSubmit = (event) => {
        event.preventDefault();
        const record = new FormData(event.target);
        const { push, showNotification } = this.props;
        const pass = record.get('pass');

        const oldkey = this.OldKey;
        const oldsalt = this.OldSalt;
        const total = this.Totals;
        const vaultid = this.VaultId;
        const vaultname = this.VaultName;

        if (oldsalt === '') {
            showNotification('Error: no verif key stored ', 'warning');
            return;
        }

        // create new salt
        var saltb = new Uint8Array(12);
        window.crypto.getRandomValues(saltb);
        const salt = btoa(ua2text(saltb));
        // create key
        const N = 16384, r = 8, p = 1, dkLen = 32;
        scrypt(text2ua(pass), saltb, N, r, p, dkLen, function(error, progress, key) {
            if (error) {
                console.log('Error: ' + error);
                showNotification('Error: ' + error , 'warning');
                return;
            } else if (key) {
                const keyb64 = btoa(ua2text(key));
                console.log('keyb64: ' + keyb64);

                const  vrfystr = ' -- Set on '  + Date();

                localStorage.setItem('salt', salt + vrfystr );
                localStorage.setItem('keyb64', keyb64);

                // Update new verifkey
                dataProvider(UPDATE,'vault', 
                    { id: vaultid, 
                        data: { 
                            verifykey: salt + encrypt_tob64string(keyb64, vrfystr), 
                            vaultname: vaultname } 
                    }) 
                    .then(() => {
                        console.error('verif key stored');
                    })
                    .catch((e) => {
                        console.error(e);
                        showNotification('Error: no verif key stored ', 'warning');
                        return;
                    });

                showNotification('New key stored');
                console.log('new key');
                
                E2Econfig.crypt_ressources.forEach(function(r) {
                    dataProvider(GET_LIST,r, {pagination: [0,total[r]],sort: 'id'})  // TODO Paginate
                        .then((data) => {
                            const resp = data.data;
                            resp.forEach(function(o) {
                                console.log(o);
                                E2Econfig.crypt_fields.forEach(function(f) {
                                    if (o.hasOwnProperty(f) === true) {
                                        const field =  o[f];
                                        //console.log(f + ' => ' + field);
                                        if ( oldsalt === field.substring(0,16) ) {
                                            const oldField = decrypt_b64string(oldkey, field.substring(16));
                                            o[f] = oldField;
                                            // Update o via REST
                                            dataProvider(UPDATE,r,{ id: o.id, data: o } )
                                                .catch((e) => {
                                                    showNotification(e,'warning');
                                                    return;
                                                });
                                        }
                                    }
                                });
                                //console.log(o);
                            });
                        });
                });
                showNotification('All fields updated');
                push('/');
            }
        });
    }

    render() {
        if (this.OldSalt !== '') 
        {
            const name = 'Active Vault: "' + this.VaultName + '"';
            return (
                <CardActions>
                    <form onSubmit={this.handleSubmit}>
                        <ViewTitle title={name} />
                        <div style={formStyle}>
                            <TextField
                                label="New Passphrase"
                                name="pass"
                                required
                            />
                            <br />
                            <br />
                            <div>Could be very long...</div>
                            <br />
                            <Button type="submit" variant="raised" color="secondary">Apply</Button>
                        </div>
                    </form>
                </CardActions>
            );
        }
        else {
            return (
                <CardActions>
                    <div style={formStyle}>
                        <ViewTitle title="No active vault" />
                        Waiting to set the Passphrase
                    </div>
                </CardActions>
            );
        }
    }
}

VaultReencrypt.propTypes = {
    showNotification: PropTypes.func,
};

export default connect(null, {
    showNotification: showNotificationAction,
    push: pushAction,
})(VaultReencrypt);

