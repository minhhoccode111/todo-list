import { env } from '$env/dynamic/public';
import { Api } from '$lib/types/api';
import ky, { HTTPError } from 'ky';
import type { ResponseMessage } from '$lib/types/api';

export const getApiError = async (error: unknown): Promise<ResponseMessage> => {
	if (error instanceof HTTPError) {
		try {
			return await error.response.json();
		} catch {
			return { message: error.message };
		}
	}
	return { message: 'An unexpected error occurred' };
};

const TOKEN = 'token';

export const getAuthToken = () =>
	typeof window !== 'undefined' ? localStorage.getItem(TOKEN) : null;

export const setAuthToken = (token: string) => {
	if (typeof window !== 'undefined') {
		localStorage.setItem(TOKEN, token);
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
	customFetch: ky.extend({
		hooks: {
			beforeRequest: [
				(request) => {
					const token = getAuthToken();
					if (token) {
						request.headers.set('Authorization', `Bearer ${token}`);
					}
				}
			]
		}
	}),
	securityWorker: async () => {
		return {};
	}
});

export { api };
