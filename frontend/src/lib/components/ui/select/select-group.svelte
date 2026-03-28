<script lang="ts">
	import { Select as SelectPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	type Props = SelectPrimitive.GroupProps & {
		class?: string;
	};

	let { ref = $bindable(null), class: className, children, ...restProps }: Props = $props();
</script>

<SelectPrimitive.Group bind:ref class={cn('p-1', className)} {...restProps}>
	{@render children?.()}
</SelectPrimitive.Group>
