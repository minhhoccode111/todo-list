# Todo List

## Learned Concepts

- Apply middlewares at global-level (engine-wide), router-group-level, and route-level
- `sqlc`
- `otter` cache
- We don't need to check for `userID` to exist before using it as foreign key
  to insert `todos`, the database will automatically return error if the
  reference `userID` doesn't exist in `users` table
- refresh token
  - store refresh token in DB to revoke
  - allow user to revoke themselves (logout all devices) with `/logout`
  - refresh endpoint with `/refresh`
  - frontend make one extra request to `/refresh` if receive a `401` response
- rate limit per IP
- unit tests
- sveltekit feels great

## Todo

- [x] register endpoint
- [x] login endpoint
- [x] auth middleware
- [x] create todo
- [x] update todo
- [x] delete todo
- [x] read todos paginate
  - [x] add cache (`~8.00 ms` → `500.00 µs`)
- [x] bearer auth for swagger
- [x] add SPA frontend using sveltekit adapter static
- [x] add refresh token
- [ ] rate limit per IP
- [ ] unit tests

## Preview

<details>
  <summary>Click me</summary>

![login](./docs/img/login.png)
![register](./docs/img/register.png)
![home](./docs/img/home.png)
![new-todo](./docs/img/new-todo.png)
![context-menu](./docs/img/context-menu.png)
![edit](./docs/img/edit.png)
![delete](./docs/img/delete.png)

</details>

## [Requirements](https://roadmap.sh/projects/todo-list-api)
