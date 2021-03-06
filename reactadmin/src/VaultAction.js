import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import CardActions from '@material-ui/core/CardActions';
import { ViewTitle } from 'react-admin';
import { GET_ONE, UPDATE, showNotification as showNotificationAction } from 'react-admin';
import { push as pushAction } from 'react-router-redux';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';

import dataProvider from './dataProvider';

import scrypt from 'scrypt-js';

import { text2ua, ua2text, encrypt_tob64string, decrypt_b64string  } from './E2Ecrypto';


const formStyle = { padding: '0 1em 3em 1em' };


class VaultAction extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.vaultid = this.props.record.id;
        this.vaultname = this.props.record.vaultname;
        const storedkey = localStorage.getItem('keyb64');
        this.CurrentHash = (storedkey) ? 'Current Salt: ' + localStorage.getItem('salt') : '';
        this.vault = (storedkey) ? 'Active Vault: "' + localStorage.getItem('vault') + '"' : '';
    }

    handleSubmit = (event) => {
        event.preventDefault();
        const record = new FormData(event.target);
        const { push, showNotification } = this.props;
        const pass = record.get('pass');
        const vaultname = record.get('vaultname');
        const vaultid = this.vaultid;

        // Get salt server, verify key with decrypt string
        //  store salt to localStorage
        //  else create new salt and push to server with ecrypt string
        dataProvider(GET_ONE,'vault', {id: vaultid}) 
            .then((data) => {
                const resp = data.data;
                var saltb = new Uint8Array(12);
                var salt = '';
                if (resp.verifykey !== '') {
                    // extract salt from verifykey
                    salt = resp.verifykey.substring(0,16);
                    saltb = text2ua(atob(salt));
                }
                else {
                    // create salt for first key
                    window.crypto.getRandomValues(saltb);
                    salt = btoa(ua2text(saltb));
                }
                console.log('salt: ' + salt);

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

                        var vrfystr = '';
                        if (resp.verifykey !== '') {
                            vrfystr = decrypt_b64string(keyb64,resp.verifykey.substring(16));
                            if ( vrfystr === '')  {
                                showNotification('Error: wrong password', 'warning');
                                return;
                            }
                        }
                        else {
                            vrfystr = ' Set on '  + Date();
                        }

                        localStorage.setItem('salt', salt + vrfystr );
                        localStorage.setItem('keyb64', keyb64);
                        localStorage.setItem('vault', vaultname);
                        localStorage.setItem('vaultid', vaultid);

                        if (resp.verifykey === '')  // push new verifykey on server
                            dataProvider(UPDATE,'vault', 
                                { id: vaultid, // Useless ID
                                    data: { vaultname: vaultname, verifykey: salt + encrypt_tob64string(keyb64, vrfystr) } 
                                }) 
                                .then(() => {
                                    showNotification('verif key stored');
                                    push('/');
                                })
                                .catch((e) => {
                                    console.error(e);
                                    showNotification('Error: no verif key stored ', 'warning');
                                    return;
                                });

                        showNotification('key stored');
                        push('/');
                    }
                });
            })
            .catch((e) => {
                console.error(e);
                showNotification('Error: no verif key find ', 'warning');
            });

    }

    render() {
        if (this.CurrentHash === '') 
        {
            return (
                <CardActions>
                    <form onSubmit={this.handleSubmit}>
                        <div style={formStyle}>
                            <TextField label="Vault Name" name="vaultname" defaultValue={this.vaultname} required />
                            <br />
                            <TextField
                                label="Set Passphrase"
                                name="pass"
                                required
                            />
                            <br />
                            <br />
                            <Button type="submit" variant="raised" color="secondary">Save</Button>
                        </div>
                    </form>
                </CardActions>
            );
        }
        else {
            return (
                <CardActions>
                    <form>
                    <ViewTitle title={this.vault} />
                    <div style={formStyle}>{this.CurrentHash}</div>
                    </form>
                </CardActions>
            );
        }
    }
}

VaultAction.propTypes = {
    record: PropTypes.object,
    showNotification: PropTypes.func,
};

export default connect(null, {
    showNotification: showNotificationAction,
    push: pushAction,
})(VaultAction);

