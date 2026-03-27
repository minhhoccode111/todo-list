import { env } from '$env/dynamic/public';
import { Api } from '$lib/types/api';

const TOKEN = 'token';

export const getAuthToken = () =>
	typeof window !== 'undefined' ? localStorage.getItem(TOKEN) : null;

export const setAuthToken = (token?: string) => {
	if (typeof window !== 'undefined') {
		localStorage.setItem(TOKEN, token ?? '');
	}
};

export const clearAuthToken = () => {
	if (typeof window !== 'undefined') {
		localStorage.removeItem(TOKEN);
	}
};

const api = new Api({
	baseUrl: env.PUBLIC_API_URL,
	baseApiParams: {
		credentials: 'same-origin'
	},
	securityWorker: async () => {
		const token = getAuthToken();
		if (token) {
			return {
				headers: {
					Authorization: `Bearer ${token}`
				}
			};
		}
		return {};
	}
});

export { api };
