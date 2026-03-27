<script lang="ts">
	import { goto } from '$app/navigation';
	import { clearAuth } from '$lib/stores/auth.svelte';
	import { resolve } from '$app/paths';

	$effect(() => {
		clearAuth();
		goto(resolve('/login'));
	});
</script>
