import React from 'react';
import { Admin, Resource } from 'react-admin';
import KeyIcon from '@material-ui/icons/VpnKey';

import { PeopleList, PeopleEdit, PeopleCreate } from './people';
import PassBoard from './PassBoard';

import dataProvider from './dataProvider';

const App = () => (
    <Admin dataProvider={dataProvider}>
        <Resource name="peoples" list={PeopleList} edit={PeopleEdit} create={PeopleCreate} />
        <Resource name="pass" options={{ label: 'Passphrase'}} list={PassBoard} icon={KeyIcon} />
    </Admin>
);

export default App;
