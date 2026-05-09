<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { cart, clearCart, getCartTotal, setCartFromResponse } from '$lib/stores/cart.js';
	import { authStore } from '$lib/stores/auth.js';
	import { selectedStoreStore, voucherStore } from '$lib/stores/voucher.js';
	import { getCart } from '$lib/api/cart.js';
	import { getVouchers } from '$lib/api/vouchers.js';
	import { applyVoucherToOrder, calculateOrder, createOrder, removeVoucherFromOrder } from '$lib/api/orders.js';
	import UiIcon from '$lib/components/UiIcon.svelte';

	let selectedStore = null;
	let vouchers = [];
	let selectedVoucher = null;
	let pendingVoucher = null;
	let loading = false;
	let loadingVouchers = true;
	let loadingCalculation = false;
	let showVoucherList = false;
	let calculatedTotal = null;
	let orderSuccess = false;

	$: cartTotal = getCartTotal($cart);
	$: orderSubtotal = safeMoney(calculatedTotal?.itemTotalMoney, cartTotal);
	$: discount = safeMoney(calculatedTotal?.discountMoney, safeMoney(calculatedTotal?.discountAmount, 0));
	$: finalAmount = safeMoney(calculatedTotal?.paymentMoney, orderSubtotal - discount);
	$: cartItems = $cart.map((item) => ({ ...item, quantity: safeMoney(item?.quantity, safeMoney(item?.amount, 1)) }));

	function safeMoney(value, fallback = 0) {
		const parsed = Number(value);
		return Number.isFinite(parsed) ? parsed : fallback;
	}

	function buildCalculationPayload(businessCode = '') {
		const payload = {
			storeCode: selectedStore.storeCode,
			deliverType: 1,
			bizChannel: 'APP'
		};

		if (businessCode) {
			payload.businessCode = businessCode;
		}

		return payload;
	}

	function buildRemoveVoucherPayload() {
		return {
			storeCode: selectedStore.storeCode,
			deliverType: 1,
			bizChannel: 'APP'
		};
	}

	function buildApplyVoucherPayload(businessCode = '') {
		const payload = buildCalculationPayload(businessCode);
		payload.businessType = 1;
		return payload;
	}

	function getVoucherBusinessCode(voucher) {
		return voucher?.accountCouponCode || voucher?.memberCouponCode || voucher?.couponCode || '';
	}

	function normalizeVoucherList(data) {
		if (Array.isArray(data)) return data;
		if (Array.isArray(data?.coupons)) return data.coupons;
		if (Array.isArray(data?.records)) return data.records;
		if (Array.isArray(data?.couponList)) return data.couponList;
		return [];
	}

	function normalizeVoucherItem(voucher) {
		return {
			...voucher,
			couponName: voucher?.couponName || voucher?.name || voucher?.title || 'Voucher',
			couponCode: voucher?.couponCode || voucher?.code || '',
			accountCouponCode: voucher?.accountCouponCode || voucher?.memberCouponCode || '',
			discountValue: safeMoney(voucher?.discountValue, safeMoney(voucher?.discountEffectivePrice, 0)),
			minOrderAmount: safeMoney(voucher?.minOrderAmount, 0)
		};
	}

	function normalizeOrderSummary(data) {
		const subtotal = safeMoney(data?.itemTotalMoney, safeMoney(data?.itemTotalLineMoney, cartTotal));
		const discountValue = safeMoney(data?.discountMoney, safeMoney(data?.discountAmount, 0));
		const total = safeMoney(data?.paymentMoney, subtotal - discountValue);

		return {
			...data,
			itemTotalMoney: subtotal,
			discountMoney: discountValue,
			paymentMoney: total,
			finalAmount: total,
			totalAmount: total
		};
	}

	async function refreshCartFromServer() {
		const cartResponse = await getCart(selectedStore.storeCode, 1);
		setCartFromResponse(cartResponse.data, $cart, true);
	}

	async function loadOrderSummary() {
		loadingCalculation = true;
		try {
			await refreshCartFromServer();
			const response = await calculateOrder(buildCalculationPayload());
			calculatedTotal = normalizeOrderSummary(response.data);
		} catch (error) {
			console.error('Failed to calculate order:', error);
		} finally {
			loadingCalculation = false;
		}
	}

	async function applyOrderVoucher(businessCode = '') {
		loadingCalculation = true;
		try {
			await refreshCartFromServer();
			const response = await applyVoucherToOrder(buildApplyVoucherPayload(businessCode));
			calculatedTotal = normalizeOrderSummary(response.data);
		} catch (error) {
			console.error('Failed to calculate order:', error);
		} finally {
			loadingCalculation = false;
		}
	}

	onMount(async () => {
		selectedStore = $selectedStoreStore;
		
		if (!selectedStore || $cart.length === 0) {
			goto('/menu');
			return;
		}

		try {
			await refreshCartFromServer();
			await loadOrderSummary();
		} catch (error) {
			console.error('Failed to load cart:', error);
		}

		// Load available vouchers
		try {
			const response = await getVouchers(selectedStore.storeCode);
			vouchers = normalizeVoucherList(response.data).map(normalizeVoucherItem);
		} catch (error) {
			console.error('Failed to load vouchers:', error);
		} finally {
			loadingVouchers = false;
		}
	});

	async function handleSelectVoucher(voucher) {
		const normalized = normalizeVoucherItem(voucher);
		selectedVoucher = normalized;
		pendingVoucher = normalized;
		voucherStore.set(normalized);
		showVoucherList = false;
		await applyOrderVoucher(getVoucherBusinessCode(normalized));
	}

	async function applyPendingVoucher() {
		if (!pendingVoucher) return;
		selectedVoucher = pendingVoucher;
		voucherStore.set(pendingVoucher);
		showVoucherList = false;
		await applyOrderVoucher(getVoucherBusinessCode(pendingVoucher));
	}

	async function removeVoucher() {
		selectedVoucher = null;
		pendingVoucher = null;
		voucherStore.set(null);
		showVoucherList = true;
		loadingCalculation = true;
		try {
			await refreshCartFromServer();
			const response = await removeVoucherFromOrder(buildRemoveVoucherPayload());
			calculatedTotal = normalizeOrderSummary(response.data);
		} catch (error) {
			console.error('Failed to calculate order:', error);
		} finally {
			loadingCalculation = false;
		}
	}

	async function handleCreateOrder() {
		loading = true;
		
		try {
			const userCode = $authStore.user?.memberCode || $authStore.user?.code || '';
			const userName = $authStore.user?.name || '';
			const orderData = {
				storeCode: selectedStore.storeCode,
				businessType: 1,
				deliverType: 1,
				bizChannel: 'APP',
				payChannelCode: '1007',
				userCode,
				userName,
			};

			const voucherCode = getVoucherBusinessCode(selectedVoucher);
			if (voucherCode) {
				orderData.businessCode = voucherCode;
			}
			
			const response = await createOrder(orderData);
			const checkoutUrl = response.data?.mobileDeeplinkCheckoutUrl || response.data?.walletsChargesResultVo?.mobileDeeplinkCheckoutUrl || '';
			const guidePicture = response.data?.guidePicture || response.data?.walletsChargesResultVo?.guidePicture || 'https://static.tomoro-coffee.id/prod/material20220721/Payment/wait1/t.png';
			const tradeOrderCode = response.data?.code || response.data?.orderCode || response.data?.paymentOrderCode || response.data?.transCode || '';
			const statusState = {
				checkoutUrl,
				guidePicture,
				tradeOrderCode,
				orderCode: tradeOrderCode,
				orderNo: tradeOrderCode,
				storeCode: selectedStore?.storeCode || ''
			};
			sessionStorage.setItem('checkoutStatusState', JSON.stringify(statusState));
			const params = new URLSearchParams();
			if (checkoutUrl) params.set('checkoutUrl', checkoutUrl);
			if (guidePicture) params.set('guidePicture', guidePicture);
			if (tradeOrderCode) params.set('tradeOrderCode', tradeOrderCode);
			if (tradeOrderCode) params.set('orderCode', tradeOrderCode);
			if (tradeOrderCode) params.set('orderNo', tradeOrderCode);
			if (selectedStore?.storeCode) params.set('storeCode', selectedStore.storeCode);

			clearCart();
			goto(`/checkout/status?${params.toString()}`);
		} catch (error) {
			console.error('Failed to create order:', error);
			alert('Gagal membuat pesanan. Silakan coba lagi.');
		} finally {
			loading = false;
		}
	}

	function formatPrice(price) {
		return new Intl.NumberFormat('id-ID', {
			style: 'currency',
			currency: 'IDR',
			minimumFractionDigits: 0
		}).format(price);
	}
