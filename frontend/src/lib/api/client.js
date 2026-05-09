import { get } from "svelte/store";
import { authStore } from "$lib/stores/auth.js";

const API_BASE = "/api";

class APIClient {
  async request(endpoint, options = {}) {
    const auth = get(authStore);

    const headers = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    // Add auth headers if available
    if (auth.token) {
      headers["Authorization"] = `Bearer ${auth.token}`;
    }
    if (auth.deviceCode) {
      headers["X-Device-Code"] = auth.deviceCode;
    }
    if (auth.wToken) {
      headers["X-WToken"] = auth.wToken;
    }
    if (auth.ucde) {
      headers["X-UCDE"] = auth.ucde;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    const data = await response.json();

    if (!response.ok || !data.success) {
      throw new Error(data.error || `API Error: ${response.status}`);
    }

    return data;
  }
}

export const api = new APIClient();
