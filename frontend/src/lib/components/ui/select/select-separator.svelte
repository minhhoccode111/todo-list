<script lang="ts">
	import { cn } from '$lib/utils.js';

	type Props = {
		class?: string;
	};

	let { class: className, ...restProps }: Props = $props();
</script>

<div class={cn('-mx-1 my-1 h-px bg-muted', className)} {...restProps}></div>
