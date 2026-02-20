# Kanban Backend API Implementation

## TL;DR

> **Quick Summary**: Build production-grade kanban backend with 15 REST controllers, real-time WebSocket updates, JWT authentication, and cron jobs using Fiber v2 + GORM v2 + PostgreSQL.
> 
> **Deliverables**:
> - 15 REST controllers with full CRUD operations
> - GORM models with relationships (User, Board, Column, Task, Comment, Label, Attachment)
> - JWT authentication (1h access token + refresh token rotation)
> - WebSocket real-time updates (board changes, task movements, comments)
> - 3 cron jobs (due reminders, soft-delete cleanup, notification cleanup)
> - TDD test suite with ≥70% coverage
> 
> **Estimated Effort**: XL (50+ hours)
> **Parallel Execution**: YES - 9 waves
> **Critical Path**: Utils → Models → Repositories → Services → Controllers → Routes → Tests → Final QA

---

## Context

### Original Request
Create controllers, database models with GORM, and routes for a comprehensive kanban board backend system with 15 controllers including Auth, User, Board, Column, Task, Comments, Labels, Attachments, Notifications, Widgets, Search, Activity Logs, Cron Jobs, and WebSockets.

### Interview Summary
**Key Discussions**:
- Architecture: Full 3-layer (Controller → Service → Repository → Models) confirmed
- Testing: TDD approach with Go testing + testify confirmed
- JWT: 1 hour expiry with refresh tokens confirmed
- WebSocket: JWT in subprotocol (Authorization header with Bearer token)
- Migrations: Versioned GORM migration files (not AutoMigrate)
- Router: Fiber v2 (consistent with existing setup)

**Research Findings**:
- Focalboard (25k+ stars) - Kanban board with WebSocket patterns
- learning-cloud-native-go/myapp (1000+ stars) - GORM + clean architecture
- filebrowser (20k+ stars) - Production JWT auth patterns
- Best practices: Repository pattern, context-based operations, soft delete, audit hooks

### Metis Review
**Identified Gaps** (addressed):
- Auth method: Email/password only (no OAuth)
- Password policy: Min 8 chars
- Permissions: Simple ownership (owner + members only)
- File constraints: 10MB max, all types, unlimited attachments
- WebSocket scope: All changes broadcast to board members
- Search scope: Task title + description only
- Cron jobs: Due reminders + soft-delete cleanup

**Guardrails Applied**:
- NO: Admin dashboard, OAuth, email notifications, external search, RBAC, public boards, virus scanning, webhooks
- NO: "User manually tests..." - all criteria executable via commands
- NO: Scope creep - strict MVP focus

---

## Work Objectives

### Core Objective
Build production-grade kanban backend REST API with real-time collaboration features, following clean architecture principles with comprehensive test coverage.

### Concrete Deliverables
- 15 Fiber controllers with full REST endpoints
- 7 GORM models with relationships and soft delete
- JWT authentication system with refresh tokens
- WebSocket server for real-time updates
- 3 cron jobs (reminders, cleanup)
- Test suite with ≥70% coverage
- GORM migration files

### Definition of Done
- [ ] All 15 controllers implemented with tests passing
- [ ] JWT auth working (register, login, refresh, logout)
- [ ] WebSocket broadcasting real-time changes
- [ ] All endpoints return standardized JSON responses
- [ ] Test coverage ≥70%
- [ ] Cron jobs registered and logging

### Must Have
- Email/password authentication only (no OAuth)
- Simple ownership permissions (no RBAC)
- 1-hour JWT access tokens with refresh token rotation
- WebSocket real-time for all board changes
- GORM versioned migrations (not AutoMigrate)
- TDD workflow (tests before implementation)

### Must NOT Have (Guardrails)
- ❌ Admin dashboard or backend UI
- ❌ OAuth/Social login (Google, GitHub, etc.)
- ❌ Email notifications system (WebSocket only)
- ❌ External search engines (SQL filtering only)
- ❌ RBAC roles (Admin/Editor/Viewer)
- ❌ Public boards or external sharing
- ❌ File virus scanning
- ❌ Webhook integrations (Slack, Jira, etc.)
- ❌ AI-slop: Excessive comments, over-abstraction, premature extraction

---

## Verification Strategy

> **ZERO HUMAN INTERVENTION** — ALL verification is agent-executed. No exceptions.
> Acceptance criteria requiring "user manually tests/confirms" are FORBIDDEN.

### Test Decision
- **Infrastructure exists**: NO (starting from scratch)
- **Automated tests**: YES (TDD)
- **Framework**: Go testing + testify + testify/mock
- **If TDD**: Each task follows RED (failing test) → GREEN (minimal impl) → REFACTOR
- **Coverage Target**: ≥70%
- **Database Tests**: sqlite in-memory for unit tests

### QA Policy
Every task MUST include agent-executed QA scenarios (see TODO template below).
Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

- **Frontend/UI**: N/A (API only)
- **TUI/CLI**: N/A (API only)
- **API/Backend**: Use Bash (curl) — Send requests, assert status + response fields
- **Library/Module**: Use Bash (bun/node REPL) — Import, call functions, compare output

---

## Execution Strategy

### Parallel Execution Waves

> Maximize throughput by grouping independent tasks into parallel waves.
> Each wave completes before the next begins.
> Target: 5-8 tasks per wave. Fewer than 3 per wave (except final) = under-splitting.

