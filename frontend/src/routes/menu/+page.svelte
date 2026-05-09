<script>
	import { onMount } from 'svelte';
	import { getMenu } from '$lib/api/menu.js';
	import { addToCart as addToCartAPI, editCart as editCartAPI, getCart } from '$lib/api/cart.js';
	import { cart, addToCart, removeFromCart, updateCartQuantity, clearCart, getCartTotal, setCartFromResponse } from '$lib/stores/cart.js';
	import { selectedStoreStore } from '$lib/stores/voucher.js';
	import { goto } from '$app/navigation';
	import UiIcon from '$lib/components/UiIcon.svelte';

	let menuCategories = [];
	let loading = true;
	let selectedCategory = null;
	let showCart = false;
	let orderSuccess = false;
	let selectedStore = null;
	let showSizeModal = false;
	let selectedProduct = null;
	let selectedSize = null;

	const IMAGE_BASE = 'https://api-service.tomoro-coffee.id';

	$: cartTotal = getCartTotal($cart);
	$: cartCount = $cart.reduce((sum, item) => sum + Number(item.quantity || item.amount || 0), 0);

	function normalizeImageUrl(url) {
		if (!url) return '';
		if (url.startsWith('http://') || url.startsWith('https://') || url.startsWith('data:')) {
			return url;
		}

		return `${IMAGE_BASE}${url.startsWith('/') ? '' : '/'}${url}`;
	}

	function firstImage(...values) {
		for (const value of values) {
			if (!value) continue;

			const candidates = Array.isArray(value)
				? value
				: String(value)
					.split(',')
					.map((item) => item.trim())
					.filter(Boolean);

			for (const candidate of candidates) {
				const normalized = normalizeImageUrl(candidate);
				if (normalized) return normalized;
			}
		}

		return '';
	}

	function getMenuImage(category) {
		return firstImage(category?.menuPicture);
	}

	function getProductImage(product) {
		return firstImage(
			product?.picture?.max,
			product?.pictureMaxUrls,
			product?.picture?.main,
			product?.pictureUrls,
			product?.pictureMiniUrls,
			product?.picture?.mini
		);
	}

	function getSizeOptions(product) {
		const variants = Array.isArray(product?.pluCodes) ? product.pluCodes : [];
		const defaultSize = String(product?.defaultSelectName || 'Small');

		if (variants.length > 0) {
			return variants.map((variant, index) => ({
				key: variant.code || `${product.code}-${index}`,
				label: variant.name || `Size ${index + 1}`,
				price: variant.price || product.price,
				linePrice: variant.linePrice || product.linePrice || 0,
				pluCode: variant.code || product.pluCode || '',
				recommended: index === 0
			}));
		}

		const basePrice = Number(product?.price || 0);

		return [
			{
				key: 'small',
				label: 'Small',
				price: basePrice,
				linePrice: product?.linePrice || 0,
				pluCode: product?.pluCode || '',
				recommended: defaultSize.toLowerCase().includes('small')
			},
			{
				key: 'medium',
				label: 'Medium',
				price: basePrice + 3000,
				linePrice: basePrice + 3000,
				pluCode: product?.pluCode || '',
				recommended: defaultSize.toLowerCase().includes('medium')
			},
			{
				key: 'large',
				label: 'Large',
				price: basePrice + 6000,
				linePrice: basePrice + 6000,
				pluCode: product?.pluCode || '',
				recommended: defaultSize.toLowerCase().includes('large')
			}
		];
	}

	function openSizeModal(product) {
		selectedProduct = product;
		const options = getSizeOptions(product);
		selectedSize = options.find((option) => option.recommended) || options[0];
		showSizeModal = true;
	}

	function closeSizeModal() {
		showSizeModal = false;
		selectedProduct = null;
		selectedSize = null;
	}

	async function confirmAddToCart() {
		if (!selectedProduct || !selectedSize) return;

		const cartItem = {
			product_id: selectedProduct.code,
			name: `${selectedProduct.name} - ${selectedSize.label}`,
			size_label: selectedSize.label,
			size_code: selectedSize.key,
			price: selectedSize.price,
			subtotal: selectedSize.price,
			cart_key: `${selectedProduct.code}:${selectedSize.key}`
		};

		try {
			addToCart(cartItem);
			await addToCartAPI({
				amount: 1,
				storeCode: selectedStore.storeCode,
				itemType: selectedProduct.type ?? 1,
				itemCode: selectedProduct.code,
				menuCode: selectedProduct.menuCode,
				pluCode: selectedSize.pluCode || selectedProduct.pluCode || '',
				mainMenuType: 1,
				productCode: selectedProduct.code,
				productName: selectedProduct.name,
				quantity: 1,
				price: selectedSize.price,
				size: selectedSize.label
			});

			const cartResponse = await getCart(selectedStore.storeCode, 1);
			setCartFromResponse(cartResponse.data, [cartItem], true);

			showCart = true;
			closeSizeModal();
		} catch (error) {
			console.error('Failed to add to cart:', error);
		}
	}

	async function removeCartItem(item) {
		if (!selectedStore) return;

		try {
			await editCartAPI({
				key: item.cart_key,
				storeCode: selectedStore.storeCode,
				operateType: 1,
				mainMenuType: 1,
				amount: 1,
				unitItemAmount: 1
			});
			removeFromCart(item.cart_key || item.product_id);
			const cartResponse = await getCart(selectedStore.storeCode, 1);
			setCartFromResponse(cartResponse.data, $cart, true);
		} catch (error) {
			console.error('Failed to remove cart item:', error);
		}
	}

	async function changeCartQuantity(item, nextQuantity) {
		if (!selectedStore) return;

		try {
			if (nextQuantity <= 0) {
				await editCartAPI({
					key: item.cart_key,
					storeCode: selectedStore.storeCode,
					operateType: 1,
					mainMenuType: 1,
					amount: 1,
					unitItemAmount: 1
				});
				removeFromCart(item.cart_key || item.product_id);
			} else {
				await editCartAPI({
					key: item.cart_key,
					storeCode: selectedStore.storeCode,
					operateType: 2,
					mainMenuType: 1,
					amount: nextQuantity,
					unitItemAmount: 1
				});
				updateCartQuantity(item.cart_key || item.product_id, nextQuantity);
			}

			const cartResponse = await getCart(selectedStore.storeCode, 1);
			setCartFromResponse(cartResponse.data, $cart, true);
		} catch (error) {
			console.error('Failed to update cart quantity:', error);
		}
	}

	onMount(async () => {
		// Check if store is selected
		selectedStore = $selectedStoreStore;
		
		if (!selectedStore) {
			goto('/outlets');
			return;
		}

		try {
			const response = await getMenu(selectedStore.storeCode);
			menuCategories = response.data.menuVos || [];
			const cartResponse = await getCart(selectedStore.storeCode, 1);
			setCartFromResponse(cartResponse.data);
			
			if (menuCategories.length > 0) {
				selectedCategory = menuCategories[0].menuCode;
			}
		} catch (error) {
			console.error('Failed to load menu:', error);
		} finally {
			loading = false;
		}
	});

	$: currentCategory = menuCategories.find(cat => cat.menuCode === selectedCategory);
	$: products = currentCategory?.items || [];

	function handleAddToCart(product) {
		openSizeModal(product);
	}

	function formatPrice(price) {
		return new Intl.NumberFormat('id-ID', {
			style: 'currency',
			currency: 'IDR',
			minimumFractionDigits: 0
		}).format(price);
	}

	function handleCheckout() {
		goto('/checkout');
	}
