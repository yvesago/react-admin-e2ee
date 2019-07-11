import React from 'react';
import CardActions from '@material-ui/core/CardActions';
import { List, Edit, Create, Datagrid, TextField, EditButton, 
    TabbedShowLayout, Tab,
    SimpleForm, TextInput } from 'react-admin';
import VaultAction from './VaultAction';
import VaultForgetAction from './VaultForgetAction';
import VaultReencrypt from './VaultReencrypt';

import { VaultActions } from './VaultActions'; // List actions


const cardActionStyle = {
    zIndex: 2,
    display: 'inline-block',
    float: 'right',
};

// Edit actions
const VaultShowActions = ({ basePath, data, ...props }) => (
    <CardActions style={cardActionStyle}>
        <VaultForgetAction { ...props } />
    </CardActions>
);

export const VaultList = (props) => (
    <List actions={<VaultActions />} {...props}>
        <Datagrid>
            <TextField source="vaultname" />
            <EditButton />
        </Datagrid>
    </List>
);

const VaultTitle = ({ record }) => {
    return <span>Vault {record ? `"${record.vaultname}"` : ''}</span>;
};

export const VaultEdit = (props) => (
    <Edit title={<VaultTitle />} actions={<VaultShowActions />} {...props}>
        <TabbedShowLayout>
            <Tab label="Set Passphrase">
                <VaultAction { ...props } />
            </Tab>
            <Tab label="Reencrypt">
                <VaultReencrypt { ...props } />
            </Tab>
        </TabbedShowLayout>
    </Edit>
);

export const VaultCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="vaultname" />
        </SimpleForm>
    </Create>
);
