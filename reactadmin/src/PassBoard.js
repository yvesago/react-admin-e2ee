import React from 'react';
import Card from '@material-ui/core/Card';
import { ViewTitle } from 'react-admin';
import PassAction from './PassAction';
import PassRemoveAction from './PassRemoveAction';
import PassReencrypt from './PassReencrypt';

export const PassBoard = ({ ...props }) => (
    <Card>
        <ViewTitle title="Passphrase" />
        <PassAction { ...props } />
        <ViewTitle title="Remove Passphrase" />
        <PassRemoveAction { ...props } />
        <ViewTitle title="Change Passphrase" />
        <PassReencrypt { ...props } />
    </Card>
);

export default PassBoard;
