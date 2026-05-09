<script>
	import { onDestroy, onMount } from 'svelte';
	import { api } from '$lib/api/client.js';
	import { authStore } from '$lib/stores/auth.js';
	import { getMemberVouchers } from '$lib/api/vouchers.js';
	import UiIcon from '$lib/components/UiIcon.svelte';

	let stats = null;
	let loading = true;
	let voucherCount = 0;
	let countdown = { hours: 5, minutes: 59, seconds: 47 };
	let countdownTimer;

	$: displayName = $authStore.user?.name?.trim() || 'GUEST';

	onMount(async () => {
		try {
			const [statsResponse, voucherResponse] = await Promise.all([
				api.getDashboardStats(),
				getMemberVouchers().catch(() => null)
			]);

			stats = statsResponse;
			voucherCount = voucherResponse?.data?.total || voucherResponse?.data?.records?.length || voucherResponse?.data?.coupons?.length || 0;
		} catch (error) {
			console.error('Failed to load dashboard:', error);
		} finally {
			loading = false;
		}

		// Countdown timer
		countdownTimer = setInterval(() => {
			countdown.seconds--;
			if (countdown.seconds < 0) {
				countdown.seconds = 59;
				countdown.minutes--;
			}
			if (countdown.minutes < 0) {
				countdown.minutes = 59;
				countdown.hours--;
			}
			if (countdown.hours < 0) {
				countdown = { hours: 5, minutes: 59, seconds: 59 };
			}
		}, 1000);
	});

	onDestroy(() => {
		if (countdownTimer) clearInterval(countdownTimer);
	});

	function formatPrice(price) {
		return new Intl.NumberFormat('id-ID', {
			style: 'currency',
			currency: 'IDR',
			minimumFractionDigits: 0
		}).format(price);
	}

	$: countdownText = `${String(countdown.hours).padStart(2, '0')} : ${String(countdown.minutes).padStart(2, '0')} : ${String(countdown.seconds).padStart(2, '0')}`;
</script>

<svelte:head>
	<title>Beranda - Coffee Order</title>
</svelte:head>

