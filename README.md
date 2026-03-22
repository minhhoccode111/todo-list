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
- [ ] create todo
