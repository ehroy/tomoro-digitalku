<script>
	import { authStore } from '$lib/stores/auth.js';
	import { logout } from '$lib/api/auth.js';
	import { goto } from '$app/navigation';
	import UiIcon from '$lib/components/UiIcon.svelte';
	
	$: user = $authStore.user;
	
	async function handleLogout() {
		logout();
		goto('/');
	}
</script>

	<div class="profile-page">
		<div class="saya-hero">
			<div class="saya-hero-deco"><UiIcon name="coffee" size={72} /></div>
			<div class="profile-row">
				<div class="avatar"><UiIcon name="user" size={28} /></div>
				<div class="profile-name">{user?.name || 'GUEST'}</div>
			</div>
		</div>

	<div class="profile-content">
		<div class="stats-card">
			<div class="stats-item">
				<div class="stats-label">Phone</div>
				<div class="stats-value">{user?.phone || '-'}</div>
			</div>
			<div class="stats-item">
				<div class="stats-label">Member Code</div>
				<div class="stats-value">{user?.memberCode || '-'}</div>
			</div>
			<div class="stats-item">
				<div class="stats-label">Points</div>
				<div class="stats-icon"><UiIcon name="points" size={22} /></div>
				<div><span class="stats-value">0</span> <span class="stats-unit">pts</span></div>
			</div>
		</div>

		<div class="menu-list">
			<div class="menu-list-item">
				<div class="menu-list-icon"><UiIcon name="location" size={18} /></div>
				<div class="menu-list-label">Alamat Saya</div>
				<div class="menu-list-arrow">›</div>
			</div>
			<div class="menu-list-item">
				<div class="menu-list-icon"><UiIcon name="refresh" size={18} /></div>
				<div class="menu-list-label">Redeem Center</div>
				<div class="menu-list-arrow">›</div>
			</div>
			<div class="menu-list-item">
				<div class="menu-list-icon"><UiIcon name="help" size={18} /></div>
				<div class="menu-list-label">Pusat Bantuan</div>
				<div class="menu-list-arrow">›</div>
			</div>
			<div class="menu-list-item">
				<div class="menu-list-icon"><UiIcon name="globe" size={18} /></div>
				<div class="menu-list-label">Pengaturan Bahasa</div>
				<div class="menu-list-arrow">›</div>
			</div>
			<div class="menu-list-item">
				<div class="menu-list-icon"><UiIcon name="terms" size={18} /></div>
				<div class="menu-list-label">Ketentuan Layanan</div>
				<div class="menu-list-arrow">›</div>
			</div>
			<button class="menu-list-item" on:click={handleLogout}>
				<div class="menu-list-icon" style="background:#FFEAEA;"><UiIcon name="logout" size={18} /></div>
				<div class="menu-list-label" style="color:#e74c3c;">Keluar</div>
				<div class="menu-list-arrow">›</div>
			</button>
		</div>
	</div>
</div>

<style>
	.profile-page {
		min-height: 100vh;
		background: var(--warm-gray);
		padding-bottom: 80px;
	}

	.saya-hero {
		background: var(--cream);
		padding: 24px 16px 30px;
		position: relative;
		overflow: hidden;
	}

	.saya-hero-deco {
		position: absolute;
		right: 20px;
		top: 10px;
		font-size: 70px;
		opacity: 0.25;
	}

	.profile-row {
		display: flex;
		align-items: center;
		gap: 14px;
	}

	.avatar {
		width: 60px;
		height: 60px;
		border-radius: 50%;
		background: #FFD0B8;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 28px;
		flex-shrink: 0;
	}

	.profile-name {
		font-size: 20px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.profile-content {
		padding: 16px 12px;
	}

	.stats-card {
		background: white;
		border-radius: var(--card-radius);
		padding: 18px 16px;
		margin-bottom: 16px;
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
		border: 1px solid #f0f0f0;
	}

	.stats-item {
		text-align: center;
		padding: 0 8px;
		border-right: 1px solid #f0f0f0;
	}
	.stats-item:last-child { border-right: none; }

	.stats-label {
		font-size: 11px;
		color: var(--text-mid);
		font-weight: 700;
		margin-bottom: 6px;
	}

	.stats-value {
		font-size: 16px;
		font-weight: 800;
		color: var(--text-dark);
		word-break: break-all;
	}

	.stats-unit {
		font-size: 11px;
		color: var(--text-light);
		font-weight: 600;
	}

	.stats-icon {
		font-size: 22px;
		margin-bottom: 4px;
	}

	.menu-list {
		background: white;
		border-radius: 16px;
		overflow: hidden;
	}

	.menu-list-item {
		display: flex;
		align-items: center;
		gap: 14px;
		padding: 16px;
		border-bottom: 1px solid #f5f5f5;
		cursor: pointer;
		transition: background 0.15s;
		border: none;
		background: white;
		width: 100%;
		text-align: left;
		font-family: var(--font);
	}
	.menu-list-item:last-child { border-bottom: none; }
	.menu-list-item:hover { background: var(--warm-gray); }

	.menu-list-icon {
		width: 36px;
		height: 36px;
		border-radius: 10px;
		background: var(--warm-gray);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 18px;
		flex-shrink: 0;
	}

	.menu-list-label {
		flex: 1;
		font-size: 14px;
		font-weight: 700;
		color: var(--text-dark);
	}

	.menu-list-arrow {
		color: var(--text-light);
		font-size: 16px;
	}
</style>
