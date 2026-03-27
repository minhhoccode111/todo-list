<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	$effect(() => {
		if (!auth.token) {
			goto(resolve('/login'));
		}
	});
</script>

<p>Welcome to Todo-list</p>
