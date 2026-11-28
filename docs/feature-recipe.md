# Feature Recipe — adding an HTTP endpoint

Copy-paste reference for implementing endpoints from `docs/instagram-features-plan.md`.
Written from the `/api/feed` implementation; mirror it for every new endpoint.

## 0. Read this first

- Layers flow **inbound → usecase → outbound**; `converter` translates `entity` ↔ `dto`.
- Every interface gets a `_mw_logger.go` decorator (enforced by `make check-rules`).
- No hardcoded table/column strings in repositories — use `pkg/constant/table` and `pkg/constant/column`. **Enforced**: any string literal equal to a known table/column name in `internal/outbound/repository/` fails `make check-rules` (decorators + tests exempt).
- Every handled error must be wrapped with `errkit.AddFuncName(err, "<call site>")` before returning, in impl files of `internal/outbound/*`, `internal/usecase`, `internal/inbound`. **Enforced**; decorators + tests exempt (they record via `telemetry.RecordError` instead). Bare `return err` propagation is exempt.

## 1. File map (exact paths)

| Step | File | What to add |
| :--- | :--- | :--- |
| entity | `internal/entity/<x>_entity.go` | struct + `<X>List []<X>` if listing (skip if it exists) |
| dto | `internal/dto/<x>_dto.go` | `<X>Request`, `<X>Response`, `<X>ResponseList` |
| converter | `internal/converter/<x>_converter.go` | `EntityXToDtoXResponse`, `EntityXListToDtoXResponseList` |
| repo iface+impl | `internal/outbound/repository/<x>_repository.go` | method on interface + impl + `//go:generate moq` already present |
| repo decorator | `internal/outbound/repository/<x>_repository_mw_logger.go` | same method, 4-step body (span, retry, record, log) |
| usecase | `internal/usecase/<x>usecase/<verb>.go` | one file per usecase method |
| usecase iface | `internal/usecase/<x>usecase/<x>_usecase.go` | add method to interface + `Impl` struct already wired |
| usecase decorator | `internal/usecase/<x>usecase/<x>_usecase_mw_logger.go` | same method, no retry |
| controller | `internal/inbound/http/<x>_controller.go` | handler + swagger annotations |
| route | `internal/inbound/http/route/route_http.go` | register on the right group |
| DI | `internal/dependency_injection/usecase_dependency_injection.go` | only if new repo/producer/impl — wire real first, then MwLogger |

## 2. Templates

### Usecase method (`get_feed.go` is the reference)

```go
func (u *ImageUsecaseImpl) GetFeed(ctx context.Context, req dto.FeedRequest) (dto.FeedResponse, error) {
	userAuth := ctxuserauth.Get(ctx)          // current user; nil only on public routes

	err := validatorkit.Validate.Struct(&req) // when the request has validate tags
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.FeedResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetFeed")
	}

	// call repo via u.DB; wrap every returned err with errkit.AddFuncName
}
```

- Tag errors with `errkit.SetCode(err, http.Status<Code>)` for known statuses.
- Always end wrappers with `errkit.AddFuncName(err, "<pkg>.(*<Impl>).<Method>")`.

### Repo method + decorator

Impl — use `column`/`table` constants, never literals:

```go
whereIn, argsIn := column.UserID.InSubquery(subquery)
query := db.WithContext(ctx).
	Where(whereIn, argsIn).
	Order(column.CreatedAt.Desc()).
	Order(column.ID.Desc()).
	Limit(limit)
```

Decorator (4 steps — enforced):

```go
func (r *ImageRepositoryMwLogger) FindFeedByFollowerID(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, followerID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindFeedByFollowerID(ctx, db, imageList, followerID, cursorCreatedAt, cursorID, limit)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{"follower": followerID}
	logkit.LogMw(ctx, fields, err)

	return err
}
```

Retry helper by layer: `DBRetry`/`DBRetryWithData` (repository), `CacheRetry` (cache),
`MessagingRetry` (messaging), `SearchRetry` (search), `StorageRetry` (storage).
Use `*WithData` when the method returns a value. Usecases: no retry.

### Controller handler

```go
// GetFeed godoc
//
//	@Summary		Get feed
//	@Tags			images
//	@Produce		json
//	@Param			limit	query	int		false	"Page size (default 20, max 100)"
//	@Param			cursor	query	string	false	"Opaque cursor; empty for first page"
//	@Security		SimpleApiKeyAuth
//	@Success		200	{object}	response.WebResponse[dto.FeedResponse]
//	@Router			/api/feed [get]
func (c *ImageController) GetFeed(ctx *fiber.Ctx) error {
	span := telemetry.StartController(ctx)
	defer span.End()

	req := dto.FeedRequest{}
	// parse params/path/body; on parse error: errkit.SetCode(..., http.StatusBadRequest)

	res, err := c.Usecase.GetFeed(ctx.UserContext(), req)
	if err != nil {
		logkit.Logger.WithContext(ctx.UserContext()).WithError(err).Error()
		return errkit.AddFuncName(err, "http.(*ImageController).GetFeed")
	}
	return response.Data(ctx, http.StatusOK, res)
}
```

- Current user: `ctxuserauth.Get(ctx.UserContext())` → `userAuth.ID`.
- Path param: `strconv.ParseInt(ctx.Params("imageId"), 10, 64)`.
- Body: `ctx.BodyParser(&req)`.

### Route registration (`route_http.go`)

- Public → `setupGuestRoute`.
- Authenticated → `setupAuthRoute` (has `AuthMiddleware`). Plain path, e.g. `router.Get("/feed", ...)`.
- Repo-native path params, e.g. `images.Get("/:imageId/likes", ...)`.

## 3. Pagination convention (Phase 3 — already built)

Shared helpers live in `pkg/pagination` and `internal/dto/pagination.go`.

- Request: embed `dto.PageRequest` (has `Limit`, `Cursor`). Read `limit`/`cursor` query params.
- `req.Page.NormalizedLimit()` → clamped limit (default 20, max 100).
- `req.Page.DecodeCursor()` → `(pagination.Cursor, hasCursor bool, err error)`; empty cursor = first page, not an error.
- Repo fetches `limit+1` rows + keyset predicate `column.KeysetBefore(cursorCreatedAt, cursorID)`.
- Usecase trims to `limit`, computes `hasMore`, builds envelope via `dto.PageFrom(list, hasMore, lastCreatedAt, lastID)`.
- Envelope: `dto.CursorPage[T]{ Data T, NextCursor string }`; `next_cursor` empty on last page.
- Stable sort: always `ORDER BY created_at DESC, id DESC`.

## 4. Commands (run after interface changes)

```bash
make generate     # regenerate internal/mock (needed after any interface change)
make swag         # regenerate api/ swagger docs (needed after controller changes)
make go-test      # writes logs/go_test.jsonl; grep for "--- FAIL"
make check-rules  # pre-commit rule checks
```

Tests: external package (`package <pkg>_test`), one `<verb>_test.go` per usecase, mocks from
`internal/mock`, fixtures in `helper_test.go`. See `internal/usecase/imageusecase/get_feed_test.go`.

## 5. Commit

Conventional Commits, e.g. `feat(feed): add GET /api/feed endpoint`.
Then tick the checkbox in `docs/instagram-features-plan.md`.