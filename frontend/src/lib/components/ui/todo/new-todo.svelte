<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import * as Select from '$lib/components/ui/select';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Calendar } from '$lib/components/ui/calendar/index.js';
	import { EntityPriorityLevel } from '$lib/types/api';
	import { getLocalTimeZone, today, type CalendarDate } from '@internationalized/date';
	import CalendarIcon from '@lucide/svelte/icons/calendar';
	import { api, getApiError } from '$lib/api/client';
	import { toast } from 'svelte-sonner';

	const { onSuccess }: { onSuccess: () => void } = $props();

	const priorities = Object.values(EntityPriorityLevel);

	let title = $state('');
	let description = $state('');
	let priority = $state<EntityPriorityLevel>();
	let dueDate = $state<CalendarDate>();

	let errors = $state<{ title?: string; description?: string }>({});
	let loading = $state(false);
	let open = $state(false);

	function formatDate(value: CalendarDate | undefined): string {
		if (!value) return '';
		return value.toDate(getLocalTimeZone()).toLocaleDateString('en-US', {
			day: '2-digit',
			month: 'long',
			year: 'numeric'
		});
	}

	function validate(): boolean {
		errors = {};

		if (!title.trim()) {
			errors.title = 'Title is required';
		} else if (title.length > 255) {
			errors.title = 'Title must be 255 characters or less';
		}

		if (!description.trim()) {
			errors.description = 'Description is required';
		} else if (description.length > 10000) {
			errors.description = 'Description must be 10000 characters or less';
		}

		return Object.keys(errors).length === 0;
	}

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		if (!validate()) return;

		loading = true;
		const dueDateStr = dueDate ? dueDate.toDate(getLocalTimeZone()).toISOString() : undefined;

		try {
			const res = await api.todos.createTodo({
				title,
				description,
				priority,
				due_date: dueDateStr
			});

			console.log(res);

			onSuccess();

			toast.success('Todo created');

			title = '';
			description = '';
			priority = undefined;
			dueDate = undefined;

			open = false;
		} catch (e) {
			const err = await getApiError(e);
			toast.error(err.message || 'Create todo failed');
		} finally {
			loading = false;
		}
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger type="button" class={buttonVariants({ variant: 'default' })}>
		New todo
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-106.25">
		<form onsubmit={handleSubmit} class="space-y-4">
			<Dialog.Header>
				<Dialog.Title>New todo</Dialog.Title>
				<Dialog.Description>
					Create a Todo item with title, description, priority, and due date.
				</Dialog.Description>
			</Dialog.Header>
			<div class="grid gap-4">
				<div class="grid gap-3">
					<Label for="title-1">Title</Label>
					<Input type="text" id="title-1" name="title" bind:value={title} maxlength={255} />
					{#if errors.title}
						<p class="text-sm text-destructive">{errors.title}</p>
					{/if}
				</div>
				<div class="grid gap-3">
					<Label for="description-1">Description</Label>
					<Textarea id="description-1" name="description" bind:value={description} />
					{#if errors.description}
						<p class="text-sm text-destructive">{errors.description}</p>
					{/if}
				</div>
				<div class="flex gap-4">
					<div class="grid flex-1 gap-3">
						<Label for="priority-1">Priority (optional)</Label>
						<Select.Root type="single" bind:value={priority}>
							<Select.Trigger class="w-full justify-between">
								<Select.Value placeholder="Choose priority" value={priority} />
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
						<Label for="due-date-1">Due Date (optional)</Label>
						<Popover.Root>
							<Popover.Trigger id="due-date-1">
								{#snippet child({ props })}
									<Button {...props} variant="outline" class="w-full justify-between font-normal">
										{dueDate ? formatDate(dueDate) : 'Select date'}
										<CalendarIcon class="ms-2 size-4" />
									</Button>
								{/snippet}
							</Popover.Trigger>
							<Popover.Content class="w-auto p-0" align="start">
								<Calendar type="single" bind:value={dueDate} minValue={today(getLocalTimeZone())} />
							</Popover.Content>
						</Popover.Root>
					</div>
				</div>
			</div>
			<Dialog.Footer>
				<Dialog.Close type="button" class={buttonVariants({ variant: 'outline' })}>
					Cancel
				</Dialog.Close>
				<Button type="submit">
					{loading ? 'Creating...' : 'Create'}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
