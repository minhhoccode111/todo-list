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
	import { api, getApiError } from '$lib';
	import { goto } from '$app/navigation';
	import { setAuth } from '$lib/stores/auth.svelte';
	import { toast } from 'svelte-sonner';
	import { resolve } from '$app/paths';

	const minlength = 1;
	const maxlength = 255;

	let email = $state('');
	let password = $state('');
	let errors = $state<{ email?: string; password?: string }>({});
	let loading = $state(false);

	const validate = (): boolean => {
		errors = {};
		if (!email.trim()) {
			errors.email = 'Email is required';
		}
		if (!password.trim()) {
			errors.password = 'Password is required';
		}
		return Object.keys(errors).length === 0;
	};

	const handleSubmit = async (e: Event) => {
		e.preventDefault();
		if (!validate()) return;

		loading = true;

		try {
			const res = await api.login.login({ email, password });
			setAuth(res.data.token);
			goto(resolve('/'));
		} catch (e) {
			const err = await getApiError(e);
			toast.error(err.message || 'Login failed');
		} finally {
			loading = false;
		}
	};
</script>

<div>
	<Card class="mx-auto max-w-95">
		<CardHeader>
			<CardTitle>Login</CardTitle>
			<CardDescription>Enter your credentials to sign in.</CardDescription>
		</CardHeader>
		<CardContent>
			<form onsubmit={handleSubmit} class="space-y-4">
				<div class="space-y-2">
					<Label for="email">Email</Label>
					<Input
						id="email"
						type="email"
						bind:value={email}
						placeholder="you@example.com"
						{minlength}
						{maxlength}
					/>
					{#if errors.email}
						<p class="text-sm text-destructive">{errors.email}</p>
					{/if}
				</div>

				<div class="space-y-2">
					<Label for="password">Password</Label>
					<Input
						id="password"
						type="password"
						bind:value={password}
						placeholder="your password"
						{minlength}
						{maxlength}
					/>
					{#if errors.password}
						<p class="text-sm text-destructive">{errors.password}</p>
					{/if}
				</div>

				<Button disabled={loading} type="submit" class="w-full">
					{loading ? 'Signing In...' : 'Sign In'}
				</Button>

				<p>Need an account? <a class="text-sky-500" href={resolve('/register')}>Sign up</a>.</p>
			</form>
		</CardContent>
	</Card>
</div>
