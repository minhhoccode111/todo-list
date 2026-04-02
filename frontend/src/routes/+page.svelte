<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { api, getApiError } from '$lib/api/client';
	import type { EntityTodo } from '$lib/types/api';
	import { Pagination } from '$lib/components/ui/pagination';
	import { NewTodo, Todo } from '$lib/components/ui/todo';
	import { toast } from 'svelte-sonner';
	import { auth } from '$lib/stores/auth.svelte';
	import { resolve } from '$app/paths';

	let todos = $state<EntityTodo[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let total = $state(0);

	let currentPage = $derived(Number(page.url.searchParams.get('page')) || 1);
	let currentLimit = $derived(Number(page.url.searchParams.get('limit')) || 10);

	async function listTodos(page: number, limit: number) {
		loading = true;
		error = null;

		try {
			const response = await api.todos.getTodos({
				page,
				limit
			});
			todos = response.data.data;
			total = response.data.total;
		} catch (err) {
			error = (await getApiError(err)).message;
			toast.error(error);
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (!auth.token) {
			goto(resolve('/login'));
			return;
		}

		listTodos(currentPage, currentLimit);
	});

	function handlePageChange(newPage: number) {
		const url = new URL(page.url);
		url.searchParams.set('page', String(newPage));
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(`?${url.searchParams.toString()}`);
	}

	function handleLimitChange(newLimit: number) {
		const url = new URL(page.url);
		url.searchParams.set('limit', String(newLimit));
		url.searchParams.set('page', '1');
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(`?${url.searchParams.toString()}`);
	}

	async function handleToggle(todo: EntityTodo) {
		const oldTodo = { ...todo };
		const optimisticState = !todo.completed;

		const index = todos.findIndex((t) => t.id === todo.id);
		if (index !== -1) todos[index] = { ...todos[index], completed: optimisticState };

		try {
			await api.todos.updateTodo(todo.id, { ...todo, completed: optimisticState });
			toast.success(optimisticState ? 'Todo completed!' : 'Todo marked as pending');
		} catch (err) {
			if (index !== -1) todos[index] = oldTodo;
			const errorMessage = (await getApiError(err)).message;
			toast.error(errorMessage || 'Failed to mark todo as completed');
		}
	}

	async function handleUpdate(updatedTodo: EntityTodo) {
		const oldTodo = { ...updatedTodo };
		const index = todos.findIndex((t) => t.id === updatedTodo.id);
		if (index !== -1) todos[index] = updatedTodo;

		try {
			await api.todos.updateTodo(updatedTodo.id, {
				title: updatedTodo.title,
				description: updatedTodo.description,
				priority: updatedTodo.priority,
				due_date: updatedTodo.due_date
			});
			toast.success('Todo updated');
		} catch (err) {
			if (index !== -1) todos[index] = oldTodo;
			const errorMessage = (await getApiError(err)).message;
			toast.error(errorMessage || 'Failed to update todo');
		}
	}

	async function handleDelete(todo: EntityTodo) {
		const index = todos.findIndex((t) => t.id === todo.id);
		const deletedTodo = todos[index];

		if (index !== -1) {
			todos.splice(index, 1);
			total -= 1;
		}

		try {
			await api.todos.deleteTodo(todo.id);
			toast.success('Todo deleted');
		} catch (err) {
			if (index !== -1) {
				todos.splice(index, 0, deletedTodo);
				total += 1;
			}
			const errorMessage = (await getApiError(err)).message;
			toast.error(errorMessage || 'Failed to delete todo');
		}
	}
</script>

<div class="mx-auto flex max-w-2xl flex-col gap-4">
	<div class="flex items-center justify-between gap-4">
		<h2 class="text-2xl font-bold">Todos</h2>

		<NewTodo onSuccess={() => listTodos(currentPage, currentLimit)} />
	</div>

	{#if loading}
		<div class="p-4">
			<p><center>Loading...</center></p>
		</div>
	{:else if error && todos.length === 0}
		<div class="p-4">
			<p class="text-destructive"><center>{error}</center></p>
		</div>
	{:else}
		<div class="flex flex-col gap-2">
			{#each todos as todo (todo.id)}
				<Todo {todo} onToggle={handleToggle} onUpdate={handleUpdate} onDelete={handleDelete} />
			{/each}

			{#if todos.length === 0}
				<p class="text-muted-foreground">No todos found.</p>
			{/if}
		</div>

		{#if total > 0}
			<div class="mt-4 flex justify-center">
				<Pagination
					page={currentPage}
					limit={currentLimit}
					{total}
					onpagechange={handlePageChange}
					onlimitchange={handleLimitChange}
				/>
			</div>

			<p class="text-center text-sm text-muted-foreground">
				Showing {(currentPage - 1) * currentLimit + 1} to {Math.min(
					currentPage * currentLimit,
					total
				)} of {total} todos
			</p>
		{/if}
	{/if}
</div>
