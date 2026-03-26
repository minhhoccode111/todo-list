<script lang="ts">
	import { goto } from '$app/navigation';
	import { clearAuth } from '$lib/stores/auth.svelte';

	$effect(() => {
		clearAuth();
		goto('/login');
	});
</script>
