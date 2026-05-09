<script>
	import { onDestroy, onMount } from 'svelte';
	import { getStores } from '$lib/api/stores.js';
	import { getUserLocation } from '$lib/stores/location.js';
	import { selectedStoreStore } from '$lib/stores/voucher.js';
	import { goto } from '$app/navigation';
	import 'leaflet/dist/leaflet.css';
	import UiIcon from '$lib/components/UiIcon.svelte';

	let stores = [];
	let loading = true;
	let mapLoading = true;
	let searchQuery = '';
	let userLocation = null;
	let activeFilter = 'all';
	let mapContainer;
	let mapInstance;
	let mapLayer;
	let leaflet;

	$: filteredStores = searchQuery 
		? stores.filter(store => 
			store.storeName.toLowerCase().includes(searchQuery.toLowerCase()) ||
			store.storeAddress.toLowerCase().includes(searchQuery.toLowerCase())
		)
		: stores;

	$: visibleStores = filteredStores.filter((store) => {
		if (activeFilter === 'open') return store.businessStatus === 0;
		if (activeFilter === 'closed') return store.businessStatus !== 0;
		return true;
	});

	$: openCount = stores.filter((store) => store.businessStatus === 0).length;
	$: closedCount = stores.length - openCount;
	$: currentLabel = userLocation
		? `${userLocation.latitude.toFixed(3)}, ${userLocation.longitude.toFixed(3)}`
		: 'Lokasi default';
	$: if (mapInstance && leaflet) {
		updateMap();
	}

	onMount(async () => {
		try {
			// Get user location
			userLocation = await getUserLocation();
			leaflet = await import('leaflet');
			
			// Load stores
			const response = await getStores(
				userLocation.latitude,
				userLocation.longitude,
				'',
				1,
				20
			);
			
			stores = response.data.records || [];
			initMap();
		} catch (error) {
			console.error('Failed to load stores:', error);
		} finally {
			loading = false;
			mapLoading = false;
		}
	});

	onDestroy(() => {
		if (mapInstance) {
			mapInstance.remove();
			mapInstance = null;
			mapLayer = null;
		}
	});
	
	function formatDistance(meters) {
		if (meters < 1000) return `${meters}m`;
		return `${(meters / 1000).toFixed(1)}km`;
	}

	function isOpen(store) {
		return store.businessStatus === 0;
	}
	
	function selectStore(store) {
		if (!isOpen(store)) return;
		selectedStoreStore.set(store);
		goto('/menu');
	}

	function toNumber(value) {
		const n = Number(value);
		return Number.isFinite(n) ? n : null;
	}

	function initMap() {
		if (!leaflet || !mapContainer || mapInstance) return;

		const userLat = toNumber(userLocation?.latitude) || -6.573982;
		const userLng = toNumber(userLocation?.longitude) || 110.684519;

		mapInstance = leaflet.map(mapContainer, {
			zoomControl: false,
			scrollWheelZoom: false,
			attributionControl: true
		}).setView([userLat, userLng], 13);

		leaflet.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
			maxZoom: 19,
			attribution: '&copy; OpenStreetMap contributors'
		}).addTo(mapInstance);

		mapLayer = leaflet.layerGroup().addTo(mapInstance);
		updateMap();
	}

	function markerIcon(kind) {
		const colors = {
			user: '#4d84ff',
			open: '#28b76b',
			closed: '#f25656'
		};

		return leaflet.divIcon({
			className: 'store-map-marker',
			html: `<span style="background:${colors[kind] || colors.open}"></span>`,
			iconSize: [18, 18],
			iconAnchor: [9, 9]
		});
	}

	function updateMap() {
		if (!mapInstance || !mapLayer || !leaflet) return;

		mapLayer.clearLayers();

		const storeBounds = [];
		const userLat = toNumber(userLocation?.latitude);
		const userLng = toNumber(userLocation?.longitude);

		if (userLat !== null && userLng !== null) {
			leaflet.marker([userLat, userLng], { icon: markerIcon('user') })
				.addTo(mapLayer)
				.bindPopup('Lokasi Anda');
		}

		let firstStorePoint = null;
		for (const store of visibleStores) {
			const lat = toNumber(store.latitude);
			const lng = toNumber(store.longitude);
			if (lat === null || lng === null) continue;

			if (!firstStorePoint) {
				firstStorePoint = [lat, lng];
			}

			const marker = leaflet.marker([lat, lng], {
				icon: markerIcon(isOpen(store) ? 'open' : 'closed')
			}).addTo(mapLayer);

			marker.bindPopup(`
				<div style="min-width:160px">
					<strong>${store.storeName}</strong><br />
					<span>${store.storeAddress || ''}</span>
				</div>
			`);

			storeBounds.push([lat, lng]);
		}

		if (storeBounds.length > 1) {
			mapInstance.fitBounds(storeBounds, { padding: [28, 28] });
		} else if (firstStorePoint) {
			mapInstance.setView(firstStorePoint, 15);
		} else if (userLat !== null && userLng !== null) {
			mapInstance.setView([userLat, userLng], 13);
		}
	}
