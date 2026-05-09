import { api } from './client.js';

export async function getStores(lat, lng, search = '', page = 1, size = 20) {
	const params = new URLSearchParams({
		lat: lat.toString(),
		lng: lng.toString(),
		search,
		page: page.toString(),
		size: size.toString()
	});
	
	return api.request(`/stores?${params}`);
}

export async function getStoreDetail(storeCode) {
	return api.request(`/stores/${storeCode}`);
}