```
Wave 1 (Foundation & Utils) - Start Immediately:
├── Task 1: JWT utilities implementation [quick]
├── Task 2: Response helper utilities [quick]
├── Task 3: Password hashing utilities [quick]
├── Task 4: Validation utilities [quick]
├── Task 5: Error types & middleware base [quick]
├── Task 6: Logger middleware [quick]
└── Task 7: CORS middleware [quick]

Wave 2 (GORM Models) - After Wave 1:
├── Task 8: User model + audit hooks [unspecified-high]
├── Task 9: Board model + relationships [unspecified-high]
├── Task 10: Column model + relationships [unspecified-high]
├── Task 11: Task model + relationships [unspecified-high]
├── Task 12: Comment model [unspecified-high]
├── Task 13: Label model [unspecified-high]
├── Task 14: Attachment model [unspecified-high]
└── Task 15: Notification model [unspecified-high]

Wave 3 (Migrations) - After Wave 2:
├── Task 16: Migration setup & structure [deep]
├── Task 17: User migration [quick]
├── Task 18: Board/Column/Task migrations [quick]
├── Task 19: Comment/Label/Attachment migrations [quick]
└── Task 20: Notification migration [quick]

Wave 4 (Repositories - Data Access) - After Wave 3:
├── Task 21: User repository [quick]
├── Task 22: Board repository [quick]
├── Task 23: Column repository [quick]
├── Task 24: Task repository [quick]
├── Task 25: Comment repository [quick]
├── Task 26: Label repository [quick]
└── Task 27: Attachment repository [quick]

Wave 5 (Services - Business Logic) - After Wave 4:
├── Task 28: Auth service [deep]
├── Task 29: User service [unspecified-high]
├── Task 30: Board service [unspecified-high]
├── Task 31: Column service [quick]
├── Task 32: Task service [deep]
├── Task 33: Comment service [quick]
├── Task 34: Label service [quick]
└── Task 35: Attachment service [quick]

Wave 6 (Auth & User Controllers) - After Wave 5:
├── Task 36: Auth controller [deep]
├── Task 37: User controller [unspecified-high]
└── Task 38: Auth middleware [deep]

Wave 7 (Board/Column/Task Controllers) - After Wave 6:
├── Task 39: Board controller [unspecified-high]
├── Task 40: Column controller [unspecified-high]
├── Task 41: Task controller [deep]
├── Task 42: InProgressRule controller [quick]
└── Task 43: ActivityLog controller [quick]

Wave 8 (Extended Controllers) - After Wave 7:
├── Task 44: Comment controller [quick]
├── Task 45: Label controller [quick]
├── Task 46: Attachment controller [quick]
├── Task 47: Notification controller [quick]
├── Task 48: Widget controller [quick]
├── Task 49: Search controller [unspecified-high]
└── Task 50: WebSocket controller [deep]

Wave 9 (Routes + Cron + Integration) - After Wave 8:
├── Task 51: Routes setup [unspecified-high]
├── Task 52: Cron jobs setup [deep]
├── Task 53: Update main.go [quick]
└── Task 54: Health check endpoints [quick]

Wave FINAL (Verification) - After ALL tasks:
├── Task F1: Run all tests & verify coverage [deep]
├── Task F2: Integration test suite [deep]
├── Task F3: Final verification wave [oracle]
└── Task F4: Code quality review [unspecified-high]

Critical Path: T1-T7 → T8-T15 → T16-T20 → T21-T27 → T28-T35 → T36-T38 → T39-T43 → T44-T50 → T51-T54 → F1-F4
Parallel Speedup: ~65% faster than sequential
Max Concurrent: 7 (Wave 2, Wave 5, Wave 8)
```

### Dependency Matrix (FULL)

- **1-7**: — — 8-15, 21-27, 36, 51, 52
- **8-15**: 1 — 16, 21-27
- **16-20**: 8-15 — 21-27
- **21-27**: 16-20 — 28-35
- **28-35**: 21-27 — 36-37, 39-50
- **36**: 28, 35 — 38, 39-50, 51
- **37**: 28 — 39-50
- **38**: 36, 37 — 39-50
- **39**: 28, 38 — 42-43, 51
- **40**: 28 — 41, 42, 51
- **41**: 28 — 42, 51
- **42**: 28, 41 — 43, 51
- **43**: 28, 39-42 — 51
- **44-50**: 28, 38 — 51
- **51**: 36-50 — 52, 53
- **52**: 51 — 53, 54
- **53**: 36-50, 52 — 54
- **54**: 36-50, 53 — F1-F4

### Agent Dispatch Summary

- **1**: **7** — T1-T7 → `quick`
- **2**: **8** — T8-T15 → `unspecified-high`
- **3**: **5** — T16-T20 → `deep` (T16), `quick` (T17-T20)
- **4**: **7** — T21-T27 → `quick`
- **5**: **8** — T28-T35 → `deep` (T28, T32), `unspecified-high` (T29, T30, T31), `quick` (T33-T35)
- **6**: **3** — T36-T38 → `deep` (T36, T38), `unspecified-high` (T37)
- **7**: **5** — T39-T43 → `unspecified-high` (T39, T41), `quick` (T40, T42, T43)
- **8**: **7** — T44-T50 → `deep` (T50), `unspecified-high` (T49), `quick` (T44-T48)
- **9**: **4** — T51-T54 → `unspecified-high` (T51), `deep` (T52), `quick` (T53-T54)
- **FINAL**: **4** — F1-F4 → `deep` (F1-F3), `unspecified-high` (F2, F4)

---

- [x] 1. **JWT Utilities Implementation**

  **What to do**:
  - Implement `GenerateToken(userID string, expiry time.Duration) (string, error)` in `utils/jwt.go`
  - Implement `ValidateToken(tokenString string) (*jwt.Claims, error)` 
  - Implement `GenerateRefreshToken(userID string) (string, error)`
  - Implement `ValidateRefreshToken(tokenString string) (string, error)`
  - Use `github.com/golang-jwt/jwt/v5` package
  - Set signing key from `os.Getenv("JWT_SECRET")`
  - Access token expiry: 1 hour, Refresh token expiry: 7 days

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Straightforward JWT implementation with well-established library
  - **Skills**: []
    - No special skills needed - standard Go JWT library

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 2-7) | Sequential
  - **Blocks**: Tasks 8-50 (all depend on JWT utilities)
  - **Blocked By**: None (can start immediately)

  **References**:
  - `utils/jwt.go:1` - Empty file to be implemented
  - Official docs: `https://pkg.go.dev/github.com/golang-jwt/jwt/v5` - JWT v5 API reference
  - Research: filebrowser (20k+ stars) - Production JWT patterns from research

  **Acceptance Criteria**:
  - [ ] Test file created: `utils/jwt_test.go`
  - [ ] `go test ./utils -v` → PASS (GenerateToken, ValidateToken, RefreshToken tests)
  - [ ] Generated tokens include correct claims (user_id, exp, iat)
  - [ ] Expired tokens return error
  - [ ] Invalid tokens return error

  **QA Scenarios**:

  ```
  Scenario: Generate and validate access token
    Tool: Bash (go test)
    Preconditions: None
    Steps:
      1. Run: go test ./utils -run TestGenerateToken -v
      2. Verify: test passes, returns JWT string with claims
      3. Run: go test ./utils -run TestValidateToken -v
      4. Verify: valid token is accepted, extracts user_id
    Expected Result: Token generation works, validation extracts correct user_id
    Failure Indicators: Test fails, token missing claims, validation returns wrong user_id
    Evidence: .sisyphus/evidence/task-1-generate-validate-jwt.txt

  Scenario: Token expiry validation
    Tool: Bash (go test)
    Preconditions: Generate token with 1 hour expiry
    Steps:
      1. Wait 61 minutes (simulate expiry in test via time mocking)
      2. Run: go test ./utils -run TestExpiredToken -v
      3. Verify: expired token returns error
    Expected Result: Expired tokens are rejected with clear error message
    Failure Indicators: Expired tokens still accepted, unclear error
    Evidence: .sisyphus/evidence/task-1-expired-jwt.txt

  Scenario: Refresh token generation and validation
    Tool: Bash (go test)
    Preconditions: None
    Steps:
      1. Run: go test ./utils -run TestRefreshToken -v
      2. Verify: refresh token with 7 day expiry is generated
      3. Run: go test ./utils -run TestValidateRefreshToken -v
      4. Verify: valid refresh token returns user_id
    Expected Result: Refresh tokens work with longer expiry
    Failure Indicators: Refresh token same as access token, validation fails
    Evidence: .sisyphus/evidence/task-1-refresh-token.txt
  ```

  **Commit**: NO (group with Tasks 2-7)

