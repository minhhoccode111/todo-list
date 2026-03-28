<script lang="ts">
	import { Select as SelectPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import ChevronUp from '@lucide/svelte/icons/chevron-up';

	type Props = SelectPrimitive.ScrollUpButtonProps & {
		class?: string;
	};

	let { ref = $bindable(null), class: className, ...restProps }: Props = $props();
</script>

<SelectPrimitive.ScrollUpButton
	bind:ref
	class={cn('flex cursor-default items-center justify-center py-1', className)}
	{...restProps}
>
	<ChevronUp class="size-4" />
</SelectPrimitive.ScrollUpButton>
