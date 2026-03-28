<script lang="ts">
	import { Select as SelectPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import Check from '@lucide/svelte/icons/check';

	type Props = {
		value: string;
		label?: string;
		disabled?: boolean;
		class?: string;
	};

	let { value, label, disabled, class: className }: Props = $props();
</script>

<SelectPrimitive.Item
	{value}
	{disabled}
	class={cn(
		'relative flex w-full cursor-default items-center rounded-sm py-1.5 pr-8 pl-2 text-sm outline-none select-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
		className
	)}
>
	{#snippet children({ selected })}
		<span class="absolute right-2 flex size-3.5 items-center justify-center">
			{#if selected}
				<Check class="size-4" />
			{/if}
		</span>
		{#if label}
			{label}
		{:else}
			{value}
		{/if}
	{/snippet}
</SelectPrimitive.Item>
