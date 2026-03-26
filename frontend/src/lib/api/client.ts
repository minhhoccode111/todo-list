import { env } from '$env/dynamic/public';

const TOKEN = 'token';

const getAuthToken = () => (typeof window !== 'undefined' ? localStorage.getItem(TOKEN) : null);

const setAuthToken = (token: string) => {
	if (typeof window !== 'undefined') {
		localStorage.setItem(TOKEN, token);
	}
};

const clearAuthToken = () => {
	if (typeof window !== 'undefined') {
		localStorage.removeItem(TOKEN);
	}
};

interface RequestOptions {
	method?: string;
	body?: unknown;
	headers?: Record<string, string>;
}

async function request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
	const { method = 'GET', body, headers = {} } = options;

	const token = getAuthToken();
	const authHeaders: Record<string, string> = token ? { Authorization: `Bearer ${token}` } : {};

	const response = await fetch(`${env.PUBLIC_API_URL}${endpoint}`, {
		method,
		headers: {
			'Content-Type': 'application/json',
			...authHeaders,
			...headers
		},
		body: body ? JSON.stringify(body) : undefined
	});

	if (!response.ok) {
		const error = await response.json().catch(() => ({ message: 'Request failed' }));
		throw new Error(error.message || `HTTP ${response.status}`);
	}

	return response.json();
}

export const api = {
	auth: {
		register: (data: { email: string; password: string; name: string }) =>
			request<{ token: string; user: { id: string; email: string; name: string } }>('/register', {
				method: 'POST',
				body: data
			}),
		login: (data: { email: string; password: string }) =>
			request<{ token: string }>('/login', {
				method: 'POST',
				body: data
			})
	},

	todos: {
		create: (data: { title: string; completed?: boolean }) =>
			request<{ id: string; title: string; completed: boolean; createdAt: string }>('/todos', {
				method: 'POST',
				body: data
			}),

		update: (id: string, data: { title?: string; completed?: boolean }) =>
			request<{ id: string; title: string; completed: boolean; createdAt: string }>(
				`/todos/${id}`,
				{
					method: 'PUT',
					body: data
				}
			),

		delete: (id: string) => request<{ message: string }>(`/todos/${id}`, { method: 'DELETE' }),

		list: (params: { page?: number; limit?: number } = {}) => {
			const query = new URLSearchParams();
			if (params.page) query.set('page', String(params.page));
			if (params.limit) query.set('limit', String(params.limit));
			const queryString = query.toString();
			return request<{
				data: Array<{ id: string; title: string; completed: boolean; createdAt: string }>;
				meta: { page: number; limit: number; total: number; totalPages: number };
			}>(`/todos?${queryString}`);
		}
	},

	token: {
		set: setAuthToken,
		clear: clearAuthToken,
		get: getAuthToken
	}
};
