# Todo List

## Learned Concepts

- Apply Gin Middlewares Global-level (engine-wide), RouterGroup-level, and Route-level
- `sqlc`
- `otter` cache
- We don't need to check for `userID` to exist before using it as foreign key
  to insert `todos`, the database will automatically return error if the
  reference `userID` doesn't exist in `users` table

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
