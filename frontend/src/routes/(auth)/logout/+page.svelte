<script lang="ts">
	import { goto } from '$app/navigation';
	import { clearAuth } from '$lib/stores/auth.svelte';
	import { api } from '$lib';
	import { resolve } from '$app/paths';

	async function logout() {
		try {
			await api.logout.logout();
		} catch {
			// clear local state regardless of server response
		}
		clearAuth();
		goto(resolve('/login'));
	}

	logout();
</script>
