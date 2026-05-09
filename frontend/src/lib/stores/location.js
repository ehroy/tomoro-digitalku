import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Location store
export const locationStore = writable({
	latitude: null,
	longitude: null,
	error: null
});

// Get user location
export async function getUserLocation() {
	if (!browser || !navigator.geolocation) {
		return { latitude: -6.573982, longitude: 110.684519 }; // Default: Jepara
	}

	return new Promise((resolve) => {
		navigator.geolocation.getCurrentPosition(
			(position) => {
				const location = {
					latitude: position.coords.latitude,
					longitude: position.coords.longitude
				};
				locationStore.set({ ...location, error: null });
				resolve(location);
			},
			(error) => {
				console.error('Geolocation error:', error);
				const defaultLocation = { latitude: -6.573982, longitude: 110.684519 };
				locationStore.set({ ...defaultLocation, error: error.message });
				resolve(defaultLocation);
			}
		);
	});
}
