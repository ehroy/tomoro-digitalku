<script>
	import '../app.css';
	import LoginModal from '$lib/components/LoginModal.svelte';
	import { authStore } from '$lib/stores/auth.js';
	import { page } from '$app/stores';
	
	let showLoginModal = false;
	let currentPath = '';
	
	$: currentPath = $page.url.pathname;
	$: isAuthenticated = $authStore.isAuthenticated;
	
	// Public routes that don't require auth
	const publicRoutes = ['/'];
	
	// Check if current route requires auth
	$: requiresAuth = !publicRoutes.includes(currentPath);
	
	// Show login modal if not authenticated and route requires auth
	$: if (requiresAuth && !isAuthenticated && typeof window !== 'undefined') {
		showLoginModal = true;
	}
</script>

<svelte:head>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
</svelte:head>

<LoginModal 
	show={showLoginModal} 
	onClose={() => {
		if (isAuthenticated) {
			showLoginModal = false;
		}
	}} 
/>

<div class="desktop-only-screen">
	<div class="desktop-only-card">
		<div class="desktop-only-badge">Mobile Only</div>
		<h1>Buka di ponsel</h1>
		<p>Aplikasi ini dioptimalkan untuk tampilan mobile agar proses order dan pembayaran lebih nyaman.</p>
		<div class="desktop-only-steps">
			<div>1. Buka browser di ponsel</div>
			<div>2. Login akun Tomoro</div>
			<div>3. Lanjutkan order dan bayar</div>
		</div>
	</div>
</div>

<div class="app-container">
	<main class="main-content">
		<slot />
	</main>
	
	<nav class="bottom-nav">
		<a href="/" class="nav-item" class:active={currentPath === '/'}>
			<div class="nav-icon">
				<svg viewBox="0 0 24 24" fill="none" stroke-width="2.5">
					<path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
					<polyline points="9 22 9 12 15 12 15 22"/>
				</svg>
			</div>
			<span>Beranda</span>
		</a>
		
		<a href="/menu" class="nav-item" class:active={currentPath === '/menu'}>
			<div class="nav-icon">
				<svg viewBox="0 0 24 24" fill="none" stroke-width="2.5">
					<path d="M17 8h1a4 4 0 0 1 0 8h-1"/>
					<path d="M3 8h14v9a4 4 0 0 1-4 4H7a4 4 0 0 1-4-4Z"/>
					<line x1="6" y1="2" x2="6" y2="4"/>
					<line x1="10" y1="2" x2="10" y2="4"/>
					<line x1="14" y1="2" x2="14" y2="4"/>
				</svg>
			</div>
			<span>Menu</span>
		</a>
		
		<a href="/history" class="nav-item" class:active={currentPath === '/history'}>
			<div class="nav-icon">
				<svg viewBox="0 0 24 24" fill="none" stroke-width="2.5">
					<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
					<polyline points="14 2 14 8 20 8"/>
					<line x1="16" y1="13" x2="8" y2="13"/>
					<line x1="16" y1="17" x2="8" y2="17"/>
					<polyline points="10 9 9 9 8 9"/>
				</svg>
			</div>
			<span>Pesanan</span>
		</a>
		
		<a href="/profile" class="nav-item" class:active={currentPath === '/profile'}>
			<div class="nav-icon">
				<svg viewBox="0 0 24 24" fill="none" stroke-width="2.5">
					<path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
					<circle cx="12" cy="7" r="4"/>
				</svg>
			</div>
			<span>Saya</span>
		</a>
	</nav>
</div>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: var(--font);
		background: var(--warm-gray);
		overflow-x: hidden;
	}

	.app-container {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
		max-width: 100vw;
		overflow-x: hidden;
	}

	.desktop-only-screen {
		display: none;
		min-height: 100vh;
		align-items: center;
		justify-content: center;
		padding: 24px;
		background: linear-gradient(180deg, #fff7f0 0%, #f6f1ec 100%);
	}

	.desktop-only-card {
		max-width: 420px;
		width: 100%;
		background: white;
		border-radius: 24px;
		padding: 24px;
		box-shadow: 0 18px 50px rgba(0, 0, 0, 0.08);
		text-align: left;
	}

	.desktop-only-badge {
		display: inline-block;
		padding: 8px 12px;
		border-radius: 999px;
		background: #fff1e8;
		color: var(--orange-dark);
		font-weight: 800;
		font-size: 12px;
		margin-bottom: 12px;
	}

	.desktop-only-card h1 {
		margin: 0 0 8px;
		font-size: 28px;
		color: #1a1a1a;
	}

	.desktop-only-card p {
		margin: 0 0 16px;
		color: #666;
		line-height: 1.5;
	}

	.desktop-only-steps {
		display: grid;
		gap: 10px;
		color: #1a1a1a;
		font-weight: 600;
	}

	.main-content {
		flex: 1;
		overflow-y: auto;
		overflow-x: hidden;
		padding-bottom: calc(68px + env(safe-area-inset-bottom));
		scrollbar-width: none;
	}
	
	.main-content::-webkit-scrollbar {
		display: none;
	}

	.bottom-nav {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		height: calc(68px + env(safe-area-inset-bottom));
		padding-bottom: env(safe-area-inset-bottom);
		background: var(--white);
		border-top: 1px solid #eee;
		display: flex;
		align-items: center;
		z-index: 100;
		max-width: 100vw;
	}

	.nav-item {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 3px;
		padding: 8px 4px calc(8px + env(safe-area-inset-bottom));
		cursor: pointer;
		transition: all 0.2s;
		text-decoration: none;
		border: none;
		background: none;
	}

	.nav-icon {
		width: 28px;
		height: 28px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.nav-icon svg {
		width: 24px;
		height: 24px;
		stroke: #bbb;
		transition: stroke 0.2s;
	}

	.nav-item span {
		font-size: 10px;
		font-weight: 600;
		color: var(--text-light);
		transition: color 0.2s;
	}

	.nav-item.active .nav-icon svg {
		stroke: var(--orange);
	}

	.nav-item.active span {
		color: var(--orange);
	}

	@media (min-width: 768px) {
		.desktop-only-screen {
			display: flex;
		}

		.app-container {
			display: none;
		}
		
		.bottom-nav {
			max-width: 480px;
			left: 50%;
			transform: translateX(-50%);
		}

		.main-content {
			padding-bottom: calc(68px + env(safe-area-inset-bottom) + 8px);
		}
	}

	@media (min-width: 1024px) {
		.bottom-nav {
			max-width: 600px;
		}
	}
</style>
