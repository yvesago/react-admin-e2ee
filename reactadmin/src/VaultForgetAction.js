import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { showNotification as showNotificationAction } from 'react-admin';
import { push as pushAction } from 'react-router-redux';
import Button from '@material-ui/core/Button';


class VaultForgetAction extends Component {
    constructor() {
        super();
        this.handleSubmit = this.handleSubmit.bind(this);
        const storedkey = localStorage.getItem('keyb64') || '';
        this.CurrentKey = storedkey;
        const vault = localStorage.getItem('vault') || '';
        this.CurrentVault = vault;
    }

    handleSubmit = (event) => {
        event.preventDefault();
        const { push, showNotification } = this.props;
        localStorage.removeItem('keyb64');
        localStorage.removeItem('salt');
        localStorage.removeItem('vault');
        localStorage.removeItem('vaultid');
        showNotification('Password forgeted');
        push('/');
    }

    render() {
        if (this.CurrentKey !== '') 
        {
            return (
                <span>&nbsp;&nbsp;&nbsp; <Button type="submit" variant="raised" color="secondary" onClick={this.handleSubmit}>Forget "{this.CurrentVault}" key</Button></span>
            );
        }
        else {
            return (<span />);	
        }
    }
}

VaultForgetAction.propTypes = {
//    push: PropTypes.func,
    showNotification: PropTypes.func,
};

export default connect(null, {
    showNotification: showNotificationAction,
    push: pushAction,
})(VaultForgetAction);

