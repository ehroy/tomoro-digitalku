import { api } from './client.js';

export async function getVouchers(storeCode = '') {
	const params = new URLSearchParams();
	if (storeCode) {
		params.set('storeCode', storeCode);
	}
	const query = params.toString();
	return api.request(query ? `/vouchers?${query}` : '/vouchers');
}

export async function getMemberVouchers() {
	return api.request('/vouchers');
}
