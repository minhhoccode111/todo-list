import { api } from '$lib';
import type { ResponseAuth } from '$lib/types/api';
import { toast } from 'svelte-sonner';

export const auth = $state<ResponseAuth>({
	token: api.token.get() ?? ''
});

export const setAuth = (s?: string) => {
	auth.token = s;
	api.token.set(s);
	toast.success('welcome');
};

export const clearAuth = () => {
	auth.token = '';
	api.token.clear();
	toast.info('goodbye');
};
