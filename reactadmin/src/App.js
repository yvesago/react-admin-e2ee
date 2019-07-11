import React from 'react';
import { Admin, Resource } from 'react-admin';
import KeyIcon from '@material-ui/icons/VpnKey';

import { PeopleList, PeopleEdit, PeopleCreate } from './people';
import { VaultList, VaultEdit, VaultCreate } from './vault';


import dataProvider from './dataProvider';

const App = () => (
    <Admin dataProvider={dataProvider}>
        <Resource name="peoples" list={PeopleList} edit={PeopleEdit} create={PeopleCreate} />
        <Resource name="vault" options={{ label: 'Vaults'}} list={VaultList} edit={VaultEdit} create={VaultCreate} icon={KeyIcon} />
    </Admin>
);

export default App;
