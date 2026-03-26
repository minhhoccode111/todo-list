<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { api } from '$lib';
	import { goto } from '$app/navigation';
	import { setAuth } from '$lib/stores/auth.svelte';
	import type { ResponseMessage } from '$lib/types/api';
	import { toast } from 'svelte-sonner';
	import { slide } from 'svelte/transition';

	let name = $state('');
	let email = $state('');
	let password = $state('');

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		try {
			const res = await api.auth.register({ name, email, password });
			setAuth(res.token);
			goto('/');
		} catch (e) {
			const err = e as ResponseMessage;
			toast.error(err.message || 'Register failed');
		}
	};
</script>

<div in:slide>
	<Card class="mx-auto w-95">
		<CardHeader>
			<CardTitle>Register</CardTitle>
			<CardDescription>Enter your credentials to sign up.</CardDescription>
		</CardHeader>
		<CardContent>
			<form onsubmit={handleSubmit} class="space-y-4">
				<div class="space-y-2">
					<Label for="name">Name</Label>
					<Input id="name" type="text" bind:value={name} placeholder="John" />
				</div>

				<div class="space-y-2">
					<Label for="email">Email</Label>
					<Input id="email" type="email" bind:value={email} placeholder="you@example.com" />
				</div>

				<div class="space-y-2">
					<Label for="password">Password</Label>
					<Input id="password" type="password" bind:value={password} placeholder="password" />
				</div>

				<Button type="submit" class="w-full">Sign Up</Button>

				<p>Already have an account? <a class="text-sky-500" href="/login">Sign in</a>.</p>
			</form>
		</CardContent>
	</Card>
</div>
