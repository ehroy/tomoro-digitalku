import { api } from './client.js';
import { authStore } from '$lib/stores/auth.js';

export async function login(phone, pin) {
	const response = await api.request('/auth/login', {
		method: 'POST',
		body: JSON.stringify({ 
			phone, 
			pin 
		})
	});
	
	// Update auth store
	authStore.set({
		token: response.data.token,
		user: response.data.user,
		deviceCode: response.data.deviceCode,
		wToken: response.data.wToken,
		ucde: response.data.ucde,
		phone,
		isAuthenticated: true
	});
	
	return response;
}

export function logout() {
	authStore.set({
		token: null,
		user: null,
		deviceCode: null,
		wToken: null,
		ucde: null,
		phone: null,
		isAuthenticated: false
	});
}