---

## TODOs

> Implementation + Test = ONE Task. Never separate.

- [ ] 8. **User Model + Audit Hooks**

  **What to do**:
  - Create `models/user.go` with User struct
  - Fields: ID (uuid), Username, Email, Password (hashed), CreatedAt, UpdatedAt, DeletedAt
  - GORM tags: `gorm:"primarykey"`, `gorm:"uniqueIndex"`, `gorm:"not null"`
  - Implement audit hooks: BeforeCreate, AfterUpdate (log to audit_log table)
  - Add refresh tokens relation: `RefreshTokens []RefreshToken`

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: GORM model with relationships and hooks - requires GORM knowledge
  - **Skills**: []
    - No special skills needed - standard GORM patterns

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 9-15) | Sequential
  - **Blocks**: Task 16 (migrations depend on models)
  - **Blocked By**: Task 1 (JWT utilities needed for auth)

  **References**:
  - `models/` - Empty directory, create `models/user.go`
  - GORM docs: `https://gorm.io/docs/models` - GORM model patterns
  - Research: learning-cloud-native-go/myapp (1000+ stars) - GORM model examples

  **Acceptance Criteria**:
  - [ ] Test file created: `models/user_test.go`
  - [ ] `go test ./models -v` → PASS (User model tests)
  - [ ] GORM tags defined correctly
  - [ ] Audit hooks trigger on create/update

  **QA Scenarios**:

  ```
  Scenario: User model GORM tags
    Tool: Bash (go test)
    Preconditions: None
    Steps:
      1. Run: go test ./models -run TestUserModel -v
      2. Verify: User struct has correct GORM tags
      3. Verify: Email has unique index
    Expected Result: User model defined correctly
    Failure Indicators: Missing GORM tags, no unique index on email
    Evidence: .sisyphus/evidence/task-8-user-model.txt

  Scenario: Audit hooks trigger
    Tool: Bash (go test)
    Preconditions: User created
    Steps:
      1. Run: go test ./models -run TestAuditHooks -v
      2. Create user via repository
      3. Verify: Audit log entry created
    Expected Result: Audit hooks work correctly
    Failure Indicators: No audit log entries
    Evidence: .sisyphus/evidence/task-8-audit-hooks.txt
  ```

  **Commit**: NO (group with Tasks 9-15)

- [ ] 9-15. **Board, Column, Task, Comment, Label, Attachment, Notification Models**
  *(7 models in one task - similar complexity)*

  **What to do**:
  - `models/board.go`: Board struct with ID, Title, UserID, Color, CreatedAt, UpdatedAt, DeletedAt
  - `models/column.go`: Column struct with ID, BoardID, Title, Order, Tasks relation
  - `models/task.go`: Task struct with ID, ColumnID, Title, Description, Deadline, CreatedAt, UpdatedAt, DeletedAt
  - `models/comment.go`: Comment struct with ID, TaskID, UserID, Content, CreatedAt
  - `models/label.go`: Label struct with ID, Name, Color, Tasks many-to-many
  - `models/attachment.go`: Attachment struct with ID, TaskID, FileName, FileURL, CreatedAt
  - `models/notification.go`: Notification struct with ID, UserID, Message, ReadAt, CreatedAt
  - Relationships: Board→Columns→Tasks (one-to-many), Tasks↔Labels (many-to-many)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Multiple GORM models with relationships
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Task 8) | Sequential
  - **Blocks**: Task 16 (migrations depend on models)
  - **Blocked By**: Task 1 (JWT utilities needed for user_id in models)

  **References**:
  - `models/` - Create files for each model
  - GORM docs: `https://gorm.io/docs/associations` - Relationship patterns

  **Acceptance Criteria**:
  - [ ] Test files created for each model
  - [ ] `go test ./models -v` → PASS (all model tests)
  - [ ] Foreign keys defined correctly
  - [ ] Many-to-many relationships work

  **QA Scenarios**:

  ```
  Scenario: Model relationships work
    Tool: Bash (go test)
    Preconditions: Models created
    Steps:
      1. Run: go test ./models -run TestRelationships -v
      2. Create board with columns and tasks
      3. Verify: Cascading delete works
    Expected Result: Relationships defined correctly
    Failure Indicators: Foreign key errors, no cascading
    Evidence: .sisyphus/evidence/task-9-15-relationships.txt
  ```

  **Commit**: NO (group with Tasks 8)

