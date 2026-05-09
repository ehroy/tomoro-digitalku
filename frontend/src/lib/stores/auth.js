import { writable, get } from 'svelte/store';
import { browser } from '$app/environment';

// Auth store
export const authStore = writable({
	token: null,
	user: null,
	phone: null,
	deviceCode: null,
	wToken: null,
	ucde: null,
	isAuthenticated: false
});

// Load from localStorage on init
if (browser) {
	const saved = localStorage.getItem('tomoro_auth');
	if (saved) {
		try {
			authStore.set(JSON.parse(saved));
		} catch (e) {
			console.error('Failed to parse saved auth:', e);
		}
	}
}

// Save to localStorage on change
authStore.subscribe(value => {
	if (browser) {
		if (value.isAuthenticated) {
			localStorage.setItem('tomoro_auth', JSON.stringify(value));
		} else {
			localStorage.removeItem('tomoro_auth');
		}
	}
});