</script>

<svelte:head>
	<title>Checkout - Coffee Order</title>
</svelte:head>

<div class="checkout-screen">
	{#if orderSuccess}
		<div class="success-overlay">
			<div class="success-card">
				<div class="success-icon"><UiIcon name="check" size={32} /></div>
				<h2>Pesanan Berhasil!</h2>
				<p>Pesanan Anda sedang diproses</p>
			</div>
		</div>
	{/if}

		<div class="checkout-header">
			<button class="back-btn" on:click={() => goto('/menu')}>
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
					<path d="m15 18-6-6 6-6"/>
				</svg>
			</button>
			<h1>Checkout</h1>
		</div>

	<div class="checkout-content">
		<!-- Store Info -->
		<div class="section-card">
			<h3>Outlet</h3>
			<div class="store-info">
				<div class="store-icon"><UiIcon name="coffee" size={20} /></div>
				<div>
					<div class="store-name">{selectedStore?.storeName}</div>
					<div class="store-address">{selectedStore?.storeAddress}</div>
				</div>
			</div>
		</div>

		<!-- Order Items -->
		<div class="section-card">
			<h3>Pesanan Anda</h3>
			<div class="order-items">
				{#each cartItems as item}
					<div class="order-item">
						<div class="item-details">
							<span class="item-name">{item.name}</span>
							{#if item.size_label}
								<span class="item-size">{item.size_label}</span>
							{/if}
							<span class="item-qty">x{item.quantity}</span>
						</div>
						<span class="item-price">{formatPrice(item.subtotal)}</span>
					</div>
				{/each}
			</div>
		</div>

		<!-- Voucher Selection -->
		<div class="section-card">
			<h3>Voucher</h3>
			
		{#if selectedVoucher}
				<div class="selected-voucher">
					<div class="voucher-info">
						<div class="voucher-icon"><UiIcon name="voucher" size={20} /></div>
						<div>
							<div class="voucher-name">{selectedVoucher.couponName}</div>
							<div class="voucher-discount">
								Diskon {formatPrice(safeMoney(selectedVoucher.discountValue, safeMoney(selectedVoucher.discountEffectivePrice, 0))) }
							</div>
							<div class="voucher-picked">Voucher aktif</div>
						</div>
					</div>
					<button class="remove-voucher-btn" on:click={removeVoucher}><UiIcon name="close" size={16} /></button>
				</div>
			{:else}
				<button class="select-voucher-btn" on:click={() => showVoucherList = !showVoucherList}>
					{showVoucherList ? 'Tutup Daftar Voucher' : 'Pilih Voucher'}
				</button>
			{/if}

			{#if pendingVoucher && pendingVoucher !== selectedVoucher}
				<div class="voucher-preview">
					<div>
						<div class="voucher-name">{pendingVoucher.couponName || 'Voucher dipilih'}</div>
						<div class="voucher-discount">
							Diskon {formatPrice(safeMoney(pendingVoucher.discountValue, safeMoney(pendingVoucher.discountEffectivePrice, 0)))}
						</div>
					</div>
					<button class="apply-voucher-btn" on:click={applyPendingVoucher} disabled={loadingCalculation}>
						{#if loadingCalculation}
							Memproses...
						{:else}
							Pakai Voucher
						{/if}
					</button>
				</div>
			{/if}

			{#if !selectedVoucher && vouchers.length > 0 && showVoucherList}
				<div class="voucher-list">
					{#if loadingVouchers}
						<div class="loading-vouchers">Loading vouchers...</div>
					{:else}
					{#each vouchers as voucher}
						<div 
							class="voucher-card"
							on:click={() => handleSelectVoucher(voucher)}
								on:keypress={() => handleSelectVoucher(voucher)}
								role="button"
								tabindex="0"
							>
								<div class="voucher-icon"><UiIcon name="voucher" size={20} /></div>
							<div class="voucher-details">
								<div class="voucher-name">{voucher.couponName}</div>
								<div class="voucher-discount">
									Diskon {formatPrice(safeMoney(voucher.discountValue, safeMoney(voucher.discountEffectivePrice, 0)))}
								</div>
								{#if voucher.minOrderAmount}
										<div class="voucher-min">
											Min. pembelian {formatPrice(voucher.minOrderAmount)}
										</div>
									{/if}
									{#if pendingVoucher === voucher}
										<div class="voucher-picked">Dipilih</div>
									{/if}
								</div>
							</div>
						{/each}
					{/if}
				</div>
			{:else if !selectedVoucher && !loadingVouchers && vouchers.length === 0}
				<div class="empty-voucher-state">
					<UiIcon name="voucher" size={18} />
					<span>Voucher tidak tersedia untuk outlet ini</span>
				</div>
			{/if}
		</div>

		<!-- Payment Method -->
		<div class="section-card">
			<h3>Metode Pembayaran</h3>
			<div class="payment-method">
				<div class="payment-icon"><UiIcon name="voucher" size={22} /></div>
				<div>
					<div class="payment-name">GoPay</div>
					<div class="payment-desc">Pembayaran via GoPay</div>
				</div>
				<div class="payment-badge">Terpilih</div>
			</div>
		</div>

		<!-- Order Summary -->
		<div class="section-card summary-card">
			<h3>Ringkasan Pembayaran</h3>
			<div class="summary-row">
				<span>Subtotal</span>
				<span>{formatPrice(orderSubtotal)}</span>
			</div>
			{#if discount > 0}
				<div class="summary-row discount">
					<span>Diskon</span>
					<span>-{formatPrice(discount)}</span>
				</div>
			{/if}
			<div class="summary-divider"></div>
			<div class="summary-row total">
				<span>Total</span>
				<span>{formatPrice(finalAmount)}</span>
			</div>
		</div>
	</div>

	<div class="checkout-footer">
		<div class="footer-total">
			<span>Total Pembayaran</span>
			<span class="footer-amount">{formatPrice(finalAmount)}</span>
		</div>
		<button 
			class="order-btn" 
			disabled={loading}
			on:click={handleCreateOrder}
		>
			{#if loading}
				<span class="spinner"></span>
				Memproses...
			{:else}
				Buat Pesanan
			{/if}
		</button>
	</div>
</div>

<style>
	.checkout-screen {
		display: flex;
		flex-direction: column;
		min-height: calc(100vh - 68px);
		background: var(--warm-gray);
	}

	.success-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.8);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		animation: fadeIn 0.3s ease;
	}

	.success-card {
		background: white;
		border-radius: 20px;
		padding: 40px;
		text-align: center;
		animation: scaleIn 0.3s ease;
	}

	@keyframes scaleIn {
		from {
			transform: scale(0.8);
			opacity: 0;
		}
		to {
			transform: scale(1);
			opacity: 1;
		}
	}

	.success-icon {
		font-size: 64px;
		margin-bottom: 16px;
	}

	.success-card h2 {
		margin: 0 0 8px 0;
		color: var(--text-dark);
		font-size: 24px;
	}

	.success-card p {
		margin: 0;
		color: var(--text-mid);
	}

	.checkout-header {
		background: white;
		padding: 16px;
		display: flex;
		align-items: center;
		gap: 12px;
		border-bottom: 1px solid #f0f0f0;
	}

	.back-btn {
		background: var(--warm-gray);
		border: none;
		width: 40px;
		height: 40px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
	}

	.checkout-header h1 {
		margin: 0;
		font-size: 20px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.checkout-content {
		flex: 1;
		overflow-y: auto;
		padding: 16px 12px 100px 12px;
		scrollbar-width: none;
	}

	.checkout-content::-webkit-scrollbar {
		display: none;
	}

	.section-card {
		background: white;
		border-radius: 16px;
		padding: 16px;
		margin-bottom: 12px;
	}

	.section-card h3 {
		margin: 0 0 12px 0;
		font-size: 16px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.store-info {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.store-icon {
		width: 48px;
		height: 48px;
		background: var(--orange-light);
		border-radius: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 24px;
		flex-shrink: 0;
	}

	.store-name {
		font-weight: 700;
		color: var(--text-dark);
		margin-bottom: 4px;
	}

	.store-address {
		font-size: 12px;
		color: var(--text-light);
	}

	.order-items {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.order-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 12px;
		background: var(--warm-gray);
		border-radius: 10px;
	}

	.item-details {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.item-name {
		font-weight: 700;
		color: var(--text-dark);
	}

	.item-size {
		font-size: 11px;
		color: var(--text-light);
		font-weight: 600;
	}

	.item-qty {
		background: var(--orange);
		color: white;
		padding: 2px 8px;
		border-radius: 8px;
		font-size: 12px;
		font-weight: 700;
	}

	.item-price {
		font-weight: 800;
		color: var(--orange);
	}

	.select-voucher-btn {
		width: 100%;
		padding: 12px;
		background: var(--orange-light);
		color: var(--orange);
		border: 2px dashed var(--orange);
		border-radius: 12px;
		font-weight: 700;
		cursor: pointer;
		font-family: var(--font);
		transition: all 0.2s;
	}

	.select-voucher-btn:active {
		transform: scale(0.98);
	}

	.selected-voucher {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 12px;
		background: var(--orange-light);
		border-radius: 12px;
		border: 2px solid var(--orange);
	}

	.voucher-preview {
		margin-top: 10px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		padding: 12px;
		border-radius: 12px;
		background: #fff8f1;
		border: 1px solid rgba(244, 147, 106, 0.35);
	}

	.apply-voucher-btn {
		border: none;
		background: var(--orange);
		color: white;
		border-radius: 10px;
		padding: 10px 14px;
		font-weight: 700;
		cursor: pointer;
		font-family: var(--font);
		white-space: nowrap;
	}

	.apply-voucher-btn:disabled {
		opacity: 0.7;
		cursor: not-allowed;
	}

	.voucher-picked {
		margin-top: 6px;
		display: inline-flex;
		align-items: center;
		padding: 3px 8px;
		border-radius: 999px;
		background: rgba(244, 147, 106, 0.16);
		color: var(--orange);
		font-size: 11px;
		font-weight: 700;
	}

	.voucher-info {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.voucher-icon {
		font-size: 24px;
	}

	.voucher-name {
		font-weight: 700;
		color: var(--text-dark);
		margin-bottom: 2px;
	}

	.voucher-discount {
		font-size: 12px;
		color: var(--orange);
		font-weight: 700;
	}

	.remove-voucher-btn {
		background: none;
		border: none;
		color: var(--text-light);
		font-size: 20px;
		cursor: pointer;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		transition: all 0.2s;
	}

	.remove-voucher-btn:hover {
		background: rgba(0, 0, 0, 0.1);
	}

	.voucher-list {
		margin-top: 12px;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.empty-voucher-state {
		margin-top: 12px;
		padding: 14px;
		border-radius: 12px;
		background: #fafafa;
		color: var(--text-light);
		display: flex;
		align-items: center;
		gap: 10px;
		font-size: 13px;
		font-weight: 600;
	}

	.voucher-card {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px;
		background: var(--warm-gray);
		border-radius: 12px;
		cursor: pointer;
		transition: all 0.2s;
		border: 2px solid transparent;
	}

	.voucher-card:active {
		transform: scale(0.98);
		border-color: var(--orange);
	}

	.voucher-details {
		flex: 1;
	}

	.voucher-min {
		font-size: 11px;
		color: var(--text-light);
		margin-top: 2px;
	}

	.payment-method {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px;
		background: var(--warm-gray);
		border-radius: 12px;
	}

	.payment-icon {
		width: 48px;
		height: 48px;
		background: white;
		border-radius: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 24px;
	}

	.payment-name {
		font-weight: 700;
		color: var(--text-dark);
		margin-bottom: 2px;
	}

	.payment-desc {
		font-size: 12px;
		color: var(--text-light);
	}

	.payment-badge {
		margin-left: auto;
		background: var(--orange);
		color: white;
		padding: 4px 12px;
		border-radius: 12px;
		font-size: 11px;
		font-weight: 700;
	}

	.summary-card {
		background: var(--cream);
		border: 2px solid var(--orange);
	}

	.summary-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 8px 0;
		color: var(--text-dark);
		font-weight: 600;
	}

	.summary-row.discount {
		color: var(--orange);
	}

	.summary-divider {
		height: 2px;
		background: var(--orange);
		margin: 8px 0;
	}

	.summary-row.total {
		font-size: 18px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.summary-row.total span:last-child {
		color: var(--orange);
	}

	.checkout-footer {
		position: fixed;
		bottom: 68px;
		left: 0;
		right: 0;
		background: white;
		padding: 16px;
		border-top: 2px solid #f0f0f0;
		box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
	}

	.footer-total {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 12px;
	}

	.footer-total span:first-child {
		color: var(--text-mid);
		font-weight: 600;
	}

	.footer-amount {
		font-size: 20px;
		font-weight: 800;
		color: var(--orange);
	}

	.order-btn {
		width: 100%;
		padding: 14px;
		background: var(--orange);
		color: white;
		border: none;
		border-radius: 12px;
		font-size: 16px;
		font-weight: 800;
		cursor: pointer;
		font-family: var(--font);
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
	}

	.order-btn:active:not(:disabled) {
		transform: scale(0.98);
	}

	.order-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.spinner {
		width: 16px;
		height: 16px;
		border: 2px solid white;
		border-top-color: transparent;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	@media (min-width: 768px) {
		.checkout-footer {
			max-width: 480px;
			left: 50%;
			transform: translateX(-50%);
		}
	}

	@media (min-width: 1024px) {
		.checkout-footer {
			max-width: 600px;
		}
	}
</style>
