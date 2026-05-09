<script>
	import { onDestroy, onMount } from 'svelte';
	import { getOrderHistory } from '$lib/api/orders.js';
	import { getMemberVouchers } from '$lib/api/vouchers.js';
	import UiIcon from '$lib/components/UiIcon.svelte';

	let orders = [];
	let vouchers = [];
	let loading = true;
	let loadingVouchers = true;
	let activeTab = 'outlet';
	let refreshTimer;
	let refreshing = false;

	function normalizeOrders(response) {
		return response?.data?.records || response?.data || [];
	}

	function getPaymentStatus(order) {
		if (order?.paymentStatus) return order.paymentStatus;
		if (order?.statusCode === 1) return 'Menunggu pembayaran';
		if (order?.statusCode === 6) return 'Dibatalkan';
		return 'Sudah dibayar';
	}

	function getOrderStatus(order) {
		if (order?.status) return order.status;
		if (order?.statusCode === 6) return 'Dibatalkan';
		if (order?.statusCode === 1) return 'Menunggu pembayaran';
		return 'Menunggu pembayaran';
	}

	async function loadHistory() {
		refreshing = true;
		try {
			const [ordersResponse, vouchersResponse] = await Promise.all([
				getOrderHistory(),
				getMemberVouchers()
			]);

			orders = normalizeOrders(ordersResponse);
			vouchers = vouchersResponse.data?.records || vouchersResponse.data?.coupons || vouchersResponse.data || [];

			if (orders.length > 0 && orders[0].createdAt) {
				orders.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
			}
		} catch (error) {
			console.error('Failed to load history data:', error);
			orders = [];
			vouchers = [];
		} finally {
			loading = false;
			loadingVouchers = false;
			refreshing = false;
		}
	}

	onMount(async () => {
		await loadHistory();
		refreshTimer = setInterval(loadHistory, 5000);
	});

	onDestroy(() => {
		if (refreshTimer) clearInterval(refreshTimer);
	});

	function formatPrice(price) {
		return new Intl.NumberFormat('id-ID', {
			style: 'currency',
			currency: 'IDR',
			minimumFractionDigits: 0
		}).format(price);
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

	function getStatusClass(status) {
		if (!status) return 'status-process';
		
		const statusLower = status.toLowerCase();
		if (statusLower.includes('cancel') || statusLower.includes('batal')) return 'status-cancelled';
		if (statusLower.includes('complete') || statusLower.includes('selesai') || statusLower.includes('done')) return 'status-done';
		return 'status-process';
	}

	function getStatusLabel(status) {
		if (!status) return 'Diproses';
		
		const statusLower = status.toLowerCase();
		if (statusLower.includes('cancel') || statusLower.includes('batal')) return 'Dibatalkan';
		if (statusLower.includes('complete') || statusLower.includes('selesai') || statusLower.includes('done')) return 'Selesai';
		if (statusLower.includes('pending')) return 'Pending';
		if (statusLower.includes('process')) return 'Diproses';
		return status;
	}

	function getItemIcon(name) {
		if (!name) return 'coffee';

		const nameLower = name.toLowerCase();
		if (nameLower.includes('croissant') || nameLower.includes('pastry') || nameLower.includes('bundle')) return 'package';
		return 'coffee';
	}

	function formatVoucherRange(voucher) {
		if (!voucher?.couponActiveStartTime && !voucher?.couponActiveEndTime) return '-';
		const start = voucher.couponActiveStartTime || '-';
		const end = voucher.couponActiveEndTime || '-';
		return `${start} - ${end}`;
	}
</script>

<svelte:head>
	<title>Pesanan - Coffee Order</title>
</svelte:head>

<div class="pesanan-screen">
		<div class="pesanan-top">
			<div class="tommunity-title">TOMMUNITY</div>
			<div class="refresh-row">
				<div class="refresh-hint">Realtime order history</div>
				<button class="refresh-btn" on:click={loadHistory} disabled={refreshing}>{refreshing ? 'Memuat...' : 'Refresh'}</button>
			</div>
			<div class="community-grid">
				<div class="community-card">
					<div class="community-card-left">
						<div class="community-card-icon"><UiIcon name="calendar" size={18} /></div>
						<div class="community-card-label">Survey</div>
					</div>
					<div class="community-card-arrow">›</div>
				</div>
				<div class="community-card">
					<div class="community-card-left">
						<div class="community-card-icon"><UiIcon name="activity" size={18} /></div>
						<div class="community-card-label">Activity</div>
					</div>
					<div class="community-card-arrow">›</div>
			</div>
		</div>
	</div>

	<div class="pesanan-content">
		<div class="pesanan-section-header">
			<div class="pesanan-section-title">Pesanan Saya</div>
			<div class="tabs">
				<div 
					class="tab" 
					class:active={activeTab === 'outlet'}
					on:click={() => activeTab = 'outlet'}
					on:keypress={() => activeTab = 'outlet'}
					role="button"
					tabindex="0"
				>
					Outlet
				</div>
				<div 
					class="tab" 
					class:active={activeTab === 'voucher'}
					on:click={() => activeTab = 'voucher'}
					on:keypress={() => activeTab = 'voucher'}
					role="button"
					tabindex="0"
				>
					Paket Voucher
				</div>
				<div 
					class="tab" 
					class:active={activeTab === 'benefit'}
					on:click={() => activeTab = 'benefit'}
					on:keypress={() => activeTab = 'benefit'}
					role="button"
					tabindex="0"
				>
					Benefit Card
				</div>
			</div>
		</div>

		{#if activeTab === 'outlet'}
			<div class="tab-content">
				{#if loading}
					<div class="loading-state">
						<div class="spinner"></div>
						<p>Loading pesanan...</p>
					</div>
				{:else if orders.length === 0}
					<div class="empty-state">
						<div class="empty-icon"><UiIcon name="package" size={48} /></div>
						<div class="empty-text">Belum ada pesanan</div>
						<a href="/menu" class="browse-btn">Mulai Pesan</a>
					</div>
				{:else}
					{#each orders as order}
						<div class="order-card">
							<div class="order-card-header">
								<div>
									<div class="order-outlet">{order.storeName || order.storeCode || 'Outlet'}</div>
									<div class="order-date">{formatDate(order.createdAt || order.created_at)}</div>
								</div>
					<div class="order-status {getStatusClass(getOrderStatus(order))}">
						{getStatusLabel(getOrderStatus(order))}
					</div>
			</div>
			<div class="payment-state">
				<span>Payment</span>
				<strong>{getPaymentStatus(order)}</strong>
			</div>

					{#if order.items && order.items.length > 0}
						{#each order.items as item}
							<div class="order-item-row">
								<div class="order-item-img"><UiIcon name={getItemIcon(item.name || item.productName)} size={20} /></div>
								<div class="order-item-info">
									<div class="order-item-name">{item.name || item.productName}</div>
								</div>
										<div style="text-align: right;">
											<div class="order-item-price">{formatPrice(item.subtotal || (item.price * item.quantity))}</div>
											<div class="order-item-qty">{item.quantity} Item{item.quantity > 1 ? 's' : ''}</div>
										</div>
									</div>
								{/each}
							{/if}

						<div class="order-total-row">
							<span>Total</span>
							<span class="order-total">{formatPrice(order.totalAmount || order.total_amount || 0)}</span>
						</div>

							<button class="reorder-btn">Pesan Lagi</button>
						</div>
					{/each}
				{/if}
			</div>
		{:else if activeTab === 'voucher'}
			<div class="tab-content">
				{#if loadingVouchers}
					<div class="loading-state">
						<div class="spinner"></div>
						<p>Loading voucher...</p>
					</div>
				{:else if vouchers.length === 0}
					<div class="empty-state">
						<div class="empty-icon"><UiIcon name="voucher" size={48} /></div>
						<div class="empty-text">Belum ada voucher</div>
					</div>
				{:else}
					{#each vouchers as voucher}
						<div class="order-card voucher-card">
							<div class="order-card-header">
								<div>
									<div class="order-outlet">{voucher.couponName || 'Voucher'}</div>
									<div class="order-date">{voucher.couponCode || voucher.accountCouponCode || '-'}</div>
								</div>
								<div class="order-status status-process">Aktif</div>
							</div>

							<div class="voucher-desc">{voucher.couponDesc || '-'}</div>
							<div class="voucher-meta">
								<div><span>Periode</span><strong>{formatVoucherRange(voucher)}</strong></div>
								<div><span>Status</span><strong>{voucher.isValidity ? 'Valid' : 'Tidak valid'}</strong></div>
							</div>
						</div>
					{/each}
				{/if}
			</div>
		{:else}
			<div class="tab-content">
				<div class="empty-state">
					<div class="empty-icon"><UiIcon name="voucher" size={48} /></div>
					<div class="empty-text">Belum ada benefit card</div>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.pesanan-screen {
		display: flex;
		flex-direction: column;
		height: calc(100vh - 68px);
	}

	.pesanan-top {
		background: var(--cream);
		padding: 20px 16px 16px;
		flex-shrink: 0;
	}

	.refresh-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 12px;
		margin: 8px 0 14px;
	}

	.refresh-hint {
		font-size: 12px;
		color: var(--text-light);
		font-weight: 700;
	}

	.refresh-btn {
		border: none;
		background: var(--orange);
		color: white;
		border-radius: 999px;
		padding: 8px 12px;
		font-size: 12px;
		font-weight: 700;
	}

	.tommunity-title {
		font-size: 18px;
		font-weight: 800;
		color: var(--text-dark);
		margin-bottom: 12px;
		letter-spacing: 0.5px;
	}

	.community-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 10px;
	}

	.community-card {
		background: #FFF6EE;
		border-radius: 14px;
		padding: 16px;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: space-between;
		transition: transform 0.2s;
	}

	.community-card:active {
		transform: scale(0.97);
	}

	.community-card-left {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.community-card-icon {
		font-size: 28px;
		margin-bottom: 4px;
	}

	.community-card-label {
		font-size: 14px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.community-card-arrow {
		color: var(--text-light);
		font-size: 20px;
	}

	.pesanan-content {
		background: var(--warm-gray);
		flex: 1;
		overflow-y: auto;
		scrollbar-width: none;
	}

	.pesanan-content::-webkit-scrollbar {
		display: none;
	}

	.pesanan-section-header {
		background: white;
		padding: 16px;
		margin-bottom: 8px;
	}

	.pesanan-section-title {
		font-size: 16px;
		font-weight: 800;
		color: var(--text-dark);
		margin-bottom: 12px;
	}

	.tabs {
		display: flex;
		gap: 20px;
		border-bottom: 2px solid #f0f0f0;
		padding-bottom: 0;
	}

	.tab {
		font-size: 13px;
		font-weight: 700;
		color: var(--text-light);
		padding-bottom: 10px;
		cursor: pointer;
		border-bottom: 2px solid transparent;
		margin-bottom: -2px;
		transition: all 0.2s;
	}

	.tab.active {
		color: var(--orange);
		border-bottom-color: var(--orange);
	}

	.tab-content {
		padding-top: 12px;
	}

	.loading-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 40px;
		color: var(--text-light);
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 3px solid #f3f3f3;
		border-top: 3px solid var(--orange);
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 12px;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	.empty-state {
		text-align: center;
		padding: 60px 20px;
		color: var(--text-light);
	}

	.empty-icon {
		font-size: 48px;
		margin-bottom: 12px;
		opacity: 0.5;
	}

	.empty-text {
		font-size: 14px;
		font-weight: 600;
		margin-bottom: 16px;
	}

	.browse-btn {
		display: inline-block;
		padding: 12px 24px;
		background: var(--orange);
		color: white;
		text-decoration: none;
		border-radius: 12px;
		font-weight: 700;
		transition: all 0.2s;
	}

	.browse-btn:active {
		transform: scale(0.95);
	}

	.order-card {
		background: white;
		border-radius: 16px;
		margin: 0 12px 12px;
		padding: 16px;
	}

	.voucher-card {
		border-left: 4px solid var(--orange);
	}

	.order-card-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 12px;
	}

	.payment-state {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 10px 12px;
		margin-bottom: 12px;
		border-radius: 12px;
		background: #fff7f0;
		font-size: 12px;
	}

	.payment-state span {
		color: var(--text-light);
		font-weight: 700;
	}

	.payment-state strong {
		color: var(--text-dark);
	}

	.order-outlet {
		font-size: 14px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.order-date {
		font-size: 11px;
		color: var(--text-light);
		font-weight: 600;
		margin-top: 2px;
	}

	.order-status {
		font-size: 12px;
		font-weight: 700;
		padding: 4px 10px;
		border-radius: 10px;
	}

	.status-cancelled {
		color: #999;
		background: #f5f5f5;
	}

	.status-done {
		color: #27ae60;
		background: #eafaf1;
	}

	.status-process {
		color: var(--orange);
		background: var(--orange-light);
	}

	.order-item-row {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 10px 0;
		border-top: 1px solid #f5f5f5;
	}

	.order-item-img {
		width: 50px;
		height: 50px;
		border-radius: 10px;
		background: var(--warm-gray);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 24px;
		flex-shrink: 0;
	}

	.order-item-info {
		flex: 1;
	}

	.order-item-name {
		font-size: 13px;
		font-weight: 700;
		color: var(--text-dark);
	}

	.order-item-price {
		font-size: 14px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.order-item-qty {
		font-size: 11px;
		color: var(--text-light);
		font-weight: 600;
	}

	.order-total-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 12px 0;
		margin-top: 8px;
		border-top: 2px solid #f0f0f0;
		font-weight: 700;
		color: var(--text-dark);
	}

	.order-total {
		font-size: 18px;
		font-weight: 800;
		color: var(--orange);
	}

	.voucher-desc {
		font-size: 13px;
		color: var(--text-mid);
		line-height: 1.5;
		margin: 12px 0;
	}

	.voucher-meta {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 10px;
		font-size: 12px;
		color: var(--text-light);
	}

	.voucher-meta div {
		background: var(--warm-gray);
		border-radius: 12px;
		padding: 10px 12px;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.voucher-meta strong {
		font-size: 12px;
		color: var(--text-dark);
	}

	.reorder-btn {
		width: 100%;
		padding: 11px;
		border: 1.5px solid #ddd;
		border-radius: 20px;
		font-size: 13px;
		font-weight: 700;
		color: var(--text-dark);
		background: white;
		cursor: pointer;
		margin-top: 12px;
		font-family: var(--font);
		transition: all 0.2s;
	}

	.reorder-btn:hover {
		border-color: var(--orange);
		color: var(--orange);
	}

	.reorder-btn:active {
		transform: scale(0.98);
	}
</style>