- [ ] 16. **Migration Setup & Structure**

  **What to do**:
  - Create `migrations/` directory structure
  - Create migration registry to track applied migrations
  - Implement up() and down() functions per migration
  - Use GORM migrations (not AutoMigrate)
  - Store migration version in `migrations` table

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Requires understanding GORM migration patterns and versioning
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: Tasks 17-20 (specific migrations depend on setup)
  - **Blocked By**: Tasks 8-15 (models must exist first)

  **References**:
  - `migrations/` - Empty directory, create structure
  - GORM docs: `https://gorm.io/docs/migration` - Migration patterns

  **Acceptance Criteria**:
  - [ ] Test file created: `migrations/migration_test.go`
  - [ ] `go test ./migrations -v` → PASS (migration tests)
  - [ ] Migration registry tracks applied versions
  - [ ] up() and down() functions work correctly

  **QA Scenarios**:

  ```
  Scenario: Migration applies and rolls back
    Tool: Bash (go test)
    Preconditions: None
    Steps:
      1. Run: go test ./migrations -run TestMigration -v
      2. Apply migration, verify table created
      3. Rollback migration, verify table dropped
    Expected Result: Migrations work correctly
    Failure Indicators: Migration fails, rollback doesn't work
    Evidence: .sisyphus/evidence/task-16-migration.txt
  ```

  **Commit**: NO (group with Tasks 17-20)

- [ ] 17-20. **User, Board/Column/Task/Comment/Label/Attachment/Notification Migrations**
  *(6 specific migrations - similar complexity)*

  **What to do**:
  - Create `migrations/00001_create_users.up.sql` and `.down.sql`
  - Create `migrations/00002_create_boards.up.sql` and `.down.sql`
  - Create `migrations/00003_create_columns.up.sql` and `.down.sql`
  - Create `migrations/00004_create_tasks.up.sql` and `.down.sql`
  - Create `migrations/00005_create_comments.up.sql` and `.down.sql`
  - Create `migrations/00006_create_labels.up.sql` and `.down.sql`
  - Create `migrations/00007_create_attachments.up.sql` and `.down.sql`
  - Create `migrations/00008_create_notifications.up.sql` and `.down.sql`
  - Include indexes, foreign keys, constraints

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Straightforward table creation SQL
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with Tasks 17-20) | Sequential
  - **Blocks**: Tasks 21-27 (repositories depend on migrations)
  - **Blocked By**: Task 16 (migration setup)

  **References**:
  - `migrations/` - Create migration files

  **Acceptance Criteria**:
  - [ ] `go run migrations/migrate up` applies all migrations
  - [ ] `go run migrations/migrate down` rolls back migrations
  - [ ] Foreign keys work (cascade delete)
  - [ ] Indexes created

  **QA Scenarios**:

  ```
  Scenario: Migrations apply successfully
    Tool: Bash (go test)
    Preconditions: Migration setup complete
    Steps:
      1. Run: go test ./migrations -run TestApplyMigrations -v
      2. Verify: All tables created with correct schema
      3. Verify: Indexes and foreign keys exist
    Expected Result: Migrations apply correctly
    Failure Indicators: Migration fails, wrong schema
    Evidence: .sisyphus/evidence/task-17-20-apply-migrations.txt

  Scenario: Migration rollback works
    Tool: Bash (go test)
    Preconditions: Migrations applied
    Steps:
      1. Run: go test ./migrations -run TestRollback -v
      2. Verify: Tables dropped correctly
    Expected Result: Rollback works
    Failure Indicators: Tables not dropped, errors
    Evidence: .sisyphus/evidence/task-17-20-rollback.txt
  ```

  **Commit**: YES (commit Tasks 16-20 as "Migrations: Setup + user/board/column/task/comment/label/attachment/notification")

- [ ] 21. **User Repository**

  **What to do**:
  - Create `repositories/user_repository.go`
  - Implement `UserRepository` interface: Create, FindByEmail, FindByID, Update, SoftDelete
  - Use GORM with config.DB
  - Return *models.User or errors

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Standard CRUD repository pattern
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with Tasks 22-27) | Sequential
  - **Blocks**: Task 28 (auth service depends on user repository)
  - **Blocked By**: Task 20 (migrations must be complete)

  **References**:
  - `repositories/` - Empty directory, create file
  - Research: learning-cloud-native-go/myapp (1000+ stars) - Repository pattern

  **Acceptance Criteria**:
  - [ ] Test file created: `repositories/user_repository_test.go`
  - [ ] `go test ./repositories -v` → PASS (all CRUD operations)
  - [ ] Context used for timeout/cancellation

  **QA Scenarios**:

  ```
  Scenario: User repository CRUD operations
    Tool: Bash (go test)
    Preconditions: Migrations applied
    Steps:
      1. Run: go test ./repositories -run TestUserCRUD -v
      2. Create user, find by email, update, soft delete
      3. Verify: All operations work correctly
    Expected Result: User repository works
    Failure Indicators: CRUD fails, context not used
    Evidence: .sisyphus/evidence/task-21-user-repo.txt
  ```

  **Commit**: NO (group with Tasks 22-27)

- [ ] 22-27. **Board, Column, Task, Comment, Label, Attachment, Notification Repositories**
  *(6 repositories - similar complexity)*

  **What to do**:
  - `repositories/board_repository.go`: Create, FindByUser, FindByID, Update, Delete, AddMember, RemoveMember
  - `repositories/column_repository.go`: Create, FindByBoard, Update, Reorder, Delete
  - `repositories/task_repository.go`: Create, FindByColumn, FindByID, Update, Move, Delete, Assign, Unassign
  - `repositories/comment_repository.go`: Create, FindByTask, FindByID, Update, Delete
  - `repositories/label_repository.go`: Create, FindByBoard, FindByID, Update, Delete, AttachToTask, DetachFromTask
  - `repositories/attachment_repository.go`: Create, FindByTask, FindByID, Delete

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Standard CRUD repository patterns
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 4 (with Task 21) | Sequential
  - **Blocks**: Task 28-35 (services depend on repositories)
  - **Blocked By**: Task 20 (migrations)

  **References**:
  - `repositories/` - Create files for each repository

  **Acceptance Criteria**:
  - [ ] Test files created for each repository
  - [ ] `go test ./repositories -v` → PASS (all repository tests)
  - [ ] Foreign key queries work correctly

  **QA Scenarios**:

  ```
  Scenario: Board repository operations
    Tool: Bash (go test)
    Preconditions: Migrations applied
    Steps:
      1. Run: go test ./repositories -run TestBoardRepo -v
      2. Create board, add member, find by user
      3. Verify: All operations work
    Expected Result: Board repository works
    Failure Indicators: Operations fail
    Evidence: .sisyphus/evidence/task-22-27-board-repo.txt

  Scenario: Task repository move operation
    Tool: Bash (go test)
    Preconditions: Tasks and columns exist
    Steps:
      1. Run: go test ./repositories -run TestTaskMove -v
      2. Create task in column A, move to column B
      3. Verify: ColumnID updated correctly
    Expected Result: Task move works
    Failure Indicators: Move fails, wrong column
    Evidence: .sisyphus/evidence/task-22-27-task-move.txt
  ```

  **Commit**: YES (commit Tasks 21-27 as "Repositories: User + board + column + task + comment + label + attachment")

