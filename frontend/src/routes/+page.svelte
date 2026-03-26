<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	$effect(() => {
		if (!auth.token) goto('/login');
	});
</script>
