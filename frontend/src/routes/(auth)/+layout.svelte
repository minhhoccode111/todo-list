<script lang="ts">
	let { children } = $props();
</script>

<div class="mx-auto max-w-2xl border border-red-500">
	{@render children()}
</div>
