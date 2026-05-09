import { api } from './client.js';

export async function addToCart(item) {
	return api.request('/cart/add', {
		method: 'POST',
		body: JSON.stringify(item)
	});
}

export async function editCart(item) {
	return api.request('/cart/edit', {
		method: 'POST',
		body: JSON.stringify(item)
	});
}

export async function getCart(storeCode, mainMenuType = 1) {
	const params = new URLSearchParams({
		storeCode,
		mainMenuType: String(mainMenuType)
	});

	return api.request(`/cart?${params.toString()}`);
}

export async function getAllCart(storeCode) {
	const params = new URLSearchParams({
		storeCode
	});

	return api.request(`/cart/all?${params.toString()}`);
}
