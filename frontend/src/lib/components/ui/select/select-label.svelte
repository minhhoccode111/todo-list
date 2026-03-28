<script lang="ts">
	import { cn } from '$lib/utils.js';

	type Props = {
		ref?: HTMLElement | null;
		class?: string;
		children?: import('svelte').Snippet;
	};

	let { ref = $bindable(null), class: className, children, ...restProps }: Props = $props();
</script>

<span bind:this={ref} class={cn('px-2 py-1.5 text-sm font-semibold', className)} {...restProps}>
	{@render children?.()}
</span>
