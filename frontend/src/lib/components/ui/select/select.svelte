<script lang="ts">
	import { Select as SelectPrimitive } from 'bits-ui';
	import type { Snippet } from 'svelte';

	type Props = {
		value?: string;
		onValueChange?: (value: string) => void;
		children?: Snippet;
		type?: 'single' | 'multiple';
	};

	let { value = $bindable(), onValueChange, children, type = 'single' }: Props = $props();

	function handleChange(val: string | string[] | undefined): void {
		if (onValueChange && typeof val === 'string') {
			onValueChange(val);
		}
	}
</script>

<SelectPrimitive.Root
	{...{
		type,
		value,
		onValueChange: handleChange
	} as SelectPrimitive.RootProps}
>
	{@render children?.()}
</SelectPrimitive.Root>
