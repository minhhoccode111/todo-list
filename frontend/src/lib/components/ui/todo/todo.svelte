<script lang="ts">
	import type { EntityTodo } from '$lib/types/api';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { cn } from '$lib/utils';
	import * as ContextMenu from '$lib/components/ui/context-menu';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';
	import * as Select from '$lib/components/ui/select';
	import * as Popover from '$lib/components/ui/popover';
	import { Calendar } from '$lib/components/ui/calendar';
	import { EntityPriorityLevel } from '$lib/types/api';
	import { getLocalTimeZone, today, CalendarDate } from '@internationalized/date';
	import CalendarIcon from '@lucide/svelte/icons/calendar';
	import EditIcon from '@lucide/svelte/icons/pencil';
	import TrashIcon from '@lucide/svelte/icons/trash';

	let {
		todo,
		onToggle,
		onDelete,
		onUpdate
	}: {
		todo: EntityTodo;
		onToggle?: (todo: EntityTodo) => void;
		onDelete?: (todo: EntityTodo) => void;
		onUpdate?: (todo: EntityTodo) => void;
	} = $props();

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

	let editOpen = $state(false);
	let deleteOpen = $state(false);

	let editTitle = $state('');
	let editDescription = $state('');
	let editPriority = $state<EntityPriorityLevel>();
	let editDueDate = $state<CalendarDate>();
	let editErrors = $state<{ title?: string; description?: string }>({});
	let editLoading = $state(false);

	$effect(() => {
		editTitle = todo.title;
		editDescription = todo.description;
		editPriority = todo.priority;
		editDueDate = todo.due_date
			? new CalendarDate(
					new Date(todo.due_date).getFullYear(),
					new Date(todo.due_date).getMonth() + 1,
					new Date(todo.due_date).getDate()
				)
			: undefined;
	});

	const priorities = Object.values(EntityPriorityLevel);

	function formatEditDate(value: CalendarDate | undefined): string {
		if (!value) return '';
		return value.toDate(getLocalTimeZone()).toLocaleDateString('en-US', {
			day: '2-digit',
			month: 'long',
			year: 'numeric'
		});
	}

	function openEditDialog() {
		editErrors = {};
		editOpen = true;
	}

	function validateEdit(): boolean {
		editErrors = {};

		if (!editTitle.trim()) {
			editErrors.title = 'Title is required';
		} else if (editTitle.length > 255) {
			editErrors.title = 'Title must be 255 characters or less';
		}

		if (!editDescription.trim()) {
			editErrors.description = 'Description is required';
		} else if (editDescription.length > 10000) {
			editErrors.description = 'Description must be 10000 characters or less';
		}

		return Object.keys(editErrors).length === 0;
	}

	async function handleEditSubmit(e: Event) {
		e.preventDefault();
		if (!validateEdit()) return;

		editLoading = true;
		const dueDateStr = editDueDate
			? editDueDate.toDate(getLocalTimeZone()).toISOString()
			: undefined;

		const updatedTodo: EntityTodo = {
			...todo,
			title: editTitle,
			description: editDescription,
			priority: editPriority || todo.priority,
			due_date: dueDateStr
		};

		onUpdate?.(updatedTodo);
		editOpen = false;
		editLoading = false;
	}

	function handleDeleteConfirm() {
		onDelete?.(todo);
		deleteOpen = false;
	}

	function handleToggle(e: MouseEvent) {
		e.stopPropagation();

		onToggle?.(todo);
	}
</script>

