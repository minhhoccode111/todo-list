<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';
	import ChevronLeft from '@lucide/svelte/icons/chevron-left';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import ChevronsLeft from '@lucide/svelte/icons/chevrons-left';
	import ChevronsRight from '@lucide/svelte/icons/chevrons-right';
	import { cn } from '$lib/utils.js';

	interface Props {
		page: number;
		limit: number;
		total: number;
		limitOptions?: number[];
		class?: string;
		onpagechange: (page: number) => void;
		onlimitchange?: (limit: number) => void;
	}

	let {
		page,
		limit,
		total,
		limitOptions = [5, 10, 20, 50],
		class: className,
		onpagechange,
		onlimitchange
	}: Props = $props();

	let totalPages = $derived(Math.ceil(total / limit));
	let canGoPrev = $derived(page > 1);
	let canGoNext = $derived(page < totalPages);

	function getPageNumbers(): (number | '...')[] {
		const pages: (number | '...')[] = [];
		const maxVisible = 5;

		if (totalPages <= maxVisible) {
			for (let i = 1; i <= totalPages; i++) {
				pages.push(i);
			}
		} else {
			if (page <= 3) {
				for (let i = 1; i <= 4; i++) pages.push(i);
				pages.push('...');
				pages.push(totalPages);
			} else if (page >= totalPages - 2) {
				pages.push(1);
				pages.push('...');
				for (let i = totalPages - 3; i <= totalPages; i++) pages.push(i);
			} else {
				pages.push(1);
				pages.push('...');
				for (let i = page - 1; i <= page + 1; i++) pages.push(i);
				pages.push('...');
				pages.push(totalPages);
			}
		}

		return pages;
	}
</script>

<nav class={cn('flex items-center gap-1', className)} aria-label="pagination">
	<Button
		variant="outline"
		size="icon-xs"
		disabled={!canGoPrev}
		onclick={() => onpagechange(1)}
		aria-label="First page"
	>
		<ChevronsLeft class="size-3" />
	</Button>

	<Button
		variant="outline"
		size="icon-xs"
		disabled={!canGoPrev}
		onclick={() => onpagechange(page - 1)}
		aria-label="Previous page"
	>
		<ChevronLeft class="size-3" />
	</Button>

	<div class="flex items-center gap-0.5">
		{#each getPageNumbers() as p (p)}
			{#if p === '...'}
				<span class="px-1 text-muted-foreground">...</span>
			{:else}
				<Button
					variant={p === page ? 'default' : 'ghost'}
					size="icon-xs"
					onclick={() => onpagechange(p)}
					aria-current={p === page ? 'page' : undefined}
				>
					{p}
				</Button>
			{/if}
		{/each}
	</div>

	<Button
		variant="outline"
		size="icon-xs"
		disabled={!canGoNext}
		onclick={() => onpagechange(page + 1)}
		aria-label="Next page"
	>
		<ChevronRight class="size-3" />
	</Button>

	<Button
		variant="outline"
		size="icon-xs"
		disabled={!canGoNext}
		onclick={() => onpagechange(totalPages)}
		aria-label="Last page"
	>
		<ChevronsRight class="size-3" />
	</Button>

	{#if onlimitchange}
		<div class="ml-2 flex items-center gap-2">
			<Select.Root
				type="single"
				value={String(limit)}
				onValueChange={(val) => {
					if (val) {
						onlimitchange?.(Number(val));
					}
				}}
			>
				<Select.Trigger class="h-6 w-16 text-xs" aria-label="Items per page">
					<Select.Value placeholder="10" value={String(limit)} />
				</Select.Trigger>
				<Select.Content>
					<Select.Group>
						{#each limitOptions as opt (opt)}
							<Select.Item value={String(opt)} label={String(opt)} />
						{/each}
					</Select.Group>
				</Select.Content>
			</Select.Root>
			<span class="text-xs text-muted-foreground">/ page</span>
		</div>
	{/if}
</nav>