- [ ] 28. **Auth Service**

  **What to do**:
  - Create `services/auth_service.go`
  - Implement `Register(email, password) (user, token, error)`
  - Implement `Login(email, password) (accessToken, refreshToken, error)`
  - Implement `RefreshToken(refreshToken) (accessToken, error)`
  - Implement `Logout(userID) error`
  - Use `utils/jwt.go` for token generation/validation
  - Use `utils/password.go` for hashing/checking
  - Store refresh token in database

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Business logic for authentication with JWT and refresh tokens
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: Task 36 (auth controller depends on auth service)
  - **Blocked By**: Tasks 1, 21 (JWT utils and user repository)

  **References**:
  - `services/` - Create file
  - `services/S3_service.go` - Pattern for service structure

  **Acceptance Criteria**:
  - [ ] Test file created: `services/auth_service_test.go`
  - [ ] `go test ./services -v` → PASS (register, login, refresh, logout tests)
  - [ ] Password hashed before storing
  - [ ] Refresh tokens stored in database

  **QA Scenarios**:

  ```
  Scenario: Register new user
    Tool: Bash (curl)
    Preconditions: App running
    Steps:
      1. Run: curl -X POST http://localhost:8080/api/auth/register -d '{"email":"test@example.com","password":"password123"}' -H "Content-Type: application/json"
      2. Verify: Returns 201 with access token
      3. Verify: User stored in database with hashed password
    Expected Result: Registration works
    Failure Indicators: Wrong status, token not returned, password not hashed
    Evidence: .sisyphus/evidence/task-28-register.txt

  Scenario: Login and get refresh token
    Tool: Bash (curl)
    Preconditions: User registered
    Steps:
      1. Run: curl -X POST http://localhost:8080/api/auth/login -d '{"email":"test@example.com","password":"password123"}' -H "Content-Type: application/json"
      2. Verify: Returns access token + refresh token
      3. Verify: Refresh token stored in DB
    Expected Result: Login works
    Failure Indicators: Wrong status, missing tokens
    Evidence: .sisyphus/evidence/task-28-login.txt

  Scenario: Refresh token rotation
    Tool: Bash (curl)
    Preconditions: Valid refresh token
    Steps:
      1. Run: curl -X POST http://localhost:8080/api/auth/refresh -d '{"refresh_token":"xxx"}' -H "Content-Type: application/json"
      2. Verify: Returns new access token
    Expected Result: Refresh works
    Failure Indicators: Invalid tokens accepted, no new token
    Evidence: .sisyphus/evidence/task-28-refresh.txt
  ```

  **Commit**: NO (group with Tasks 29-35)