<div class="beranda-screen">
	<div class="beranda-hero">
		<div class="hero-top">
			<div class="hero-logo-wrap">
				<div class="hero-mark"><UiIcon name="coffee" size={14} /></div>
				<div class="hero-logo">KOPI NUSANTARA</div>
			</div>
			<button class="qr-btn">
				<UiIcon name="qr" size={14} />
				QR
			</button>
		</div>

		<div class="hero-banner">
			<div>
				<div class="banner-tag">HEALING<br>SERIES</div>
				<div class="banner-subtitle">Healing Journey, Flavour Echoes</div>
			</div>
			<div class="banner-prices">
				<div class="banner-badge">
					<div class="price-name">Matcha Latte</div>
					<div class="price">Rp 29k</div>
				</div>
				<div class="banner-badge">
					<div class="price-name">Oat Latte</div>
					<div class="price">Rp 33k</div>
				</div>
			</div>
		</div>

		<div class="banner-dots">
			<div class="dot active"></div>
			<div class="dot"></div>
			<div class="dot"></div>
			<div class="dot"></div>
		</div>
	</div>

	<div class="card-sheet">
		<div class="greet-row">
			<div class="greet-name">Hi, {displayName}</div>
			<div class="greet-stats">
				<div class="stat-pill">
					<div class="stat-icon icon-voucher"><UiIcon name="voucher" size={16} /></div>
					<span>{voucherCount}</span>
				</div>
				<div class="stat-pill">
					<div class="stat-icon icon-points"><UiIcon name="points" size={16} /></div>
					<span>0 pts</span>
				</div>
			</div>
		</div>

		<div class="big-btns">
			<a href="/menu" class="big-btn">
				<div class="big-btn-icon"><UiIcon name="coffee" size={22} /></div>
				<div class="big-btn-label">Order</div>
			</a>
			<button class="big-btn">
				<div class="big-btn-icon"><UiIcon name="voucher" size={22} /></div>
				<div class="big-btn-label">Voucher+</div>
			</button>
		</div>

		<div class="feature-grid">
			<div class="feature-item">
				<div class="feature-icon">
					<UiIcon name="calendar" size={18} />
					<span class="feature-badge">10%</span>
				</div>
				<div class="feature-label">Monthly Promo</div>
			</div>
			<div class="feature-item">
				<div class="feature-icon"><UiIcon name="group" size={18} /></div>
				<div class="feature-label">Referral</div>
			</div>
			<div class="feature-item">
				<div class="feature-icon"><UiIcon name="activity" size={18} /></div>
				<div class="feature-label">Spin & Win</div>
			</div>
			<div class="feature-item">
				<div class="feature-icon"><UiIcon name="check" size={18} /></div>
				<div class="feature-label">Mission</div>
			</div>
		</div>

		{#if stats}
			<div class="promo-timer-card">
				<div>
					<div class="timer-label"><UiIcon name="clock" size={14} /> Berakhir dalam</div>
					<div class="promo-desc">New User Promo!</div>
				</div>
				<div class="timer-value">{countdownText}</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.beranda-screen {
		display: flex;
		flex-direction: column;
		min-height: calc(100vh - 68px);
	}

	.beranda-hero {
		background: linear-gradient(135deg, #F4936A 0%, #E8763A 100%);
		min-height: 250px;
		padding: 18px 16px 58px;
		position: relative;
		overflow: hidden;
	}

	.hero-top {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 16px;
	}

	.hero-logo {
		font-size: 12px;
		font-weight: 800;
		color: rgba(255, 255, 255, 0.9);
		letter-spacing: 1px;
	}

	.hero-logo-wrap {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.hero-mark {
		width: 26px;
		height: 26px;
		border-radius: 8px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.18);
		color: white;
	}

	.qr-btn {
		background: white;
		border-radius: 20px;
		padding: 6px 14px;
		display: flex;
		align-items: center;
		gap: 6px;
		font-size: 12px;
		font-weight: 700;
		color: var(--orange);
		cursor: pointer;
		border: none;
		font-family: var(--font);
		transition: transform 0.2s;
	}

	.qr-btn:active {
		transform: scale(0.95);
	}

	.hero-banner {
		background: rgba(255, 255, 255, 0.15);
		border-radius: 18px;
		padding: 12px 14px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		backdrop-filter: blur(4px);
	}

	.banner-tag {
		font-size: 22px;
		font-weight: 800;
		color: white;
		line-height: 1.1;
	}

	.banner-subtitle {
		font-size: 11px;
		color: rgba(255, 255, 255, 0.8);
		margin-top: 4px;
		font-weight: 600;
	}

	.banner-prices {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.banner-badge {
		background: rgba(255, 255, 255, 0.25);
		border-radius: 10px;
		padding: 8px 12px;
		text-align: right;
	}

	.price-name {
		font-size: 10px;
		color: rgba(255, 255, 255, 0.9);
		font-weight: 600;
	}

	.price {
		font-size: 16px;
		color: white;
		font-weight: 800;
	}

	.banner-dots {
		display: flex;
		gap: 5px;
		justify-content: center;
		margin-top: 10px;
	}

	.dot {
		width: 6px;
		height: 6px;
		border-radius: 3px;
		background: rgba(255, 255, 255, 0.4);
		transition: all 0.3s;
	}

	.dot.active {
		width: 18px;
		background: white;
	}

	.card-sheet {
		background: white;
		border-radius: 24px 24px 0 0;
		margin-top: -16px;
		padding: 20px 16px 24px;
		flex: 1;
	}

	.greet-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 18px;
		gap: 12px;
	}

	.greet-name {
		font-size: 20px;
		font-weight: 800;
		color: var(--text-dark);
		line-height: 1.1;
		max-width: 52%;
	}

	.greet-stats {
		display: flex;
		gap: 12px;
		align-items: center;
	}

	.stat-pill {
		display: flex;
		align-items: center;
		gap: 5px;
		font-size: 13px;
		font-weight: 700;
		color: var(--text-dark);
	}

	.stat-icon {
		width: 22px;
		height: 22px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 11px;
	}

	.icon-voucher {
		background: #FFE5D0;
	}

	.icon-points {
		background: #FFD700;
	}

	.big-btns {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 10px;
		margin-bottom: 20px;
	}

	.big-btn {
		background: var(--cream);
		border-radius: 14px;
		padding: 16px;
		display: flex;
		align-items: center;
		gap: 12px;
		cursor: pointer;
		border: none;
		font-family: var(--font);
		transition: transform 0.15s;
		text-decoration: none;
	}

	.big-btn:active {
		transform: scale(0.97);
	}

	.big-btn-icon {
		width: 40px;
		height: 40px;
		background: var(--orange-light);
		border-radius: 10px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 20px;
		flex-shrink: 0;
	}

	.big-btn-label {
		font-size: 15px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.feature-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 8px;
		margin-bottom: 20px;
	}

	.feature-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 6px;
		cursor: pointer;
	}

	.feature-icon {
		width: 52px;
		height: 52px;
		background: var(--warm-gray);
		border-radius: 14px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 22px;
		position: relative;
		transition: transform 0.2s;
	}

	.feature-item:active .feature-icon {
		transform: scale(0.95);
	}

	.feature-badge {
		position: absolute;
		top: -4px;
		right: -4px;
		background: var(--orange);
		color: white;
		font-size: 9px;
		font-weight: 800;
		border-radius: 8px;
		padding: 2px 5px;
	}

	.feature-label {
		font-size: 10px;
		font-weight: 700;
		color: var(--text-mid);
		text-align: center;
		line-height: 1.2;
	}

	.promo-timer-card {
		background: var(--orange-light);
		border-radius: 14px;
		padding: 14px 16px;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.timer-label {
		font-size: 11px;
		color: var(--text-mid);
		font-weight: 600;
		display: inline-flex;
		align-items: center;
		gap: 6px;
	}

	.promo-desc {
		font-size: 14px;
		font-weight: 700;
		color: var(--text-dark);
		margin-top: 4px;
	}

	.timer-value {
		background: var(--orange);
		color: white;
		font-size: 13px;
		font-weight: 800;
		padding: 5px 12px;
		border-radius: 20px;
		letter-spacing: 1px;
	}

	@media (min-width: 768px) {
		.feature-grid {
			grid-template-columns: repeat(4, 1fr);
		}
	}
</style>