</script>

<svelte:head>
	<title>Outlets - Coffee Order</title>
</svelte:head>

<div class="outlets-screen">
	<div class="outlets-header">
		<div class="page-kicker">Pilih Outlet</div>
		<h1 class="page-title">Outlet Terdekat</h1>
		<div class="map-card">
			<div class="map-card-top">
				<div>
					<div class="map-label">Lokasi Aktif</div>
					<div class="map-value">{currentLabel}</div>
			</div>
			<div class="map-chip"><UiIcon name="location" size={14} /> {stores.length} Outlet</div>
		</div>
		<div class="map-visual">
			{#if mapLoading}
				<div class="map-loading">Memuat peta...</div>
			{/if}
			<div bind:this={mapContainer} class="leaflet-map"></div>
		</div>
	</div>
		<div class="controls-row">
			<div class="filter-tabs">
				<button class:active={activeFilter === 'all'} on:click={() => activeFilter = 'all'}>Semua</button>
				<button class:active={activeFilter === 'open'} on:click={() => activeFilter = 'open'}>Buka ({openCount})</button>
				<button class:active={activeFilter === 'closed'} on:click={() => activeFilter = 'closed'}>Tutup ({closedCount})</button>
			</div>
			<div class="search-box">
				<UiIcon name="search" size={16} />
				<input 
					type="text" 
					placeholder="Cari outlet"
					bind:value={searchQuery}
				/>
			</div>
		</div>
	</div>

	<div class="outlets-content">
		{#if loading}
			<div class="loading-state">
				<div class="spinner"></div>
				<p>Loading outlets...</p>
			</div>
		{:else if visibleStores.length === 0}
			<div class="empty-state">
				<div class="empty-icon"><UiIcon name="location" size={48} /></div>
				<div class="empty-text">Tidak ada outlet ditemukan</div>
			</div>
		{:else}
			{#each visibleStores as store, index}
				<button 
					class="outlet-card" 
					class:disabled={!isOpen(store)}
					style="animation-delay: {index * 0.1}s"
					on:click={() => selectStore(store)}
					disabled={!isOpen(store)}
					type="button"
				>
					<div class="outlet-header">
						<div class="outlet-icon"><UiIcon name="coffee" size={18} /></div>
						<div class="outlet-status" class:open={store.businessStatus === 0}>
							{isOpen(store) ? 'Buka' : 'Tutup'}
						</div>
					</div>

					<div class="outlet-body">
						<h3 class="outlet-name">{store.storeName}</h3>
						{#if !isOpen(store)}
							<div class="closed-note">Outlet sedang tutup</div>
						{/if}
						
						<div class="outlet-info">
						<div class="info-row">
							<span class="info-icon"><UiIcon name="location" size={14} /></span>
							<span class="info-text">{store.storeAddress}</span>
						</div>

						<div class="info-row">
							<span class="info-icon"><UiIcon name="phone" size={14} /></span>
							<span class="info-text">{store.storePhone}</span>
						</div>

						{#if store.distance}
							<div class="info-row">
								<span class="info-icon"><UiIcon name="delivery" size={14} /></span>
								<span class="info-text">{formatDistance(store.distance)}</span>
							</div>
						{/if}
						
						{#if store.isDelivery === 1}
							<div class="info-row">
								<span class="info-icon"><UiIcon name="delivery" size={14} /></span>
								<span class="info-text">Delivery tersedia</span>
							</div>
						{/if}
						</div>
					</div>

						<div class="outlet-footer">
							<button class="select-btn" disabled={!isOpen(store)}>
								Pilih Outlet
							</button>
						</div>
				</button>
			{/each}
		{/if}
	</div>
</div>

<style>
	.outlets-screen {
		display: flex;
		flex-direction: column;
		height: calc(100vh - 68px);
		background: var(--warm-gray);
	}

	.outlets-header {
		background: linear-gradient(180deg, #fffaf5 0%, #ffffff 100%);
		padding: 18px 16px 14px;
		flex-shrink: 0;
		border-bottom: 1px solid #f2e8df;
	}

	.page-kicker {
		font-size: 11px;
		font-weight: 800;
		letter-spacing: 1.4px;
		text-transform: uppercase;
		color: var(--orange);
		margin-bottom: 6px;
	}

	.page-title {
		font-size: 26px;
		font-weight: 800;
		color: var(--text-dark);
		margin: 0 0 14px 0;
		line-height: 1.05;
	}

	.map-card {
		background: linear-gradient(135deg, #f7efe7 0%, #fff 100%);
		border: 1px solid #f0e2d8;
		border-radius: 24px;
		padding: 14px;
		box-shadow: 0 12px 28px rgba(244, 123, 32, 0.08);
		margin-bottom: 14px;
	}

	.map-card-top {
		display: flex;
		justify-content: space-between;
		gap: 12px;
		align-items: flex-start;
		margin-bottom: 12px;
	}

	.map-label {
		font-size: 11px;
		color: var(--text-light);
		font-weight: 700;
	}

	.map-value {
		font-size: 14px;
		font-weight: 800;
		color: var(--text-dark);
		margin-top: 4px;
	}

	.map-chip {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 8px 10px;
		background: white;
		border-radius: 999px;
		font-size: 11px;
		font-weight: 800;
		color: var(--orange);
		box-shadow: 0 6px 16px rgba(0, 0, 0, 0.06);
	}

	.map-visual {
		position: relative;
		height: 170px;
		border-radius: 20px;
		overflow: hidden;
		background: linear-gradient(135deg, #dff2f7 0%, #e9f9ee 50%, #f8efe4 100%);
	}

	.leaflet-map {
		width: 100%;
		height: 100%;
	}

	.map-loading {
		position: absolute;
		inset: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.55);
		backdrop-filter: blur(4px);
		font-size: 13px;
		font-weight: 700;
		color: var(--text-dark);
		z-index: 2;
	}

	.controls-row {
		display: grid;
		gap: 12px;
	}

	.filter-tabs {
		display: flex;
		gap: 8px;
		overflow-x: auto;
		padding-bottom: 2px;
	}

	.filter-tabs button {
		border: none;
		background: #f4f1ee;
		color: var(--text-mid);
		padding: 10px 14px;
		border-radius: 999px;
		font-family: var(--font);
		font-weight: 800;
		font-size: 13px;
		white-space: nowrap;
	}

	.filter-tabs button.active {
		background: var(--orange);
		color: white;
		box-shadow: 0 6px 16px rgba(244, 123, 32, 0.22);
	}

	.search-box {
		background: white;
		border-radius: 999px;
		padding: 12px 14px;
		display: flex;
		align-items: center;
		gap: 8px;
		border: 1px solid #eee1d6;
		box-shadow: 0 10px 20px rgba(0, 0, 0, 0.04);
	}

	.search-box input {
		flex: 1;
		border: none;
		background: none;
		font-size: 14px;
		font-weight: 600;
		color: var(--text-dark);
		font-family: var(--font);
		outline: none;
	}

	.search-box input::placeholder {
		color: #bbb;
	}

	.outlets-content {
		flex: 1;
		overflow-y: auto;
		padding: 16px 12px;
		scrollbar-width: none;
	}

	.outlets-content::-webkit-scrollbar {
		display: none;
	}

	.loading-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 60px 20px;
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
	}

	.outlet-card {
		background: white;
		border-radius: 22px;
		overflow: hidden;
		margin-bottom: 14px;
		transition: transform 0.2s, box-shadow 0.2s, opacity 0.2s;
		cursor: pointer;
		border: 1px solid #f3ece6;
		padding: 0;
		text-align: left;
		font-family: var(--font);
		animation: slideUp 0.5s ease forwards;
		opacity: 0;
		box-shadow: 0 10px 24px rgba(0, 0, 0, 0.05);
	}

	.outlet-card.disabled {
		opacity: 0.55;
		cursor: not-allowed;
		filter: grayscale(0.2);
	}

	@keyframes slideUp {
		to {
			opacity: 1;
			transform: translateY(0);
		}
		from {
			opacity: 0;
			transform: translateY(20px);
		}
	}

	.outlet-card:active {
		transform: scale(0.98);
	}

	.outlet-header {
		background: linear-gradient(135deg, #f6f0ea 0%, #fff8f2 100%);
		padding: 14px 16px;
		display: flex;
		justify-content: space-between;
		align-items: center;
		position: relative;
		overflow: hidden;
	}

	.outlet-icon {
		width: 40px;
		height: 40px;
		border-radius: 14px;
		background: rgba(244, 123, 32, 0.1);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--orange);
		box-shadow: inset 0 0 0 1px rgba(244, 123, 32, 0.12);
		position: relative;
		z-index: 1;
	}

	.outlet-status {
		padding: 6px 12px;
		border-radius: 12px;
		font-weight: 700;
		font-size: 12px;
		backdrop-filter: blur(10px);
		position: relative;
		z-index: 1;
	}

	.outlet-status.open {
		background: rgba(39, 174, 96, 0.9);
		color: white;
	}

	.outlet-status:not(.open) {
		background: rgba(220, 53, 69, 0.9);
		color: white;
	}

	.outlet-body {
		padding: 16px;
	}

	.outlet-name {
		font-size: 16px;
		font-weight: 800;
		color: var(--text-dark);
		margin: 0 0 12px 0;
	}

	.closed-note {
		display: inline-flex;
		align-items: center;
		padding: 6px 10px;
		border-radius: 999px;
		background: #f5f5f5;
		color: #999;
		font-size: 11px;
		font-weight: 800;
		margin-bottom: 10px;
	}

	.outlet-info {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.info-row {
		display: flex;
		align-items: flex-start;
		gap: 10px;
	}

	.info-icon {
		font-size: 16px;
		flex-shrink: 0;
	}

	.info-text {
		color: var(--text-mid);
		font-size: 13px;
		line-height: 1.4;
		font-weight: 600;
	}

	.outlet-footer {
		padding: 12px 16px 16px;
		background: linear-gradient(180deg, #fff 0%, #fcf8f4 100%);
	}

	.select-btn {
		width: 100%;
		padding: 13px;
		background: var(--orange);
		color: white;
		border: none;
		border-radius: 16px;
		font-weight: 700;
		font-size: 14px;
		cursor: pointer;
		font-family: var(--font);
		transition: all 0.2s;
		box-shadow: 0 10px 22px rgba(244, 123, 32, 0.24);
	}

	.select-btn:active {
		transform: scale(0.97);
	}

	.select-btn:disabled {
		background: #e8e8e8;
		color: #aaa;
		box-shadow: none;
		cursor: not-allowed;
	}

	@media (min-width: 768px) {
		.outlets-screen {
			max-width: 520px;
			margin: 0 auto;
			box-shadow: 0 0 40px rgba(0, 0, 0, 0.08);
		}
	}

	@media (min-width: 1024px) {
		.outlets-screen {
			max-width: 620px;
		}
	}
</style>
