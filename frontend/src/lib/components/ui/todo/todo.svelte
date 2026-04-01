<script lang="ts">
	import type { EntityTodo } from '$lib/types/api';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { cn } from '$lib/utils';

	let { todo }: { todo: EntityTodo } = $props();

	const priorityVariant = $derived(
		todo.priority === 'high' ? 'destructive' : todo.priority === 'med' ? 'default' : 'secondary'
	);

	const priorityLabel = $derived(todo.priority.charAt(0).toUpperCase() + todo.priority.slice(1));

	const formatDate = (dateStr: string): string => {
		return new Date(dateStr).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	};

	$effect(() => {
		if (todo.completed) {
			$inspect(todo);
		}
	});
</script>

<Card.Root class={cn('transition-all hover:shadow-md', todo.completed && 'opacity-60')}>
	<Card.Header class="flex flex-row items-start justify-between space-y-0 pb-2">
		<div class="flex items-center gap-3">
			<div
				class={cn(
					'flex size-5 items-center justify-center rounded-full border-2',
					todo.completed
						? 'border-primary bg-primary text-primary-foreground'
						: 'border-muted-foreground'
				)}
			>
				{#if todo.completed}
					<svg
						xmlns="http://www.w3.org/2000/svg"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="3"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="size-3"
					>
						<polyline points="20 6 9 17 4 12"></polyline>
					</svg>
				{/if}
			</div>
			<Card.Title
				class={cn('text-base font-medium', todo.completed && 'text-muted-foreground line-through')}
			>
				{todo.title}
			</Card.Title>
		</div>
		<Badge variant={priorityVariant}>{priorityLabel}</Badge>
	</Card.Header>
	{#if todo.description}
		<Card.Content class="pt-0">
			<p class="text-sm text-muted-foreground">{todo.description}</p>
		</Card.Content>
	{/if}
	<Card.Footer class="flex items-center justify-between pt-0 text-xs text-muted-foreground">
		{#if todo.due_date}
			<span>Due: {formatDate(todo.due_date)}</span>
			<span>Updated: {formatDate(todo.updated_at)}</span>
		{:else}
			<span>Updated: {formatDate(todo.updated_at)}</span>
		{/if}
		<span>Created: {formatDate(todo.created_at)}</span>
	</Card.Footer>
</Card.Root>
