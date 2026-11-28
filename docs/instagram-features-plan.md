# Instagram-like Feature Plan

## Goal

Turn the existing clean-architecture backend into a fuller Instagram-like API.
The backend already has: auth (JWT), follow, image upload → S3, like, comment, get likes/comments, Kafka notifications, Elasticsearch search sync, outbox, idempotency, retry/DLQ, tracing.

This plan adds the **missing read and write HTTP APIs** and one **housekeeping** item.
Work is done **one checklist item at a time**.

## Scope

**In scope**
- Missing read APIs (data already exists, endpoints don't).
- Missing write APIs (delete/edit/unfollow/unlike), where the data model already supports it.
- A pagination convention for list endpoints.

**Out of scope (do not build)**
- Stories, reels/video, direct messages/chat.
- Saved/bookmarks, user tags, hashtags.
- Profile avatars/bio, private accounts, follow requests.
- Frontend of any kind (this repo is backend-only).

## How to work an item

Follow the repo's existing conventions (see `AGENTS.md` and the `architecture` skill):

1. Read the closest existing feature that is similar and mirror it (e.g. feed mirrors get-image + follow).
2. entity → dto → usecase (+ `_mw_logger`) → outbound repository (+ `_mw_logger`) → DI wire → controller → route.
3. Run `make generate` and `make swag` after interface changes.
4. Add/adjust tests (`make go-test`).
5. Commit with a Conventional Commit message.
6. Tick the checkbox below and move to the next unchecked item.

Keep each item small. Prefer reusing existing repositories/usecases over new ones.

---

## Checklist

Work items one at a time using the `instagram-feature` skill (the step-by-step recipe for adding an endpoint).

### Phase 1 — Read APIs (data already exists)

- [x] `GET /api/feed` — posts from users the current user follows, newest first (paginated). **Highest priority; the core IG feature.**
- [x] `GET /api/images/:id` — single post detail.
- [x] `GET /api/users/:id/images` — a user's post grid (their profile feed), paginated.
- [x] `GET /api/users/:id` — public profile: username, name, follower/following counts.
- [x] `GET /api/users/:id/followers` — list of followers (paginated).
- [x] `GET /api/users/:id/following` — list of followings (paginated).
- [x] `GET /api/explore` — discover posts; wire the existing Elasticsearch image search to an HTTP endpoint.

### Phase 2 — Write APIs

- [ ] `DELETE /api/images/:id` — delete own post (only the owner may delete).
- [ ] `PATCH /api/images/:id` — edit own post caption (only the owner may edit).
- [ ] Unlike — `DELETE /api/images/:id/like` (only like currently exists); decrement like count.
- [ ] `DELETE /api/images/:id/comments/:commentId` — delete a comment (comment owner or image owner).
- [ ] Unfollow — `DELETE /api/users/:id/follow` (only follow exists); decrement follower stats.

### Phase 3 — Housekeeping

- [ ] **Cursor-based pagination convention**, applied consistently to every Phase 1 + 2 list endpoint (feed, user grid, followers, following, explore, likes, comments).

  **Expected outcome**
  - Every list endpoint accepts `limit` (default + max, e.g. 20 / 100) and `cursor` (opaque string; absent = first page).
  - Response envelope is shared across all list endpoints, e.g. `{ "data": [...], "next_cursor": "..." }` where `next_cursor` is empty/null on the last page.
  - Cursor encodes the sort position (e.g. `created_at` + `id` for descending feeds) so new rows between pages don't cause duplicates or skips.
  - Sorting is stable: always tie-break by `id` when ordering by `created_at`.
  - Reusable across domains — one shared helper/dto for encoding/decoding the cursor and building the query, not per-endpoint logic.
  - Documented once (in this repo's conventions) and referenced by each endpoint.

---

## Notes / decisions for the implementing agent

- Reuse `ImageRepository`, `FollowRepository`, `UserRepository`, `LikeRepository`, `CommentRepository` where possible; add methods instead of new repositories.
- The feed query is essentially: images whose `user_id` is in (follows of current user), ordered by `created_at desc`, paginated.
- `/api/explore` should reuse the existing `ImageSearch` (Elasticsearch) outbound already wired for search sync.
- Keep responses consistent with existing DTO shapes (`ImageResponse`, etc.).
- Do not introduce new domains listed under **Out of scope**.