</script>

<svelte:head>
	<title>Menu - Coffee Order</title>
</svelte:head>

<div class="menu-screen">
	{#if orderSuccess}
		<div class="success-toast">
			<UiIcon name="check" size={16} /> Item ditambahkan ke keranjang!
		</div>
	{/if}

	<div class="menu-top">
		<div class="store-info">
			<button class="back-btn" on:click={() => goto('/outlets')}>
				<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
					<path d="m15 18-6-6 6-6"/>
				</svg>
			</button>
			<div class="store-details">
				<h2>{selectedStore?.storeName || 'Pilih Outlet'}</h2>
				<p>{selectedStore?.storeAddress || ''}</p>
			</div>
		</div>
	</div>

	<div class="menu-body">
		<div class="menu-sidebar">
			{#each menuCategories as category}
				<div 
					class="sidebar-item" 
					class:active={selectedCategory === category.menuCode}
					on:click={() => selectedCategory = category.menuCode}
					on:keypress={() => selectedCategory = category.menuCode}
					role="button"
					tabindex="0"
				>
					{#if category.menuDocument}
						<div class="sidebar-badge">{category.menuDocument}</div>
					{/if}
					{category.menuName}
				</div>
			{/each}
		</div>

		<div class="menu-content">
			{#if loading}
				<div class="loading-state">
					<div class="spinner"></div>
					<p>Loading menu...</p>
				</div>
		{:else if currentCategory}
				{#if getMenuImage(currentCategory)}
			<div class="menu-banner">
				<img src={getMenuImage(currentCategory)} alt={currentCategory.menuName} />
						<div class="banner-overlay">
							<h3>{currentCategory.menuName}</h3>
							{#if currentCategory.menuDesc}
								<p>{currentCategory.menuDesc}</p>
							{/if}
						</div>
					</div>
				{/if}

				<div class="section-title">{currentCategory.menuName}</div>

				{#each products as product}
					<div class="product-card" class:sold-out={product.isSellOut === 0}>
						<div class="product-img">
							{#if getProductImage(product)}
								<img src={getProductImage(product)} alt={product.name} loading="lazy" />
							{:else}
								<div class="placeholder-img"><UiIcon name="coffee" size={24} /></div>
							{/if}
						</div>
						<div class="product-info">
							<div class="product-name">{product.name}</div>
							{#if product.desc}
								<div class="product-desc">{product.desc}</div>
							{/if}
							<div class="product-pricing">
								<span class="product-price">{formatPrice(product.price)}</span>
								{#if product.linePrice && product.linePrice > product.price}
									<span class="product-price-old">{formatPrice(product.linePrice)}</span>
								{/if}
							</div>
							{#if product.menuSecondaryTag}
								<div class="product-tag">{product.menuSecondaryTag}</div>
							{/if}
						</div>
						<button 
							class="add-btn" 
							disabled={product.isSellOut === 0 || product.isCanAddCart === 0}
							on:click={() => handleAddToCart(product)}
						>
							+
						</button>
					</div>
				{/each}
			{/if}
		</div>
	</div>

	{#if cartCount > 0}
		<button class="cart-float" on:click={() => showCart = !showCart}>
			<UiIcon name="cart" size={18} /> {cartCount} Item{cartCount > 1 ? 's' : ''}
		</button>
	{/if}

				{#if showCart}
		<div class="cart-overlay" on:click={() => showCart = false} role="presentation"></div>
		<div class="cart-panel">
			<div class="cart-header">
				<h2>Keranjang</h2>
				<button class="close-btn" on:click={() => showCart = false}><UiIcon name="close" size={18} /></button>
			</div>

			<div class="cart-items">
				{#if $cart.length === 0}
					<p class="empty-cart">Keranjang kosong</p>
				{:else}
					{#each $cart as item}
						<div class="cart-item">
							<div class="item-info">
								<h4>{item.name}</h4>
								{#if item.size_label}
									<p class="item-size">Ukuran: {item.size_label}</p>
								{/if}
								<p class="item-price">{formatPrice(item.price)}</p>
							</div>
					<div class="item-controls">
						<button 
							class="qty-btn"
							on:click={() => changeCartQuantity(item, Number(item.quantity || item.amount || 0) - 1)}
						>
							−
						</button>
						<span class="quantity">{item.quantity || item.amount || 0}</span>
					<button 
						class="qty-btn"
						on:click={() => changeCartQuantity(item, Number(item.quantity || item.amount || 0) + 1)}
						>
							+
						</button>
					<button 
						class="remove-btn"
						on:click={() => removeCartItem(item)}
					>
						<UiIcon name="trash" size={18} />
					</button>
							</div>
							<div class="item-subtotal">{formatPrice(item.subtotal)}</div>
						</div>
					{/each}
				{/if}
			</div>

			{#if $cart.length > 0}
				<div class="cart-footer">
					<div class="total-row">
						<span class="total-label">Total:</span>
						<span class="total-amount">{formatPrice(cartTotal)}</span>
					</div>
					<button class="checkout-btn" on:click={handleCheckout}>
						Checkout
					</button>
			</div>
		{/if}
		</div>
	{/if}

	{#if showSizeModal}
		<div class="size-overlay" on:click={closeSizeModal} role="presentation"></div>
		<div class="size-modal">
			<div class="size-modal-header">
				<div>
					<div class="size-modal-title">Pilih Ukuran</div>
					<div class="size-modal-name">{selectedProduct?.name}</div>
				</div>
				<button class="close-btn" on:click={closeSizeModal}><UiIcon name="close" size={18} /></button>
			</div>

			<div class="size-options">
				{#each getSizeOptions(selectedProduct) as option}
					<button class="size-option" class:active={selectedSize?.key === option.key} on:click={() => selectedSize = option}>
						<div class="size-option-top">
							<span>{option.label}</span>
							{#if option.recommended}
								<small>Rekomendasi</small>
							{/if}
						</div>
						<div class="size-option-price">{formatPrice(option.price)}</div>
					</button>
				{/each}
			</div>

			<button class="confirm-size-btn" on:click={confirmAddToCart}>
				Tambah ke Keranjang
			</button>
		</div>
	{/if}
</div>

<style>
	.menu-screen {
		display: flex;
		flex-direction: column;
		height: calc(100vh - 68px);
		position: relative;
	}

	.success-toast {
		position: fixed;
		top: 20px;
		left: 50%;
		transform: translateX(-50%);
		background: #27ae60;
		color: white;
		padding: 12px 24px;
		border-radius: 20px;
		font-weight: 700;
		font-size: 14px;
		z-index: 1000;
		animation: slideDown 0.3s ease;
		box-shadow: 0 4px 20px rgba(39, 174, 96, 0.3);
	}

	@keyframes slideDown {
		from {
			transform: translateX(-50%) translateY(-100%);
			opacity: 0;
		}
		to {
			transform: translateX(-50%) translateY(0);
			opacity: 1;
		}
	}

	.menu-top {
		background: white;
		padding: 16px;
		flex-shrink: 0;
		border-bottom: 1px solid #f0f0f0;
	}

	.store-info {
		display: flex;
		align-items: center;
		gap: 12px;
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
		flex-shrink: 0;
	}

	.store-details h2 {
		margin: 0;
		font-size: 16px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.store-details p {
		margin: 4px 0 0 0;
		font-size: 12px;
		color: var(--text-light);
	}

	.menu-body {
		display: flex;
		flex: 1;
		overflow: hidden;
	}

	.menu-sidebar {
		width: 80px;
		background: white;
		overflow-y: auto;
		scrollbar-width: none;
		border-right: 1px solid #f0f0f0;
		flex-shrink: 0;
	}

	.menu-sidebar::-webkit-scrollbar {
		display: none;
	}

	.sidebar-item {
		padding: 14px 8px;
		text-align: center;
		font-size: 10px;
		font-weight: 700;
		color: var(--text-light);
		cursor: pointer;
		border-left: 3px solid transparent;
		line-height: 1.2;
		transition: all 0.2s;
	}

	.sidebar-item.active {
		color: var(--orange);
		border-left-color: var(--orange);
		background: var(--orange-light);
	}

	.sidebar-badge {
		background: var(--orange);
		color: white;
		font-size: 8px;
		font-weight: 800;
		padding: 2px 5px;
		border-radius: 6px;
		display: inline-block;
		margin-bottom: 3px;
	}

	.menu-content {
		flex: 1;
		overflow-y: auto;
		padding: 16px;
		scrollbar-width: none;
		background: var(--warm-gray);
	}

	.menu-content::-webkit-scrollbar {
		display: none;
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

	.menu-banner {
		position: relative;
		border-radius: 14px;
		overflow: hidden;
		margin-bottom: 16px;
		height: 120px;
	}

	.menu-banner img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.banner-overlay {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		padding: 12px;
		background: linear-gradient(to top, rgba(0,0,0,0.7), transparent);
		color: white;
	}

	.banner-overlay h3 {
		margin: 0;
		font-size: 16px;
		font-weight: 800;
	}

	.banner-overlay p {
		margin: 4px 0 0 0;
		font-size: 11px;
	}

	.section-title {
		font-size: 14px;
		font-weight: 800;
		color: var(--text-dark);
		margin-bottom: 10px;
	}

	.product-card {
		background: white;
		border-radius: var(--card-radius);
		overflow: hidden;
		margin-bottom: 10px;
		display: flex;
		align-items: center;
		padding: 12px;
		gap: 12px;
		transition: transform 0.15s;
	}

	.product-card:active {
		transform: scale(0.98);
	}

	.product-card.sold-out {
		opacity: 0.6;
	}

	.product-img {
		width: 70px;
		height: 70px;
		border-radius: 10px;
		flex-shrink: 0;
		overflow: hidden;
		background: var(--warm-gray);
	}

	.product-img img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.placeholder-img {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 32px;
	}

	.product-info {
		flex: 1;
	}

	.product-name {
		font-size: 13px;
		font-weight: 700;
		color: var(--text-dark);
		margin-bottom: 3px;
	}

	.product-desc {
		font-size: 11px;
		color: var(--text-light);
		margin-bottom: 6px;
		line-height: 1.3;
	}

	.product-pricing {
		display: flex;
		align-items: center;
		gap: 6px;
	}

	.product-price {
		font-size: 15px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.product-price-old {
		font-size: 11px;
		color: var(--text-light);
		text-decoration: line-through;
	}

	.product-tag {
		display: inline-block;
		margin-top: 4px;
		padding: 2px 8px;
		background: #f5f5f5;
		color: #999;
		font-size: 10px;
		font-weight: 700;
		border-radius: 6px;
	}

	.add-btn {
		width: 32px;
		height: 32px;
		background: var(--orange);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
		font-size: 20px;
		cursor: pointer;
		flex-shrink: 0;
		border: none;
		font-family: var(--font);
		transition: transform 0.2s;
	}

	.add-btn:active:not(:disabled) {
		transform: scale(0.9);
	}

	.add-btn:disabled {
		background: #ccc;
		cursor: not-allowed;
	}

	.cart-float {
		position: fixed;
		bottom: 80px;
		right: 20px;
		background: var(--orange);
		color: white;
		border-radius: 20px;
		padding: 10px 16px;
		font-size: 13px;
		font-weight: 800;
		display: flex;
		align-items: center;
		gap: 6px;
		box-shadow: 0 4px 14px rgba(244, 123, 32, 0.4);
		z-index: 99;
		border: none;
		cursor: pointer;
		font-family: var(--font);
		transition: transform 0.2s;
	}

	.cart-float:active {
		transform: scale(0.95);
	}

	.cart-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 98;
		animation: fadeIn 0.3s ease;
	}

	@keyframes fadeIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}

	.cart-panel {
		position: fixed;
		top: 0;
		right: 0;
		width: 100%;
		max-width: 450px;
		height: calc(100vh - 68px);
		bottom: 68px;
		background: white;
		box-shadow: -5px 0 30px rgba(0, 0, 0, 0.2);
		z-index: 99;
		display: flex;
		flex-direction: column;
		animation: slideIn 0.3s ease;
	}

	@keyframes slideIn {
		from { transform: translateX(100%); }
		to { transform: translateX(0); }
	}

	.cart-header {
		padding: 20px;
		border-bottom: 2px solid #f0f0f0;
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: var(--orange);
		color: white;
	}

	.cart-header h2 {
		margin: 0;
		font-size: 18px;
		font-weight: 800;
	}

	.close-btn {
		background: none;
		border: none;
		color: white;
		font-size: 24px;
		cursor: pointer;
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		transition: all 0.2s;
	}

	.close-btn:hover {
		background: rgba(255, 255, 255, 0.2);
	}

	.cart-items {
		flex: 1;
		overflow-y: auto;
		padding: 16px;
		scrollbar-width: none;
	}

	.cart-items::-webkit-scrollbar {
		display: none;
	}

	.empty-cart {
		text-align: center;
		color: var(--text-light);
		padding: 40px 20px;
		font-size: 14px;
	}

	.cart-item {
		background: var(--warm-gray);
		border-radius: 14px;
		padding: 14px;
		margin-bottom: 12px;
	}

	.item-info h4 {
		margin: 0 0 4px 0;
		color: var(--text-dark);
		font-size: 14px;
		font-weight: 700;
	}

	.item-price {
		color: var(--text-mid);
		margin: 0;
		font-size: 12px;
		font-weight: 600;
	}

	.item-controls {
		display: flex;
		align-items: center;
		gap: 8px;
		margin-top: 10px;
	}

	.qty-btn {
		width: 32px;
		height: 32px;
		border: 2px solid var(--orange);
		background: white;
		color: var(--orange);
		border-radius: 8px;
		cursor: pointer;
		font-size: 16px;
		font-weight: 700;
		transition: all 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.qty-btn:active {
		background: var(--orange);
		color: white;
	}

	.quantity {
		min-width: 30px;
		text-align: center;
		font-weight: 700;
		color: var(--text-dark);
		font-size: 14px;
	}

	.remove-btn {
		margin-left: auto;
		background: none;
		border: none;
		font-size: 18px;
		cursor: pointer;
		padding: 4px;
		transition: transform 0.2s;
	}

	.remove-btn:active {
		transform: scale(1.2);
	}

	.item-subtotal {
		text-align: right;
		font-size: 16px;
		font-weight: 800;
		color: var(--orange);
		margin-top: 8px;
		padding-top: 8px;
		border-top: 1px solid #e0e0e0;
	}

	.item-size {
		font-size: 11px;
		color: var(--text-light);
		margin: 4px 0 0;
	}

	.size-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.45);
		z-index: 1200;
	}

	.size-modal {
		position: fixed;
		left: 50%;
		bottom: 88px;
		transform: translateX(-50%);
		width: calc(100% - 24px);
		max-width: 480px;
		background: white;
		border-radius: 24px;
		padding: 18px;
		z-index: 1201;
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.25);
	}

	.size-modal-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 12px;
		margin-bottom: 16px;
	}

	.size-modal-title {
		font-size: 18px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.size-modal-name {
		font-size: 12px;
		color: var(--text-light);
		margin-top: 4px;
	}

	.size-options {
		display: grid;
		gap: 10px;
		margin-bottom: 16px;
	}

	.size-option {
		border: 1.5px solid #ececec;
		border-radius: 16px;
		padding: 12px 14px;
		background: #fff;
		text-align: left;
		font-family: var(--font);
	}

	.size-option.active {
		border-color: var(--orange);
		background: rgba(244, 123, 32, 0.06);
	}

	.size-option-top {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 10px;
		font-weight: 700;
		color: var(--text-dark);
	}

	.size-option-top small {
		font-size: 11px;
		color: var(--orange);
	}

	.size-option-price {
		margin-top: 4px;
		font-size: 13px;
		font-weight: 800;
		color: var(--text-dark);
	}

	.confirm-size-btn {
		width: 100%;
		padding: 14px;
		border: none;
		border-radius: 16px;
		background: var(--orange);
		color: white;
		font-weight: 800;
		font-family: var(--font);
	}

	.cart-footer {
		padding: 20px;
		border-top: 2px solid #f0f0f0;
		background: var(--cream);
		padding-bottom: calc(20px + env(safe-area-inset-bottom));
	}

	.total-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 16px;
	}

	.total-label {
		font-size: 16px;
		font-weight: 700;
		color: var(--text-dark);
	}

	.total-amount {
		font-size: 20px;
		font-weight: 800;
		color: var(--orange);
	}

	.checkout-btn {
		width: 100%;
		padding: 14px;
		background: var(--orange);
		color: white;
		border: none;
		border-radius: 20px;
		font-size: 15px;
		font-weight: 800;
		cursor: pointer;
		transition: all 0.2s;
		font-family: var(--font);
		box-shadow: 0 4px 14px rgba(244, 123, 32, 0.3);
	}

	.checkout-btn:active {
		transform: scale(0.98);
	}

	@media (min-width: 768px) {
		.cart-float {
			right: calc(50% - 240px + 20px);
		}

		.cart-panel {
			right: calc(50% - 240px);
			max-width: 480px;
			width: 480px;
		}

		.size-modal {
			bottom: 88px;
		}
	}

	@media (min-width: 1024px) {
		.cart-float {
			right: calc(50% - 300px + 20px);
		}

		.cart-panel {
			right: calc(50% - 300px);
			max-width: 600px;
			width: 600px;
		}
	}
</style>
