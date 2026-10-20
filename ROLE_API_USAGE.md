# OneLogin Role API Usage Guide

## Overview

How the OneLogin role endpoints behave, as `onelogin_roles` relies on them. The
membership behaviour below was verified against a development tenant on
2026-10-20.

## Endpoints

- `POST /api/2/roles` - create a role. Accepts `name`, `apps`, `users` and `admins`.
- `GET /api/2/roles/{id}` - returns `id`, `name`, and the `admins`, `apps` and `users` arrays.
- `PUT /api/2/roles/{id}` - update a role (see below).
- `DELETE /api/2/roles/{id}` - delete a role. Members do not need removing first.
- `GET /api/2/roles/{id}/apps`, `/users`, `/admins` - membership, paginated. The
  next page comes from the `After-Cursor` response header, and a cursor cannot be
  combined with `limit` or `page`.

`roleRead` reads membership from the paginated sub-endpoints.

## Updating membership

Each of `apps`, `users` and `admins` in a `PUT` body replaces that membership, and
one left out of the body is not touched:

| `PUT` body                               | Result                                          |
| ---------------------------------------- | ----------------------------------------------- |
| `{"name": "Engineering"}`                | memberships unchanged                           |
| `{"name": "Engineering", "users": []}`   | every user removed; apps and admins unchanged   |
| `{"name": "Engineering", "users": [2]}`  | users become exactly `[2]`, whoever was there   |

`roleUpdate` therefore sends a membership only when the plan changes it, and an
emptied set goes out as `[]`. The three attributes are computed as well as
optional, so one the configuration leaves out never changes and is never sent.

onelogin-go-sdk v4.12.0 and later leave a nil slice out of the body
(onelogin/onelogin-go-sdk#114). Earlier versions sent it as `[]`, so any update
that left a membership out also cleared it.
