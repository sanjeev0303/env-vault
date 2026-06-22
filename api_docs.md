# Env Vault API Documentation

Base URL: `/api`

## Authentication (`/api/auth`)
All authentication endpoints (except `/logout`, `/me`, and `/reauth`) are public.

* `POST /register`: Registers a new user (`email`, `password`, `name`).
* `POST /login`: Authenticates user, sets `refresh_token` cookie, returns `access_token` JSON.
* `POST /refresh`: Uses `refresh_token` cookie to issue a new access token.
* `POST /logout`: Revokes current session and clears `refresh_token` cookie.
* `GET /me`: Returns current authenticated user.
* `GET /csrf`: Returns a `csrf_token` cookie and JSON token.
* `POST /forgot-password`: Initiates password reset flow.
* `POST /reset-password`: Completes password reset flow.
* `POST /reauth`: Requires `password` and issues a short-lived `reauth_token` for sensitive actions.

## Projects (`/api/projects`)
Requires `Authorization: Bearer <access_token>`

* `GET /`: Lists all projects the user is a member of.
* `POST /`: Creates a new project.
* `GET /{id}`: Gets project details.
* `DELETE /{id}`: Deletes a project (Requires `Owner` role).

## Environments (`/api/projects/{id}/environments`)
* `GET /`: Lists environments for a project.
* `POST /`: Creates an environment.
* `DELETE /{envID}`: Deletes an environment.

## Secrets (`/api/projects/{id}/environments/{envID}/secrets`)
* `GET /`: Lists secret keys and metadata (cursor-based pagination).
* `POST /`: Creates a new secret (`key`, `value`).
* `PUT /{secretID}`: Updates an existing secret.
* `DELETE /{secretID}`: Deletes a secret.
* `GET /{secretID}/reveal`: Decrypts and reveals a single secret value.
* `GET /{secretID}/history`: Lists version history for a secret.
* `POST /{secretID}/rollback`: Rolls back to a previous version (`version`).
* `GET /export`: Exports all decrypted secrets for an environment (Requires `X-Reauth-Token` header).

## Members (`/api/projects/{id}/members`)
* `GET /`: Lists all project members.
* `POST /`: Adds a member (`email`, `role`).
* `PUT /{userID}`: Updates a member's role.
* `DELETE /{userID}`: Removes a member from the project.

## System
* `GET /api/metrics`: Public metrics endpoint (alloc_bytes, goroutines, etc).
