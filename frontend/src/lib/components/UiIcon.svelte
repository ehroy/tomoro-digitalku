<script>
	export let name = 'coffee';
	export let size = 20;
	export let strokeWidth = 2.2;
	export let className = '';

	$: viewBox = '0 0 24 24';
</script>

<svg
	class={className}
	width={size}
	height={size}
	viewBox={viewBox}
	fill="none"
	stroke="currentColor"
	stroke-width={strokeWidth}
	stroke-linecap="round"
	stroke-linejoin="round"
	aria-hidden="true"
>
	{#if name === 'coffee'}
		<path d="M5 8h11a3 3 0 0 1 0 6H5z" />
		<path d="M18 9h1a2 2 0 0 1 0 4h-1" />
		<path d="M7 4c0 1.5-1 1.5-1 3" />
		<path d="M11 4c0 1.5-1 1.5-1 3" />
		<path d="M15 4c0 1.5-1 1.5-1 3" />
		<path d="M6 20h10" />
	{/if}

	{#if name === 'voucher'}
		<path d="M5 6h14a2 2 0 0 1 2 2v2a2 2 0 0 0 0 4v2a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-2a2 2 0 0 0 0-4V8a2 2 0 0 1 2-2Z" />
		<path d="M12 8v8" />
		<path d="M10 10a2 2 0 1 1 4 0c0 1.2-.7 1.7-1.8 2.4-.8.5-1.2.9-1.2 1.6" />
		<circle cx="12" cy="16.5" r="0.7" fill="currentColor" stroke="none" />
	{/if}

	{#if name === 'user'}
		<circle cx="12" cy="8" r="3.5" />
		<path d="M4 20a8 8 0 0 1 16 0" />
	{/if}

	{#if name === 'location'}
		<path d="M12 21s6-5.2 6-11a6 6 0 1 0-12 0c0 5.8 6 11 6 11Z" />
		<circle cx="12" cy="10" r="2.25" />
	{/if}

	{#if name === 'refresh'}
		<path d="M20 12a8 8 0 0 0-14.9-4" />
		<path d="M5 4v4h4" />
		<path d="M4 12a8 8 0 0 0 14.9 4" />
		<path d="M19 20v-4h-4" />
	{/if}

	{#if name === 'help'}
		<circle cx="12" cy="12" r="9" />
		<path d="M9.5 9a2.75 2.75 0 1 1 3.9 2.5c-.9.4-1.4 1-1.4 2" />
		<path d="M12 17h.01" />
	{/if}

	{#if name === 'globe'}
		<circle cx="12" cy="12" r="9" />
		<path d="M3 12h18" />
		<path d="M12 3a15 15 0 0 1 0 18" />
		<path d="M12 3a15 15 0 0 0 0 18" />
	{/if}

	{#if name === 'terms'}
		<path d="M7 3h7l5 5v13a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1Z" />
		<path d="M14 3v6h6" />
		<path d="M9 12h6" />
		<path d="M9 16h6" />
	{/if}

	{#if name === 'logout'}
		<path d="M10 17H6a2 2 0 0 1-2-2V9a2 2 0 0 1 2-2h4" />
		<path d="M14 7l5 5-5 5" />
		<path d="M19 12H9" />
	{/if}

	{#if name === 'trash'}
		<path d="M4 7h16" />
		<path d="M9 7V5.5A1.5 1.5 0 0 1 10.5 4h3A1.5 1.5 0 0 1 15 5.5V7" />
		<path d="M7 7l1 13h8l1-13" />
		<path d="M10 11v5" />
		<path d="M14 11v5" />
	{/if}

	{#if name === 'alert'}
		<path d="M12 4 2.8 20h18.4L12 4Z" />
		<path d="M12 9v5" />
		<path d="M12 17h.01" />
	{/if}

	{#if name === 'close'}
		<path d="M6 6l12 12" />
		<path d="M18 6 6 18" />
	{/if}

	{#if name === 'cart'}
		<path d="M5 6h2l1.2 8.5a2 2 0 0 0 2 1.7h6.8a2 2 0 0 0 2-1.6L20 9H8" />
		<circle cx="10" cy="19" r="1.5" />
		<circle cx="17" cy="19" r="1.5" />
	{/if}

	{#if name === 'package'}
		<path d="M21 8.5 12 4 3 8.5 12 13l9-4.5Z" />
		<path d="M3 8.5V17l9 4.5 9-4.5V8.5" />
		<path d="M12 13v8.5" />
	{/if}

	{#if name === 'calendar'}
		<path d="M7 3v3" />
		<path d="M17 3v3" />
		<rect x="3" y="5" width="18" height="16" rx="2" />
		<path d="M3 10h18" />
	{/if}

	{#if name === 'group'}
		<circle cx="9" cy="8" r="3" />
		<circle cx="17" cy="10" r="2.5" />
		<path d="M3.5 20a5.5 5.5 0 0 1 11 0" />
		<path d="M13 20a4.5 4.5 0 0 1 7 0" />
	{/if}

	{#if name === 'activity'}
		<path d="M3 12h4l2-5 4 10 2-5h6" />
	{/if}

	{#if name === 'delivery'}
		<path d="M3 7h11v10H3z" />
		<path d="M14 10h3l3 3v4h-6" />
		<circle cx="7" cy="18" r="1.5" />
		<circle cx="17" cy="18" r="1.5" />
	{/if}

	{#if name === 'phone'}
		<path d="M8 3h3l1 5-2 1a12 12 0 0 0 5 5l1-2 5 1v3a2 2 0 0 1-2 2A18 18 0 0 1 6 5a2 2 0 0 1 2-2Z" />
	{/if}

	{#if name === 'check'}
		<path d="M20 6 9 17l-5-5" />
	{/if}

	{#if name === 'clock'}
		<circle cx="12" cy="12" r="9" />
		<path d="M12 7v5l3 2" />
	{/if}

	{#if name === 'points'}
		<circle cx="12" cy="12" r="8" />
		<path d="M12 7v10" />
		<path d="M7.5 11.5h9" />
	{/if}

	{#if name === 'qr'}
		<path d="M4 4h6v6H4zM14 4h6v6h-6zM4 14h6v6H4z" />
		<path d="M14 14h2v2h-2zM18 14h2v2h-2zM14 18h6" />
	{/if}

	{#if name === 'search'}
		<circle cx="11" cy="11" r="6.5" />
		<path d="m16 16 4 4" />
	{/if}
</svg>
