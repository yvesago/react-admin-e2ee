import jsonServerProvider from './e2ec-ra-data-json-server';

import { AppConfig } from './AppConfig';

export default jsonServerProvider(AppConfig.API_URL);
