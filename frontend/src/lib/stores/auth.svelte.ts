import { api } from '$lib';
import { get } from 'svelte/store';

const createAuthStore = () => {
	const auth = $state({
		token: '',
		isAuth: false
	});

	return {};
};

export const auth = $state(api.token.get());
