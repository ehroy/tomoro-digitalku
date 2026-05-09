import { api } from './client.js';

export async function getMenu(storeCode) {
	const params = new URLSearchParams({ storeCode });
	return api.request(`/menu?${params}`);
}