<ContextMenu.Root>
	<ContextMenu.Trigger>
		<Card.Root
			class={cn('transition-all hover:opacity-95 hover:shadow-md', todo.completed && 'opacity-60')}
		>
			<Card.Header class="flex flex-row items-start justify-between space-y-0 pb-2">
				<div class="flex items-center gap-3">
					<button
						type="button"
						onclick={handleToggle}
						class={cn(
							'flex size-5 cursor-pointer items-center justify-center rounded-full border-2',
							todo.completed
								? 'border-primary bg-primary text-primary-foreground hover:border-muted-foreground'
								: 'border-muted-foreground hover:border-primary'
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
					</button>
					<Card.Title
						class={cn(
							'text-base font-medium',
							todo.completed && 'text-muted-foreground line-through'
						)}
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
	</ContextMenu.Trigger>
	<ContextMenu.Content>
		<ContextMenu.Item onclick={openEditDialog}>
			<EditIcon class="size-4" />
			Edit
		</ContextMenu.Item>
		<ContextMenu.Item onclick={() => (deleteOpen = true)} variant="destructive">
			<TrashIcon class="size-4" />
			Delete
		</ContextMenu.Item>
	</ContextMenu.Content>
</ContextMenu.Root>

<Dialog.Root bind:open={editOpen}>
	<Dialog.Content class="sm:max-w-106.25">
		<form onsubmit={handleEditSubmit} class="space-y-4">
			<Dialog.Header>
				<Dialog.Title>Edit todo</Dialog.Title>
				<Dialog.Description>
					Update the todo item with title, description, priority, and due date.
				</Dialog.Description>
			</Dialog.Header>
			<div class="grid gap-4">
				<div class="grid gap-3">
					<Label for="edit-title">Title</Label>
					<Input type="text" id="edit-title" name="title" bind:value={editTitle} maxlength={255} />
					{#if editErrors.title}
						<p class="text-sm text-destructive">{editErrors.title}</p>
					{/if}
				</div>
				<div class="grid gap-3">
					<Label for="edit-description">Description</Label>
					<Textarea id="edit-description" name="description" bind:value={editDescription} />
					{#if editErrors.description}
						<p class="text-sm text-destructive">{editErrors.description}</p>
					{/if}
				</div>
				<div class="flex gap-4">
					<div class="grid flex-1 gap-3">
						<Label for="edit-priority">Priority (optional)</Label>
						<Select.Root type="single" bind:value={editPriority}>
							<Select.Trigger class="w-full justify-between">
								<Select.Value placeholder="Choose priority" value={editPriority} />
							</Select.Trigger>
							<Select.Content>
								<Select.Group>
									<Select.Item value="" label="None" />
									{#each priorities as p (p)}
										<Select.Item value={p} label={p} />
									{/each}
								</Select.Group>
							</Select.Content>
						</Select.Root>
					</div>
					<div class="grid flex-1 gap-3">
						<Label for="edit-due-date">Due Date (optional)</Label>
						<Popover.Root>
							<Popover.Trigger id="edit-due-date">
								{#snippet child({ props })}
									<Button {...props} variant="outline" class="w-full justify-between font-normal">
										{editDueDate ? formatEditDate(editDueDate) : 'Select date'}
										<CalendarIcon class="ms-2 size-4" />
									</Button>
								{/snippet}
							</Popover.Trigger>
							<Popover.Content class="w-auto p-0" align="start">
								<Calendar
									type="single"
									bind:value={editDueDate}
									minValue={today(getLocalTimeZone())}
								/>
							</Popover.Content>
						</Popover.Root>
					</div>
				</div>
			</div>
			<Dialog.Footer>
				<Dialog.Close type="button" class={buttonVariants({ variant: 'outline' })}>
					Cancel
				</Dialog.Close>
				<Button disabled={editLoading} type="submit">
					{editLoading ? 'Saving...' : 'Save'}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<Dialog.Root bind:open={deleteOpen}>
	<Dialog.Content class="sm:max-w-106.25">
		<Dialog.Header>
			<Dialog.Title>Delete Todo</Dialog.Title>
			<Dialog.Description>
				Are you sure you want to delete this todo? This action cannot be undone.
			</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Dialog.Close class={buttonVariants({ variant: 'outline' })}>Cancel</Dialog.Close>
			<Button variant="destructive" onclick={handleDeleteConfirm}>Delete</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
