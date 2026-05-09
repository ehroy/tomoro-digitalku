import { api } from './client.js';

export async function calculateOrder(orderData) {
	return api.request('/orders/calculate', {
		method: 'POST',
		body: JSON.stringify(orderData)
	});
}

export async function applyVoucherToOrder(orderData) {
	return api.request('/orders/voucher/apply', {
		method: 'POST',
		body: JSON.stringify(orderData)
	});
}

export async function removeVoucherFromOrder(orderData) {
	return api.request('/orders/voucher/remove', {
		method: 'POST',
		body: JSON.stringify(orderData)
	});
}

export async function createOrder(orderData) {
	return api.request('/orders/create', {
		method: 'POST',
		body: JSON.stringify({
			...orderData,
			payment: 'gopay' // Auto-select GoPay
		})
	});
}

export async function getOrderHistory() {
	return api.request('/orders/history');
}

export async function getPayStatus(tradeOrderCode) {
	return api.request(`/orders/pay-status?tradeOrderCode=${encodeURIComponent(tradeOrderCode)}`);
}
