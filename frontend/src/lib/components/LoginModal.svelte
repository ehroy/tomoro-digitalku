<script>
	import { login } from '$lib/api/auth.js';
	import { goto } from '$app/navigation';
	import UiIcon from '$lib/components/UiIcon.svelte';
	
	export let show = false;
	export let onClose = () => {};
	
	let phone = '';
	let pin = '';
	let loading = false;
	let error = null;
	
	async function handleLogin() {
		if (!phone || !pin) {
			error = 'Phone dan PIN harus diisi';
			return;
		}
		
		loading = true;
		error = null;
		
		try {
			await login(phone, pin);
			onClose();
			goto('/');
		} catch (e) {
			error = e.message || 'Login gagal. Periksa phone dan PIN Anda.';
		} finally {
			loading = false;
		}
	}
	
	function handleOverlayClick(e) {
		if (e.target === e.currentTarget) {
			onClose();
		}
	}
</script>

{#if show}
	<div class="modal-overlay" on:click={handleOverlayClick} role="presentation">
		<div class="modal-content">
			<div class="modal-hero">
				<div class="hero-icon"><UiIcon name="coffee" size={18} /></div>
				<div>
					<div class="hero-kicker">Coffee Order</div>
					<h2>Masuk ke akunmu</h2>
					<p>Pesan, cek status, dan akses voucher dalam satu tempat.</p>
				</div>
			</div>

			<div class="modal-header">
				<button class="close-btn" on:click={onClose}><UiIcon name="close" size={18} /></button>
			</div>
			
			<div class="modal-body">
			{#if error}
				<div class="error-message">
					<UiIcon name="alert" size={16} /> {error}
				</div>
			{/if}
				
				<form on:submit|preventDefault={handleLogin}>
					<div class="form-group">
						<label for="phone">Nomor Telepon</label>
						<input
							type="tel"
						id="phone"
						bind:value={phone}
						placeholder="Contoh: 08xxxxxxxxxx"
						disabled={loading}
						required
					/>
				</div>
					
					<div class="form-group">
						<label for="pin">PIN</label>
						<input
							type="password"
						id="pin"
						bind:value={pin}
						placeholder="Masukkan PIN akun"
						disabled={loading}
						required
					/>
				</div>
					
					<button type="submit" class="login-btn" disabled={loading}>
						{#if loading}
							<span class="spinner"></span>
							Logging in...
						{:else}
							Login
						{/if}
					</button>
				</form>

				<div class="login-note">
					<UiIcon name="check" size={14} />
					Gunakan akun yang terdaftar di Tomoro.
				</div>

				<div class="buy-account-box">
					<div>
						<div class="buy-account-title">Belum punya akun?</div>
						<div class="buy-account-desc">Beli akun dulu untuk mulai order.</div>
					</div>
					<a class="buy-account-link" href="https://digitalku-murah.com/" target="_blank" rel="noreferrer">
						Beli akun
					</a>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.6);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
		animation: fadeIn 0.3s ease;
	}
	
	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}
	
	.modal-content {
		background: white;
		border-radius: 24px;
		width: 90%;
		max-width: 400px;
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
		animation: slideUp 0.3s ease;
		overflow: hidden;
	}
	
	@keyframes slideUp {
		from {
			transform: translateY(50px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}
	
	.modal-hero {
		padding: 18px 18px 16px;
		display: flex;
		align-items: flex-start;
		gap: 12px;
		background: linear-gradient(135deg, #fff4eb 0%, #ffe1ce 100%);
	}

	.hero-icon {
		width: 42px;
		height: 42px;
		border-radius: 14px;
		background: var(--orange);
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		box-shadow: 0 12px 24px rgba(244, 123, 32, 0.25);
	}

	.hero-kicker {
		font-size: 11px;
		font-weight: 800;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		color: #a54c15;
	}

	.modal-hero h2 {
		margin: 4px 0 4px;
		font-size: 22px;
		font-weight: 900;
		color: #1a1a1a;
	}

	.modal-hero p {
		margin: 0;
		font-size: 13px;
		line-height: 1.5;
		color: #6b5b4f;
	}

	.modal-header {
		padding: 0 18px 0;
		display: flex;
		justify-content: flex-end;
		align-items: center;
		background: white;
	}
	
	.close-btn {
		background: none;
		border: none;
		color: var(--text-dark);
		font-size: 24px;
		cursor: pointer;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		transition: background 0.2s;
	}
	
	.close-btn:hover {
		background: rgba(0, 0, 0, 0.06);
	}
	
	.modal-body {
		padding: 14px 18px 18px;
	}
	
	.error-message {
		background: #ffe5e5;
		color: #dc3545;
		padding: 12px;
		border-radius: 10px;
		margin-bottom: 16px;
		font-size: 14px;
		font-weight: 600;
	}
	
	.form-group {
		margin-bottom: 16px;
	}
	
	.form-group label {
		display: block;
		margin-bottom: 8px;
		font-weight: 700;
		color: var(--text-dark);
		font-size: 14px;
	}
	
	.form-group input {
		width: 100%;
		padding: 12px;
		border: 2px solid #e0e0e0;
		border-radius: 12px;
		font-size: 16px;
		font-family: var(--font);
		transition: border-color 0.2s;
	}
	
	.form-group input:focus {
		outline: none;
		border-color: var(--orange);
	}
	
	.form-group input:disabled {
		background: #f5f5f5;
		cursor: not-allowed;
	}
	
	.login-btn {
		width: 100%;
		padding: 14px;
		background: linear-gradient(135deg, var(--orange), var(--orange-dark));
		color: white;
		border: none;
		border-radius: 16px;
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
	
	.login-btn:hover:not(:disabled) {
		background: var(--orange-dark);
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(244, 123, 32, 0.4);
	}
	
	.login-btn:disabled {
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

	.login-note {
		margin-top: 14px;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		font-size: 12px;
		font-weight: 600;
		color: #8a6d3b;
	}

	.buy-account-box {
		margin-top: 12px;
		padding: 14px;
		border-radius: 16px;
		background: #fff7f0;
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 12px;
	}

	.buy-account-title {
		font-size: 13px;
		font-weight: 800;
		color: #1a1a1a;
	}

	.buy-account-desc {
		font-size: 12px;
		color: #7a6a5d;
		margin-top: 2px;
	}

	.buy-account-link {
		flex-shrink: 0;
		padding: 10px 14px;
		border-radius: 999px;
		background: var(--orange);
		color: white;
		text-decoration: none;
		font-size: 12px;
		font-weight: 800;
	}

	.buy-account-link:hover {
		background: var(--orange-dark);
	}

	@media (max-width: 480px) {
		.modal-content {
			width: calc(100% - 24px);
		}
	}
</style>
