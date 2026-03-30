<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
</script>

<Dialog.Root>
	<form>
		<Dialog.Trigger type="button" class={buttonVariants({ variant: 'default' })}>
			New todo
		</Dialog.Trigger>
		<Dialog.Content class="sm:max-w-106.25">
			<Dialog.Header>
				<Dialog.Title>New todo</Dialog.Title>
				<Dialog.Description>
					Create a Todo item with title, description, priority (optional), and due date (optional).
				</Dialog.Description>
			</Dialog.Header>
			<div class="grid gap-4">
				<div class="grid gap-3">
					<Label for="title-1">Title</Label>
					<Input id="title-1" name="title" defaultValue="" />
				</div>
				<div class="grid gap-3">
					<Label for="description-1">Description</Label>
					<Input id="description-1" name="description" defaultValue="" />
				</div>
				<div class="grid gap-3">
					<Label for="priority-1">Priority (optional)</Label>
					<Input id="priority-1" name="priority" defaultValue="" />
				</div>
				<div class="grid gap-3">
					<Label for="due-date-1">Due Date (optional)</Label>
					<Input id="due-date-1" name="due-date" defaultValue="" />
				</div>
			</div>
			<Dialog.Footer>
				<Dialog.Close type="button" class={buttonVariants({ variant: 'outline' })}>
					Cancel
				</Dialog.Close>
				<Button type="submit">Create</Button>
			</Dialog.Footer>
		</Dialog.Content>
	</form>
</Dialog.Root>
