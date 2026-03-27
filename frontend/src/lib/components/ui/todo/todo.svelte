<script lang="ts">
	import type { EntityTodo } from '$lib/types/api';

	let { todo }: { todo: EntityTodo } = $props();
</script>

<div class="">
	<p>title: {todo.title}</p>
	<p>description: {todo.description}</p>
	<p>priority: {todo.priority}</p>
	<p>completed: {todo.completed}</p>
	<p>due_date: {todo.due_date}</p>
	<p>created_at: {todo.created_at}</p>
	<p>updated_at: {todo.updated_at}</p>
</div>
