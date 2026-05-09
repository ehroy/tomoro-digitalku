import { writable } from 'svelte/store';

export const cart = writable([]);

function safeNumber(value, fallback = 0) {
	const parsed = Number(value);
	return Number.isFinite(parsed) ? parsed : fallback;
}

function pickCartItems(data) {
	if (!data) return [];
	if (Array.isArray(data)) return data;
	if (Array.isArray(data.cartItemList)) return data.cartItemList;
	if (Array.isArray(data.records)) return data.records;
	if (Array.isArray(data.cartItems)) return data.cartItems;
	if (Array.isArray(data.items)) return data.items;
	if (Array.isArray(data.list)) return data.list;
	if (Array.isArray(data.data)) return data.data;
	return [];
}

export function normalizeCartItem(item, fallback = {}) {
	const quantity = safeNumber(
		item?.itemAmount,
		safeNumber(item?.amount, safeNumber(item?.quantity, safeNumber(fallback.quantity, 1)))
	);
	const price = safeNumber(
		item?.pluPrice,
		safeNumber(item?.paymentPrice, safeNumber(item?.price, safeNumber(fallback.price, 0)))
	);
	const productId = item?.itemCode || item?.productCode || item?.code || fallback.product_id || '';
	const sizeCode = item?.pluCode || item?.sizeCode || fallback.size_code || '';

	return {
		cart_key: item?.key || item?.cartKey || item?.cart_key || item?.code || `${productId}:${sizeCode || 'default'}`,
		product_id: productId,
		name: item?.itemName || item?.productName || item?.name || fallback.name || '',
		size_label: item?.showPluName || item?.pluName || item?.sizeLabel || item?.size_name || fallback.size_label || '',
		size_code: sizeCode,
		quantity,
		price,
		subtotal: safeNumber(item?.pluTotalMoney, safeNumber(item?.subtotal, price * quantity))
	};
}

export function normalizeCartResponse(data, fallbackItems = []) {
	const items = pickCartItems(data);
	return items.map((item, index) => normalizeCartItem(item, fallbackItems[index] || {}));
}

export function setCartFromResponse(data, fallbackItems = [], preserveFallback = false) {
	const items = normalizeCartResponse(data, fallbackItems);
	const nextItems = preserveFallback && items.length === 0 && fallbackItems.length > 0
		? fallbackItems
		: items;

	cart.set(nextItems);
	return nextItems;
}

export function addToCart(product, quantity = 1) {
	cart.update(items => {
		const productId = product.product_id || product.id || product.code;
		const sizeCode = product.size_code || product.sizeCode || product.size || '';
		const cartKey = product.cart_key || `${productId}:${sizeCode || 'default'}`;
		const existingItem = items.find(item => item.cart_key === cartKey);
		
		if (existingItem) {
			existingItem.quantity += quantity;
			existingItem.subtotal = existingItem.price * existingItem.quantity;
			return [...items];
		}
		
		return [...items, {
			cart_key: cartKey,
			product_id: productId,
			name: product.name,
			size_label: product.size_label || product.sizeLabel || '',
			size_code: sizeCode,
			quantity,
			price: product.price,
			subtotal: product.price * quantity
		}];
	});
}

export function removeFromCart(cartKeyOrProductId) {
	cart.update(items => items.filter(item => item.cart_key !== cartKeyOrProductId && item.product_id !== cartKeyOrProductId));
}

export function updateCartQuantity(cartKeyOrProductId, quantity) {
	cart.update(items => {
		const item = items.find(item => item.cart_key === cartKeyOrProductId || item.product_id === cartKeyOrProductId);
		if (item) {
			item.quantity = quantity;
			item.subtotal = item.price * quantity;
		}
		return [...items];
	});
}

export function clearCart() {
	cart.set([]);
}

export function getCartTotal(items) {
	return items.reduce((total, item) => total + item.subtotal, 0);
}
