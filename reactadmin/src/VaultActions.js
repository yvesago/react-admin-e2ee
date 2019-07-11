import React from 'react';
import { CreateButton, ExportButton } from 'react-admin';
import CardActions from '@material-ui/core/CardActions';
import VaultForgetAction from './VaultForgetAction';

const cardActionStyle = {
    zIndex: 2,
    display: 'inline-block',
    float: 'right',
};

// PeopleActions: a menu to show and forget current key
export const VaultActions = ({ resource, filters, displayedFilters, filterValues, basePath, showFilter, currentSort, exporter }) => (
    <CardActions style={cardActionStyle}>
        {filters && React.cloneElement(filters, { resource, showFilter, displayedFilters, filterValues, context: 'button' }) }
        <ExportButton
            resource={resource}
            sort={currentSort}
            filter={filterValues}
            exporter={exporter}
            maxResults={10000}
        />
        <CreateButton basePath={basePath} />
        <VaultForgetAction />
    </CardActions>
);

