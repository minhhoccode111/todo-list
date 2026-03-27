<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import type { EntityTodo, EntityTodos } from '$lib/types/api';
	import { api, getApiError } from '$lib';
	import { toast } from 'svelte-sonner';
	import { Todo } from '$lib/components/ui/todo';

	let page = $state(1);
	let limit = $state(10);
	let total = $state(0);
	let loading = $state(false);

	let todos = $state<EntityTodo[]>([]);

	const getTodos = async () => {
		loading = true;
		try {
			const res = await api.todos.getTodos({ page, limit });
			const data = (await res.json()) as EntityTodos;
			total = data.total;
			todos = data.data;
		} catch (e) {
			const err = await getApiError(e);
			toast.error(err.message || 'Get todos failed');
		} finally {
			loading = false;
		}
	};

	$effect(() => {
		if (!auth.token) {
			goto(resolve('/login'));
			return;
		}

		getTodos();
	});
</script>

<div class="mx-auto max-w-2xl p-8">
	{#if loading}
		<p><center>Loading...</center></p>
	{:else if todos.length == 0}
		<p><center>No todos here yet.</center></p>
	{/if}

	{#each todos as todo (todo.id)}
		<Todo {todo} />
	{/each}
</div>