- [ ] 29-35. **User, Board, Column, Task, Comment, Label, Attachment Services**
  *(6 services - business logic layer)*

  **What to do**:
  - `services/user_service.go`: GetProfile, UpdateProfile, DeleteAccount
  - `services/board_service.go`: Create, Get, Update, Delete, AddMember, RemoveMember, ListForUser
  - `services/column_service.go`: Create, Get, Update, Reorder, Delete
  - `services/task_service.go`: Create, Get, Update, Delete, Move, Reorder, Assign, Unassign, Complete
  - `services/comment_service.go`: Create, Get, Update, Delete
  - `services/label_service.go`: Create, Get, Update, Delete
  - `services/attachment_service.go`: Create, Get, Delete
  - All services use repositories, not direct DB access
  - All services return custom errors

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high` (T29, T30), `quick` (T31, T33-T35)
    - Reason: Business logic services vary in complexity
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with Task 28) | Sequential
  - **Blocks**: Tasks 36-37 (auth/user controllers) and 39-50 (extended controllers)
  - **Blocked By**: Tasks 21-27 (repositories)

  **References**:
  - `services/` - Create service files
  - `services/S3_service.go` - Pattern for service structure

  **Acceptance Criteria**:
  - [ ] Test files created for each service
  - [ ] `go test ./services -v` → PASS (all service tests)
  - [ ] Services use repositories, not DB directly
  - [ ] Error handling works correctly

  **QA Scenarios**:

  ```
  Scenario: Board service business logic
    Tool: Bash (go test)
    Preconditions: Repositories exist
    Steps:
      1. Run: go test ./services -run TestBoardService -v
      2. Create board with initial columns
      3. Verify: Business logic applied (e.g., validation)
    Expected Result: Board service works
    Failure Indicators: Service fails, no business logic
    Evidence: .sisyphus/evidence/task-29-35-board-service.txt

  Scenario: Task service move validation
    Tool: Bash (go test)
    Preconditions: Tasks and columns exist
    Steps:
      1. Run: go test ./services -run TestTaskMoveValidation -v
      2. Try move task to invalid column (wrong board)
      3. Verify: Returns validation error
    Expected Result: Validation works
    Failure Indicators: Invalid moves accepted
    Evidence: .sisyphus/evidence/task-29-35-task-validation.txt
  ```

  **Commit**: YES (commit Tasks 28-35 as "Services: Auth + user + board + column + task + comment + label + attachment")

- [ ] 36. **Auth Controller**

  **What to do**:
  - Create `controllers/auth_controller.go`
  - Implement: Register, Login, Logout, RefreshToken, ForgotPassword, ResetPassword endpoints
  - Use `services/auth_service.go` for business logic
  - Use `utils/response.go` for responses
  - Validate inputs with `utils/validator.go`

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Authentication endpoints with JWT and refresh tokens
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: Task 38 (auth middleware depends on auth controller)
  - **Blocked By**: Tasks 28, 2, 4, 5 (auth service, response utils, validation, errors)

  **References**:
  - `controllers/` - Empty directory, create file
  - `main.go:35-40` - Existing route pattern (follow this)

  **Acceptance Criteria**:
  - [ ] Test file created: `controllers/auth_controller_test.go`
  - [ ] `go test ./controllers -v` → PASS (all auth endpoint tests)
  - [ ] All endpoints return standardized JSON
  - [ ] JWT tokens returned correctly

  **QA Scenarios**:

  ```
  Scenario: Register endpoint
    Tool: Bash (curl)
    Preconditions: App running
    Steps:
      1. Run: curl -X POST http://localhost:8080/api/auth/register -d '{"email":"new@example.com","password":"password123"}' -H "Content-Type: application/json"
      2. Verify: Returns 201 with { "success": true, "data": { "access_token": "...", "refresh_token": "..." } }
    Expected Result: Register works
    Failure Indicators: Wrong status, invalid response format
    Evidence: .sisyphus/evidence/task-36-register-endpoint.txt

  Scenario: Login endpoint
    Tool: Bash (curl)
    Preconditions: User exists
    Steps:
      1. Run: curl -X POST http://localhost:8080/api/auth/login -d '{"email":"test@example.com","password":"password123"}' -H "Content-Type: application/json"
      2. Verify: Returns 200 with tokens
    Expected Result: Login works
    Failure Indicators: Wrong credentials not handled
    Evidence: .sisyphus/evidence/task-36-login-endpoint.txt

  Scenario: Protected endpoint returns 401 without JWT
    Tool: Bash (curl)
    Preconditions: Auth middleware applied
    Steps:
      1. Run: curl http://localhost:8080/api/boards (no auth header)
      2. Verify: Returns 401 with error message
    Expected Result: Auth required
    Failure Indicators: Returns data without auth
    Evidence: .sisyphus/evidence/task-36-auth-required.txt
  ```

  **Commit**: NO (group with Tasks 37, 38)

- [ ] 37. **User Controller**

  **What to do**:
  - Create `controllers/user_controller.go`
  - Implement: GetMe, UpdateMe, DeleteMe, GetByID, Search endpoints
  - Use `services/user_service.go`
  - Handle soft delete for DeleteMe

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: User profile management with search
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 6 (with Task 36) | Sequential
  - **Blocks**: Task 51 (routes setup)
  - **Blocked By**: Tasks 29, 2, 4, 5 (user service, response utils, validation, errors)

  **References**:
  - `controllers/` - Create file

  **Acceptance Criteria**:
  - [ ] Test file created: `controllers/user_controller_test.go`
  - [ ] `go test ./controllers -v` → PASS (all user endpoint tests)
  - [ ] Soft delete works correctly

  **QA Scenarios**:

  ```
  Scenario: Get current user profile
    Tool: Bash (curl)
    Preconditions: User logged in
    Steps:
      1. Get token from login
      2. Run: curl http://localhost:8080/api/users/me -H "Authorization: Bearer <token>"
      3. Verify: Returns 200 with user data
    Expected Result: Get profile works
    Failure Indicators: Wrong user data
    Evidence: .sisyphus/evidence/task-37-get-profile.txt

  Scenario: Update profile
    Tool: Bash (curl)
    Preconditions: User logged in
    Steps:
      1. Run: curl -X PUT http://localhost:8080/api/users/me -d '{"name":"New Name"}' -H "Authorization: Bearer <token>" -H "Content-Type: application/json"
      2. Verify: Profile updated
    Expected Result: Update works
    Failure Indicators: Profile not updated
    Evidence: .sisyphus/evidence/task-37-update-profile.txt
  ```

  **Commit**: NO (group with Tasks 36, 38)

- [ ] 38. **Auth Middleware**

  **What to do**:
  - Create `middleware/auth.go`
  - Implement JWT validation for protected routes
  - Extract user_id from JWT claims and store in context
  - Return 401 for invalid/expired tokens
  - Use `utils/jwt.go` for token validation

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: JWT middleware with context injection
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: Tasks 39-50 (all controllers depend on auth middleware)
  - **Blocked By**: Tasks 1, 5 (JWT utils, errors)

  **References**:
  - `middleware/` - Create file

  **Acceptance Criteria**:
  - [ ] Test file created: `middleware/auth_test.go`
  - [ ] `go test ./middleware -v` → PASS (auth middleware tests)
  - [ ] Valid tokens pass, invalid fail
  - [ ] User ID stored in context

  **QA Scenarios**:

  ```
  Scenario: Valid token passes middleware
    Tool: Bash (curl)
    Preconditions: App running with auth middleware
    Steps:
      1. Run: curl http://localhost:8080/api/boards -H "Authorization: Bearer <valid_token>"
      2. Verify: Request passes middleware
    Expected Result: Valid token accepted
    Failure Indicators: Valid token rejected
    Evidence: .sisyphus/evidence/task-38-valid-token.txt

  Scenario: Invalid token returns 401
    Tool: Bash (curl)
    Preconditions: App running
    Steps:
      1. Run: curl http://localhost:8080/api/boards -H "Authorization: Bearer <invalid_token>"
      2. Verify: Returns 401
    Expected Result: Invalid token rejected
    Failure Indicators: Invalid token accepted
    Evidence: .sisyphus/evidence/task-38-invalid-token.txt
  ```

  **Commit**: YES (commit Tasks 36-38 as "Auth + User Controllers + Auth Middleware")

- [ ] 39-50. **Board, Column, Task, InProgressRule, ActivityLog, Comment, Label, Attachment, Notification, Widget, Search, WebSocket Controllers**
  *(12 controllers - extended functionality)*

  **What to do**:
  - `controllers/board_controller.go`: All board CRUD + members endpoints
  - `controllers/column_controller.go`: All column CRUD + reorder endpoints
  - `controllers/task_controller.go`: All task CRUD + move/reorder/assign/unassign/complete endpoints
  - `controllers/inprogressrule_controller.go`: Check active task, validate move, board in-progress tasks
  - `controllers/activitylog_controller.go`: Board activity history, task activity history
  - `controllers/comment_controller.go`: Comment CRUD on task
  - `controllers/label_controller.go`: Label CRUD + attach/detach to task
  - `controllers/attachment_controller.go`: Upload (S3), list, delete, download
  - `controllers/notification_controller.go`: Get all, mark read, read all, set preferences, snooze
  - `controllers/widget_controller.go`: Summary, deadlines, inprogress, progress endpoints
  - `controllers/search_controller.go`: Search tasks, boards, members
  - `controllers/websocket_controller.go`: WS /ws/boards/:id, WS /ws/notifications, WS /ws/chat/:taskId

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high` (T39, T41, T49, T50), `quick` (T40, T42, T43-T48)
    - Reason: Varying complexity across controllers
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 7 (T39-T43) and Wave 8 (T44-T50) | Sequential
  - **Blocks**: Task 51 (routes setup)
  - **Blocked By**: Tasks 29-35 (services), 38 (auth middleware)

  **References**:
  - `controllers/` - Create files
  - `services/S3_service.go` - For attachment upload pattern

  **Acceptance Criteria**:
  - [ ] Test files created for each controller
  - [ ] `go test ./controllers -v` → PASS (all controller tests)
  - [ ] All endpoints return standardized JSON

  **QA Scenarios**:

  ```
  Scenario: Board CRUD endpoints
    Tool: Bash (curl)
    Preconditions: Auth middleware, board service exist
    Steps:
      1. Create board: curl -X POST http://localhost:8080/api/boards -d '{"title":"My Board","color":"#FF5733"}' -H "Authorization: Bearer <token>" -H "Content-Type: application/json"
      2. Get boards: curl http://localhost:8080/api/boards -H "Authorization: Bearer <token>"
      3. Verify: Returns 201 with board data, 200 with list
    Expected Result: Board CRUD works
    Failure Indicators: Wrong status, invalid response
    Evidence: .sisyphus/evidence/task-39-50-board-crud.txt

  Scenario: WebSocket board updates broadcast
    Tool: Bash (wscat or similar)
    Preconditions: WebSocket server running
    Steps:
      1. Connect: wscat ws://localhost:8080/ws/boards/1?token=<jwt>
      2. Create task via API in another terminal
      3. Verify: WebSocket receives task update message
    Expected Result: Real-time updates work
    Failure Indicators: No WebSocket message
    Evidence: .sisyphus/evidence/task-39-50-websocket-broadcast.txt

  Scenario: File upload to S3
    Tool: Bash (curl)
    Preconditions: Auth, S3 configured
    Steps:
      1. Run: curl -X POST http://localhost:8080/api/tasks/1/attachments -F "file=@test.txt" -H "Authorization: Bearer <token>"
      2. Verify: Returns S3 URL
      3. Verify: File accessible via S3 URL
    Expected Result: Upload works
    Failure Indicators: Upload fails, no URL returned
    Evidence: .sisyphus/evidence/task-39-50-s3-upload.txt

  Scenario: Search tasks
    Tool: Bash (curl)
    Preconditions: Tasks exist
    Steps:
      1. Run: curl "http://localhost:8080/api/search/tasks?q=important" -H "Authorization: Bearer <token>"
      2. Verify: Returns matching tasks
    Expected Result: Search works
    Failure Indicators: No results, wrong filtering
    Evidence: .sisyphus/evidence/task-39-50-search.txt
  ```

  **Commit**: YES (commit Tasks 39-50 as "Controllers: Extended + WebSocket")

