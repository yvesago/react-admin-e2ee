import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import CardActions from '@material-ui/core/CardActions';
import { showNotification as showNotificationAction } from 'react-admin';
import { push as pushAction } from 'react-router-redux';
import Button from '@material-ui/core/Button';


const formStyle = { padding: '0 1em 3em 1em' };

class RemovePassAction extends Component {
    constructor() {
        super();
        this.handleSubmit = this.handleSubmit.bind(this);
        const storedkey = localStorage.getItem('keyb64') || '';
        this.CurrentKey = storedkey;
    }

    handleSubmit = (event) => {
        event.preventDefault();
        const { push, showNotification } = this.props;
        localStorage.removeItem('keyb64');
        localStorage.removeItem('salt');
        showNotification('Password removed');
        push('/');
    }

    render() {
        if (this.CurrentKey !== '') 
        {
            return (
                <CardActions>
                    <form onSubmit={this.handleSubmit}>
                        <div style={formStyle}>
                            <Button type="submit" variant="raised" color="primary">Remove</Button>
                        </div>
                    </form>
                </CardActions>
            );
        }
        else {
            return (
                <CardActions>
                    <form><div style={formStyle}>Waiting for current passphrase</div></form>
                </CardActions>
            );	
        }
    }
}

RemovePassAction.propTypes = {
//    push: PropTypes.func,
    showNotification: PropTypes.func,
};

export default connect(null, {
    showNotification: showNotificationAction,
    push: pushAction,
})(RemovePassAction);

