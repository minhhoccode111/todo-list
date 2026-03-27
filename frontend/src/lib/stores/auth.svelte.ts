import { getAuthToken, setAuthToken, clearAuthToken } from '$lib';
import type { ResponseAuth } from '$lib/types/api';
import { toast } from 'svelte-sonner';

export const auth = $state<ResponseAuth>({
	token: getAuthToken() ?? ''
});

export const setAuth = (s?: string) => {
	auth.token = s ?? '';
	setAuthToken(s);
	toast.success('welcome');
};

export const clearAuth = () => {
	auth.token = '';
	clearAuthToken();
	toast.info('goodbye');
};