- [ ] 51. **Routes Setup**

  **What to do**:
  - Create `routes/routes.go`
  - Group routes: /api/auth, /api/users, /api/boards, /api/columns, /api/tasks, /api/comments, /api/labels, /api/attachments, /api/notifications, /api/widgets, /api/search, /api/activity, /ws
  - Apply auth middleware to protected routes
  - Apply CORS and logger middleware
  - Wire controllers to routes
  - WebSocket upgrade handlers

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Route wiring for 15+ controllers and WebSocket
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: Tasks 52, 53 (cron and main depend on routes)
  - **Blocked By**: Tasks 36-50 (controllers)

  **References**:
  - `routes/` - Create file
  - `main.go:26-45` - Route registration pattern

  **Acceptance Criteria**:
  - [ ] Test file created: `routes/routes_test.go`
  - [ ] `go test ./routes -v` → PASS (route tests)
  - [ ] All routes registered
  - [ ] Auth middleware applied correctly

  **QA Scenarios**:

  ```
  Scenario: All routes accessible
    Tool: Bash (curl)
    Preconditions: App running
    Steps:
      1. Run: curl http://localhost:8080/health (no auth)
      2. Run: curl http://localhost:8080/api/boards -H "Authorization: Bearer <token>"
      3. Verify: Both work correctly
    Expected Result: Routes wired correctly
    Failure Indicators: Routes not found, auth blocking public routes
    Evidence: .sisyphus/evidence/task-51-routes-wired.txt
  ```

  **Commit**: NO (group with Tasks 52-54)

- [ ] 52. **Cron Jobs Setup**

  **What to do**:
  - Create `jobs/cron_jobs.go`
  - Implement: `SetupCronJobs()` function
  - Job 1: Task due reminders (fire WebSocket notification 1 hr before due)
  - Job 2: Soft-delete cleanup (permanently delete >30 days old)
  - Job 3: Notification cleanup (delete read notifications >7 days)
  - Use `github.com/robfig/cron/v3` package
  - Register cron in `main.go`

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Cron jobs with scheduling and WebSocket notifications
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential
  - **Blocks**: Task 53 (main update)
  - **Blocked By**: Tasks 50 (WebSocket controller for notifications)

  **References**:
  - `jobs/` - Empty directory, create file
  - `main.go:15-45` - App initialization (add cron setup)

  **Acceptance Criteria**:
  - [ ] Test file created: `jobs/cron_jobs_test.go`
  - [ ] `go test ./jobs -v` → PASS (cron job tests)
  - [ ] Cron jobs registered on startup
  - [ ] Jobs execute at correct intervals

  **QA Scenarios**:

  ```
  Scenario: Task due reminder job
    Tool: Bash (logs)
    Preconditions: Cron running, task with due date exists
    Steps:
      1. Wait for cron job to execute (1 hr before due)
      2. Check logs: Verify WebSocket notification sent
      3. Verify: Notification contains task details
    Expected Result: Reminders work
    Failure Indicators: Job doesn't run, no notification sent
    Evidence: .sisyphus/evidence/task-52-due-reminder.txt

  Scenario: Soft-delete cleanup job
    Tool: Bash (logs)
    Preconditions: Cron running, soft-deleted items exist
    Steps:
      1. Wait for cron job to execute (daily)
      2. Check logs: Verify cleanup job ran
      3. Check DB: Soft-deleted items >30 days removed
    Expected Result: Cleanup works
    Failure Indicators: Job doesn't run, items not deleted
    Evidence: .sisyphus/evidence/task-52-cleanup.txt
  ```

  **Commit**: NO (group with Tasks 53-54)

