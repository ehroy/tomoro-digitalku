<script>
	import { onDestroy, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getOrderHistory, getPayStatus } from '$lib/api/orders.js';
	import UiIcon from '$lib/components/UiIcon.svelte';

	let checkoutUrl = '';
	let guidePicture = '';
	let tradeOrderCode = '';
	let orderCode = '';
	let orderNo = '';
	let storeCode = '';
	let latestOrder = null;
	let payStatus = null;
	let loading = true;
	let refreshing = false;
	let error = '';
	let intervalId;

	function parseOrders(response) {
		return response?.data?.records || response?.data || [];
	}

	function matchesOrder(order) {
		if (!order) return false;
		const candidateCodes = [tradeOrderCode, orderCode, orderNo, order?.paymentOrderCode, order?.transCode].filter(Boolean);
		if (candidateCodes.length === 0) return true;
		const orderFields = [order.orderCode, order.paymentOrderCode, order.transCode].filter(Boolean);
		return candidateCodes.some((code) => orderFields.includes(code));
	}

	function formatDate(dateString) {
		if (!dateString) return '-';
		return new Date(dateString).toLocaleString('id-ID', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatPrice(value) {
		const amount = Number(value || 0);
		return new Intl.NumberFormat('id-ID', {
			style: 'currency',
			currency: 'IDR',
			minimumFractionDigits: 0
		}).format(amount);
	}

	function getStatusLabel(status) {
		if (!status) return 'Menunggu pembayaran';
		const statusLower = String(status).toLowerCase();
		if (statusLower.includes('cancel') || statusLower.includes('batal')) return 'Dibatalkan';
		if (statusLower.includes('complete') || statusLower.includes('selesai') || statusLower.includes('done')) return 'Selesai';
		if (statusLower.includes('pending')) return 'Menunggu pembayaran';
		if (statusLower.includes('process')) return 'Diproses';
		return status;
	}

	function getPaymentLabel(order) {
		if (!order) return 'Menunggu pembayaran';
		if (payStatus?.paymentStatus) return payStatus.paymentStatus;
		return order.paymentStatus || (order.status === 'Dibatalkan' ? 'Dibatalkan' : 'Menunggu pembayaran');
	}

	function getOrderItems(order) {
		return order?.items || [];
	}

	function formatPaymentCode(value) {
		return value && value !== '0' ? value : '-';
	}

	function loadStoredState() {
		try {
			const raw = sessionStorage.getItem('checkoutStatusState');
			return raw ? JSON.parse(raw) : {};
		} catch {
			return {};
		}
	}

	function saveStoredState(nextState) {
		try {
			sessionStorage.setItem('checkoutStatusState', JSON.stringify(nextState));
		} catch {
			// ignore storage failures
		}
	}

	function getOrderState(order) {
		if (!order) return 'pending';
		if (payStatus?.isPaid === true || payStatus?.paymentStatus === 'Sudah dibayar' || order.paymentStatus === 'Sudah dibayar' || order.status === 'Selesai') return 'paid';
		if (order.status === 'Dibatalkan') return 'cancelled';
		return 'pending';
	}

	function getStateLabel(order) {
		const state = getOrderState(order);
		if (state === 'paid') return 'Pesanan berhasil dibayar';
		if (state === 'cancelled') return 'Pesanan dibatalkan';
		return 'Menunggu pembayaran';
	}

	async function refreshStatus() {
		refreshing = true;
		error = '';
		try {
			const requests = [];
			if (tradeOrderCode) requests.push(getPayStatus(tradeOrderCode));
			requests.push(getOrderHistory());
			const responses = await Promise.allSettled(requests);
			const payStatusResponse = tradeOrderCode ? responses[0] : null;
			const historyResponse = tradeOrderCode ? responses[1] : responses[0];
			payStatus = payStatusResponse?.status === 'fulfilled' ? (payStatusResponse.value?.data?.data || payStatusResponse.value?.data || null) : null;
			if (!checkoutUrl && payStatus?.redirectUrl) {
				checkoutUrl = payStatus.redirectUrl;
				saveStoredState({ checkoutUrl, guidePicture, tradeOrderCode, orderCode, orderNo, storeCode });
			}
			const historyData = historyResponse?.status === 'fulfilled' ? historyResponse.value : null;
			const orders = parseOrders(historyData);
			latestOrder = orders.find(matchesOrder) || orders[0] || null;
			if (payStatus?.isPaid === true || payStatus?.paymentStatus === 'Sudah dibayar' || latestOrder?.paymentStatus === 'Sudah dibayar' || latestOrder?.status === 'Selesai') {
				if (intervalId) clearInterval(intervalId);
				intervalId = null;
			}
		} catch (err) {
			error = 'Gagal memuat status pesanan';
			console.error(err);
		} finally {
			loading = false;
			refreshing = false;
		}
	}

	onMount(() => {
		const stored = loadStoredState();
		const params = new URLSearchParams(window.location.search);
		checkoutUrl = params.get('checkoutUrl') || stored.checkoutUrl || '';
		guidePicture = params.get('guidePicture') || stored.guidePicture || 'https://static.tomoro-coffee.id/prod/material20220721/Payment/wait1/t.png';
		tradeOrderCode = params.get('tradeOrderCode') || stored.tradeOrderCode || '';
		orderCode = params.get('orderCode') || stored.orderCode || '';
		orderNo = params.get('orderNo') || stored.orderNo || '';
		storeCode = params.get('storeCode') || stored.storeCode || '';

		saveStoredState({ checkoutUrl, guidePicture, tradeOrderCode, orderCode, orderNo, storeCode });

		refreshStatus();
		intervalId = setInterval(refreshStatus, 5000);
	});

	onDestroy(() => {
		if (intervalId) clearInterval(intervalId);
	});
</script>

<svelte:head>
	<title>Status Pembayaran - Coffee Order</title>
</svelte:head>

<div class="status-screen">
	<div class="status-card">
		<div class="status-hero">
			<img src={guidePicture} alt="Panduan pembayaran" />
		</div>

		<div class="status-body">
			<div class="status-title">Cek Status Order</div>
			<div class="status-subtitle">Status pembayaran dipantau langsung dari Tomoro pay status API.</div>

			{#if latestOrder}
				<div class:state-paid={getOrderState(latestOrder) === 'paid'} class:state-cancelled={getOrderState(latestOrder) === 'cancelled'} class="state-banner">
					{getStateLabel(latestOrder)}
				</div>
			{/if}

			{#if error}
				<div class="status-error">{error}</div>
			{/if}

			<div class="status-meta">
				<div><span>Order</span><strong>{latestOrder?.orderCode || tradeOrderCode || orderCode || orderNo || '-'}</strong></div>
				<div><span>Outlet</span><strong>{latestOrder?.storeName || latestOrder?.storeCode || storeCode || '-'}</strong></div>
				<div><span>Status Order</span><strong>{getStatusLabel(latestOrder?.status)}</strong></div>
				<div><span>Status Bayar</span><strong>{getPaymentLabel(latestOrder)}</strong></div>
				<div><span>Waktu</span><strong>{formatDate(latestOrder?.createdAt)}</strong></div>
			</div>

			{#if payStatus}
				<div class="detail-card">
					<div class="detail-card-title">Status Pembayaran</div>
					<div class="detail-grid">
						<div><span>Trade Order Code</span><strong>{payStatus.tradeOrderCode || tradeOrderCode || '-'}</strong></div>
						<div><span>Status</span><strong>{payStatus.paymentStatus || 'Menunggu pembayaran'}</strong></div>
						<div><span>Paid</span><strong>{payStatus.isPaid ? 'Ya' : 'Tidak'}</strong></div>
						<div><span>Pay Money</span><strong>{formatPrice(payStatus.payMoney || latestOrder?.paymentAmount)}</strong></div>
						<div><span>Channel</span><strong>{payStatus.payChannelName || '-'}</strong></div>
						<div><span>Currency</span><strong>{payStatus.currencyCode || payStatus.currencyUnit || '-'}</strong></div>
					</div>
				</div>
			{/if}

			{#if latestOrder}
				<div class="detail-card">
					<div class="detail-card-title">Detail Tomoro</div>
					<div class="detail-grid">
						<div><span>Payment Order Code</span><strong>{formatPaymentCode(latestOrder.paymentOrderCode)}</strong></div>
						<div><span>Trans Code</span><strong>{formatPaymentCode(latestOrder.transCode)}</strong></div>
						<div><span>Pay Type</span><strong>{latestOrder.payType ?? '-'}</strong></div>
						<div><span>Pay Channel</span><strong>{latestOrder.payChannelCode || '-'}</strong></div>
						<div><span>Voucher</span><strong>{latestOrder.voucherCode || '-'}</strong></div>
						<div><span>Total</span><strong>{formatPrice(latestOrder.totalAmount)}</strong></div>
						<div><span>Diskon</span><strong>{formatPrice(latestOrder.discountAmount)}</strong></div>
						<div><span>Bayar</span><strong>{formatPrice(latestOrder.paymentAmount)}</strong></div>
					</div>
				</div>
			{/if}

			{#if getOrderItems(latestOrder).length > 0}
				<div class="items-card">
					<div class="detail-card-title">Item Pesanan</div>
					<div class="items-list">
						{#each getOrderItems(latestOrder) as item}
							<div class="item-row">
								<div class="item-main">
									<div class="item-name">{item.productName}</div>
									<div class="item-code">{item.productCode}</div>
								</div>
								<div class="item-side">
									<div>{item.quantity}x</div>
									<strong>{formatPrice(item.subtotal || item.price * item.quantity)}</strong>
								</div>
							</div>
						{/each}
					</div>
				</div>
			{/if}

			<div class="status-actions">
				{#if getOrderState(latestOrder) !== 'paid'}
					{#if checkoutUrl}
						<a class="primary-btn" href={checkoutUrl} target="_blank" rel="noreferrer">
							<UiIcon name="voucher" size={16} />
							Buka Link Bayar
						</a>
					{:else}
						<button class="primary-btn" disabled>
							<UiIcon name="voucher" size={16} />
							Link tidak tersedia
						</button>
					{/if}
				{:else}
					<div class="paid-note">Pembayaran selesai</div>
				{/if}
				<button class="secondary-btn" on:click={refreshStatus} disabled={refreshing}>
					{#if refreshing}
						Mengecek...
					{:else}
						Refresh Status
					{/if}
				</button>
			</div>

			<div class="footer-link">
				<button class="link-btn" on:click={() => goto('/history')}>Lihat Riwayat Pesanan</button>
			</div>
		</div>
	</div>
</div>

<style>
	.status-screen {
		min-height: calc(100vh - 68px);
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 16px;
		background: linear-gradient(180deg, #fff7f0 0%, #f6f1ec 100%);
	}

	.status-card {
		width: 100%;
		max-width: 480px;
		background: white;
		border-radius: 24px;
		overflow: hidden;
		box-shadow: 0 18px 50px rgba(0, 0, 0, 0.08);
	}

	.status-hero img {
		width: 100%;
		display: block;
		aspect-ratio: 1 / 1;
		object-fit: cover;
		background: #f4f4f4;
	}

	.status-body {
		padding: 18px;
	}

	.status-title {
		font-size: 22px;
		font-weight: 800;
		color: #1a1a1a;
	}

	.status-subtitle {
		margin-top: 6px;
		color: #666;
		font-size: 13px;
		line-height: 1.5;
	}

	.state-banner {
		margin-top: 14px;
		padding: 12px 14px;
		border-radius: 14px;
		font-size: 13px;
		font-weight: 800;
		background: #fff6eb;
		color: #b45309;
	}

	.state-banner.state-paid {
		background: #ecfdf3;
		color: #027a48;
	}

	.state-banner.state-cancelled {
		background: #fff1f3;
		color: #c01048;
	}

	.status-error {
		margin-top: 12px;
		padding: 10px 12px;
		border-radius: 12px;
		background: #fff0f0;
		color: #b42318;
		font-size: 13px;
	}

	.status-meta {
		margin-top: 16px;
		display: grid;
		gap: 12px;
	}

	.status-meta div {
		display: flex;
		justify-content: space-between;
		gap: 12px;
		font-size: 13px;
	}

	.status-meta span {
		color: #888;
	}

	.status-meta strong {
		color: #1a1a1a;
		text-align: right;
	}

	.status-actions {
		margin-top: 18px;
		display: grid;
		gap: 10px;
	}

	.detail-card,
	.items-card {
		margin-top: 16px;
		padding: 14px;
		border-radius: 18px;
		background: #faf7f4;
		border: 1px solid #eee3da;
	}

	.detail-card-title {
		font-size: 14px;
		font-weight: 800;
		color: #1a1a1a;
		margin-bottom: 12px;
	}

	.detail-grid {
		display: grid;
		gap: 10px;
	}

	.detail-grid div,
	.item-row {
		display: flex;
		justify-content: space-between;
		gap: 12px;
		font-size: 13px;
	}

	.detail-grid span,
	.item-code {
		color: #888;
	}

	.detail-grid strong,
	.item-side strong {
		color: #1a1a1a;
		text-align: right;
	}

	.items-list {
		display: grid;
		gap: 12px;
	}

	.item-main {
		min-width: 0;
	}

	.item-name {
		font-weight: 700;
		color: #1a1a1a;
	}

	.item-side {
		text-align: right;
		min-width: 72px;
	}

	.primary-btn,
	.secondary-btn,
	.link-btn {
		border: none;
		border-radius: 16px;
		font-weight: 700;
		cursor: pointer;
		font-family: inherit;
	}

	.primary-btn,
	.secondary-btn {
		padding: 14px 16px;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		text-decoration: none;
	}

	.primary-btn {
		background: #f47b20;
		color: white;
	}

	.secondary-btn {
		background: #f3f3f3;
		color: #1a1a1a;
	}

	.link-btn {
		background: transparent;
		color: #f47b20;
		padding: 4px 0 0;
		font-size: 13px;
	}

	.primary-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.paid-note {
		padding: 14px 16px;
		border-radius: 16px;
		background: #ecfdf3;
		color: #027a48;
		font-weight: 800;
		text-align: center;
	}
</style>
