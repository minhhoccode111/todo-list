# Todo List

## Learned Concepts

- Apply middlewares at global-level (engine-wide), router-group-level, and route-level
- `sqlc`
- `otter` cache
- We don't need to check for `userID` to exist before using it as foreign key
  to insert `todos`, the database will automatically return error if the
  reference `userID` doesn't exist in `users` table
- refresh tokens
  - `/login` and `/register` now set http-only cookies with refresh tokens, alongside access tokens
    - refresh tokens don't need to be JWT, because refresh tokens don't need to extract user_id, expired_at in the claims like JWT, everything can be retrieved from the database
    - refresh tokens will be hashed and stored in database on our server to be able to revoke anytime
    - refresh tokens should be sent in http-only cookies
    - setup frontend to automatically try `/refresh` with `afterResponse` hook (`ky`) when `401` happens
  - `/refresh`: takes refresh token then return new access token and refresh token and invalidate the used one
  - `/logout`: logout current session via http-only cookies
  - `/logout/all`: logout all device sessions
  - `/logout/:id`: logout a device by its session id
  - delivering the refresh-token cookie when frontend and API are on different
    hosts. If the cookie never reaches the server, the `/sessions` handler hashes an
    empty cookie and matches nothing (so `is_current` is always `false`), `/refresh`
    401s (so a reload logs you out), and `/logout/all` 401s behind the dead refresh
    (so the session list never changes). Things to get right:
    - **origin vs site** are different. _Origin_ = scheme + host + **port**; _site_ =
      scheme + registrable domain (eTLD+1), **port does not count**.
      - `localhost:3000` → `localhost:8080` is cross-**origin** but same-**site**.
      - `minhhoccode111.github.io` → a separate API host is cross-**site** (every
        `*.github.io` is its own site via the Public Suffix List).
    - the browser only attaches cookies to a cross-origin request when the frontend
      uses `credentials: 'include'` (`'same-origin'` omits them) **and** the server
      replies with `Access-Control-Allow-Credentials: true` and the reflected
      `Origin` (a wildcard `*` is rejected for credentialed requests).
    - `SameSite` controls cross-**site** sending: `Lax` is fine for same-site (dev),
      but cross-site `fetch`/XHR needs `SameSite=None`. The catch: **browsers reject a
      `SameSite=None` cookie that isn't also `Secure`**, so `None` + `Secure=false`
      (our dev default) means the cookie is silently dropped — which is exactly the bug
      we hit. (And forcing `Secure` on just to keep `None` is the wrong dev fix: outside
      the `localhost` exception, `Secure` cookies aren't accepted over plain HTTP.) So we
      tie the two together: `Secure` on ⇒ `None` (HTTPS prod, cross-site), `Secure` off
      ⇒ `Lax` (HTTP dev, same-site).
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

## Resources

- [Project Requirements](https://roadmap.sh/projects/todo-list-api)
- [Refresh Tokens - What are they and when to use them](https://auth0.com/blog/refresh-tokens-what-are-they-and-when-to-use-them/)
  - refresh token is a credential artifact that lets a client application get new access tokens without having to ask the user to login again
