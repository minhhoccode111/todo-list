import { env } from '$env/dynamic/public';
import { Api } from '$lib/types/api';
import ky, { HTTPError } from 'ky';
import type { ResponseMessage, ResponseAuth } from '$lib/types/api';

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
	if (typeof window !== 'undefined') localStorage.setItem(TOKEN, token);
};

export const clearAuthToken = () => {
	if (typeof window !== 'undefined') localStorage.removeItem(TOKEN);
};

const baseUrl = env.PUBLIC_API_URL;
let isRefreshing = false;

const apiKy = ky.extend({
	hooks: {
		beforeRequest: [
			(request) => {
				const token = getAuthToken();
				if (token) {
					request.headers.set('Authorization', `Bearer ${token}`);
				}
			}
		],
		afterResponse: [
			async (request, _options, response) => {
				if (response.status !== 401) return response;
				if (isRefreshing) return response;
				if (request.url.endsWith('/refresh')) return response;

				isRefreshing = true;
				try {
					const refreshRes = await fetch(`${baseUrl}/refresh`, {
						method: 'POST',
						credentials: 'include'
					});

					if (!refreshRes.ok) {
						clearAuthToken();
						if (typeof window !== 'undefined') {
							window.location.href = '/login';
						}
						return response;
					}

					const data = (await refreshRes.json()) as ResponseAuth;
					setAuthToken(data.token);

					const retryReq = request.clone();
					retryReq.headers.set('Authorization', `Bearer ${data.token}`);
					return ky(retryReq);
				} catch {
					clearAuthToken();
					if (typeof window !== 'undefined') {
						window.location.href = '/login';
					}
					return response;
				} finally {
					isRefreshing = false;
				}
			}
		]
	}
});

const api = new Api({
	baseUrl,
	baseApiParams: {
		credentials: 'include'
	},
	customFetch: apiKy,
	securityWorker: async () => {
		return {};
	}
});

export { api };