- [ ] 53. **Update main.go**

  **What to do**:
  - Import routes: `routes.SetupRoutes(app)`
  - Import jobs: `jobs.SetupCronJobs(app)`
  - Update middleware usage if needed
  - Ensure DB and S3 connect before routes/jobs

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Simple main.go update
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 9 (with Tasks 51-52) | Sequential
  - **Blocks**: Task 54 (health checks)
  - **Blocked By**: Tasks 51-52 (routes and cron)

  **References**:
  - `main.go:15-46` - Existing main.go (enhance)

  **Acceptance Criteria**:
  - [ ] `go run main.go` starts without errors
  - [ ] Routes and cron jobs registered
  - [ ] Health check works

  **QA Scenarios**:

  ```
  Scenario: App starts successfully
    Tool: Bash (go run)
    Preconditions: All code compiled
    Steps:
      1. Run: go run main.go
      2. Verify: Server starts on port 8080
      3. Run: curl http://localhost:8080/health
    Expected Result: App starts and health check works
    Failure Indicators: App fails to start, health check 404
    Evidence: .sisyphus/evidence/task-53-app-starts.txt
  ```

  **Commit**: NO (group with Task 54)

- [ ] 54. **Health Check Endpoints**

  **What to do**:
  - Add `/health/db` - check database connection
  - Add `/health/s3` - check S3 connection
  - Update `/health` - return full health status
  - Return JSON: `{ "status": "healthy", "checks": { "db": "ok", "s3": "ok" } }`

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Simple health check endpoints
  - **Skills**: []
    - No special skills needed

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 9 (with Tasks 51-53) | Sequential
  - **Blocks**: Final verification wave
  - **Blocked By**: Tasks 51-53 (routes, cron, main)

  **References**:
  - `main.go:35-40` - Existing health endpoint (enhance)

  **Acceptance Criteria**:
  - [ ] Test file created: `main_test.go` or separate health_test.go
  - [ ] `go test . -v` → PASS (health check tests)
  - [ ] All health endpoints return correct status

  **QA Scenarios**:

  ```
  Scenario: Health check endpoints
    Tool: Bash (curl)
    Preconditions: App running
    Steps:
      1. Run: curl http://localhost:8080/health
      2. Verify: Returns 200 with full health status
      3. Run: curl http://localhost:8080/health/db
      4. Verify: Returns DB connection status
    Expected Result: Health checks work
    Failure Indicators: Health check fails, wrong format
    Evidence: .sisyphus/evidence/task-54-health-checks.txt
  ```

  **Commit**: YES (commit Tasks 51-54 as "Integration: Routes + Cron + Main + Health Checks")

---

## Final Verification Wave

> 4 review agents run in PARALLEL. ALL must APPROVE. Rejection → fix → re-run.

- [ ] F1. **Run All Tests & Verify Coverage** — `deep`
  Run `go test ./... -cover` and verify coverage ≥70%. Check for any test failures. Fix failing tests before proceeding.

- [ ] F2. **Integration Test Suite** — `deep`
  Run full integration tests simulating real workflows: register → login → create board → add task → move task → comment → delete.

- [ ] F3. **Final Verification Wave** — `oracle`
  Read the plan end-to-end. For each "Must Have": verify implementation exists (read file, curl endpoint, run command). For each "Must NOT Have": search codebase for forbidden patterns — reject with file:line if found. Check evidence files exist in .sisyphus/evidence/. Compare deliverables against plan.
  Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT: APPROVE/REJECT`

- [ ] F4. **Code Quality Review** — `unspecified-high`
  Run `go build`, `go vet`, and linter. Review all changed files for: empty catches, console.log in prod, commented-out code, unused imports. Check AI slop: excessive comments, over-abstraction, generic names.
  Output: `Build [PASS/FAIL] | Vet [PASS/FAIL] | Lint [PASS/FAIL] | Files [N clean/N issues] | VERDICT`

---

## Commit Strategy

- **1**: `feat(foundation): JWT, response, password, validation, errors, middleware` — utils/, middleware/
- **2**: `feat(models): GORM models with relationships` — models/
- **3**: `feat(migrations): Versioned GORM migrations` — migrations/
- **4**: `feat(repositories): Data access layer` — repositories/
- **5**: `feat(services): Business logic layer` — services/
- **6**: `feat(controllers): Auth and User controllers` — controllers/auth_controller.go, controllers/user_controller.go, middleware/auth.go
- **7**: `feat(controllers): Board, Column, Task, Activity controllers` — controllers/board_controller.go, controllers/column_controller.go, controllers/task_controller.go, controllers/inprogressrule_controller.go, controllers/activitylog_controller.go
- **8**: `feat(controllers): Extended + WebSocket controllers` — controllers/comment_controller.go, controllers/label_controller.go, controllers/attachment_controller.go, controllers/notification_controller.go, controllers/widget_controller.go, controllers/search_controller.go, controllers/websocket_controller.go
- **9**: `feat(integration): Routes, Cron, Main, Health checks` — routes/, jobs/, main.go

---

## Success Criteria

### Verification Commands
```bash
# Run all tests
go test ./... -v

# Check coverage
go test ./... -cover
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

# Build
go build

# Vet
go vet ./...

# Run application
go run main.go
```

### Final Checklist
- [ ] All 15 controllers implemented with tests passing
- [ ] All 7 GORM models created with migrations
- [ ] JWT authentication working (register, login, refresh, logout)
- [ ] WebSocket broadcasting real-time changes
- [ ] All endpoints return standardized JSON responses
- [ ] Test coverage ≥70%
- [ ] Cron jobs registered and logging
- [ ] Health check endpoints working
- [ ] All "Must Have" present
- [ ] All "Must NOT Have" absent

> EVERY task MUST have: Recommended Agent Profile + Parallelization info + QA Scenarios.
> **A task WITHOUT QA Scenarios is INCOMPLETE. No exceptions.**

