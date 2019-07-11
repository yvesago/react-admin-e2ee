import React from 'react';
import { List, Edit, Create, Datagrid, TextField, EditButton,
    DisabledInput, LongTextInput, DateInput,
    SimpleForm, TextInput } from 'react-admin';
import { VaultActions } from './VaultActions';

export const PeopleList = (props) => (
    <List actions={<VaultActions />} {...props}>
        <Datagrid>
            <TextField source="id" />
            <TextField source="name" />
            <EditButton />
        </Datagrid>
    </List>
);

const PeopleTitle = ({ record }) => {
    return <span>People {record ? `"${record.name}"` : ''}</span>;
};

export const PeopleEdit = (props) => (
    <Edit title={<PeopleTitle />} {...props}>
        <SimpleForm>
            <DisabledInput source="id" />
            <TextInput source="name" />
            <LongTextInput source="xaddress" />
            <DateInput label="Date of birth" source="xdob" />
        </SimpleForm>
    </Edit>
);

export const PeopleCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="name" />
            <LongTextInput source="xaddress" />
        </SimpleForm>
    </Create>
);
