<script lang="ts">
	import { goto } from '$app/navigation';
	import { api } from '$lib';

	import { auth } from '$lib/stores/auth.svelte';

	$effect(() => {
		auth = '';
		goto('/login');
	});
</script>
