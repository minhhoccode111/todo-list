import { api } from '$lib';

export const auth = $state({
	token: api.token.get()
});

export const setAuth = (s: string) => {
	auth.token = s;
	api.token.set(s);
};

export const clearAuth = () => {
	auth.token = '';
	api.token.clear();
};
