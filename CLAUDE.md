# CLAUDE.md — MasteryPath AI Context

> **Rule #1: Always produce PRODUCTION-READY output.** No placeholders, no TODO stubs, no "you can add this later." Every code change must be deployable.

## Project Identity

**MasteryPath** — A comprehensive skill tracking and goal management application.
**Stack**: GEM-R (Go · Expo · MongoDB · React Native)
**Author**: Deepak Thomas · **License**: MIT

## Quick Reference

| Layer | Technology | Version | Notes |
|-------|-----------|---------|-------|
| Backend Runtime | Go | 1.25.5 | Container-aware GOMAXPROCS, experimental json/v2 |
| HTTP Framework | Fiber | v2.52.10 | Express-inspired, consider v3 migration path |
| Database | MongoDB | Server 8.0+ compat | Using Go driver v1.17.6 (final 1.x) |
| Mobile Framework | Expo | SDK 54 | Last SDK with Legacy Architecture support |
| UI Library | React Native | 0.81.5 | New Architecture is default |
| React | React | 19.1 | Concurrent features, Suspense |
| State/Data | TanStack Query | v5 | Server state management |
| HTTP Client | Axios | 1.13.2 | Consider fetch API migration |
| API Docs | Swagger/Swag | v1.16.6 | OpenAPI 2.0 via godoc annotations |

## Project Structure

```
MasteryPath/
├── backend/                    # Go API Server
│   ├── cmd/api/main.go         # Entry point, DI, server bootstrap
│   ├── internal/
│   │   ├── config/config.go    # Environment config (MONGO_URI, PORT)
│   │   ├── database/mongo.go   # MongoDB connection
│   │   ├── handlers/           # HTTP handlers (GoalHandler, SkillHandler, ProgressHandler)
│   │   ├── models/             # Data models (Goal, Skill, ProgressItem)
│   │   └── routes/routes.go    # Route registration
│   ├── docs/                   # Generated Swagger docs
│   ├── scripts/                # Seed scripts
│   ├── Dockerfile              # Multi-stage build (Alpine)
│   └── Makefile                # swagger, run, build, generate, seed
├── frontend/                   # Expo / React Native
│   ├── app/                    # Expo Router (file-based routing)
│   │   ├── _layout.tsx         # Root layout (QueryClient, ThemeProvider)
│   │   ├── (tabs)/             # Tab navigation
│   │   │   ├── _layout.tsx     # Tab config (Goals, Explore)
│   │   │   ├── index.tsx       # Goals screen (CRUD + optimistic updates)
│   │   │   └── explore.tsx     # Explore screen (placeholder)
│   │   └── modal.tsx           # Modal screen
│   ├── components/             # Reusable components
│   ├── constants/theme.ts      # Color palette + Fonts
│   └── hooks/                  # Custom hooks (useColorScheme, useThemeColor)
├── salesforce/                 # Salesforce Apex
│   ├── classes/                # AccountPreventDeleteHandler
│   ├── triggers/               # AccountPreventDelete trigger
│   └── tests/                  # AccountPreventDeleteTest
├── .ai/                        # AI context files (READ THESE FIRST)
├── .agent/                     # Agent configs & workflows
├── CLAUDE.md                   # THIS FILE — Primary AI context
├── README.md                   # User-facing README
├── CONTRIBUTING.md             # Contribution guidelines
└── DEPLOYMENT.md               # Production deployment guide
```

## Architecture Patterns (MUST FOLLOW)

### Backend (Go)

1. **Handler Pattern**: Each resource has its own handler struct that receives `*mongo.Collection` via constructor injection.
   - `NewGoalHandler(col)`, `NewSkillHandler(skillCol, metaCol)`, `NewProgressHandler(col)`
2. **Materialized Path**: Skills use an `ancestors` array for hierarchical tree queries. When updating a skill's parent, always cascade ancestor updates to descendants.
3. **Fibonacci Weightage**: ProgressItem uses Fibonacci values (1, 3, 5, 8, 13, 21) for weightage. Always validate via `ValidateWeightage()`.
4. **MongoDB Aggregation**: Use `$setWindowFields` for computed fields like `weight_percent`. Use `$lookup` for parent population.
5. **Swagger Annotations**: Every handler function MUST have complete Swagger/godoc annotations.
6. **Consistent Error Responses**: Always return `fiber.Map{"error": "message"}` for errors.
7. **Empty Slice Convention**: Return `[]Type{}` instead of `nil` for empty lists.

### Frontend (React Native / Expo)

1. **Expo Router**: File-based routing via `app/` directory. Tab navigation via `(tabs)/` group.
2. **TanStack Query**: ALL server data fetching MUST use `useQuery` / `useMutation`.
   - Use `queryClient.invalidateQueries()` on mutation success.
   - Implement optimistic updates for toggle-like operations.
3. **Theme System**: Use `Colors` and `Fonts` from `@/constants/theme`. Support light/dark mode.
4. **Functional Components Only**: No class components. Use hooks.
5. **API Layer**: API functions are defined inline in screens. Migrate to `services/` directory (see suggestions).

### API Contract

- Base: `GET /api`
- Goals: `GET|POST /api/goals`, `GET|PUT|DELETE /api/goals/:id`
- Skills: `GET|POST /api/skills`, `GET|PUT|DELETE /api/skills/:id`, `GET /api/skills/:id/tree`
- Progress: `POST /api/progress`, `GET /api/progress/skill/:skillId`, `PUT|DELETE /api/progress/:id`

## Critical Rules for AI

### Always
- ✅ Produce **production-ready** code — no stubs, no "add later"
- ✅ Add Swagger annotations to every new Go endpoint
- ✅ Use TanStack Query for all data fetching in frontend
- ✅ Return empty arrays `[]` not `null` from API list endpoints
- ✅ Validate all ObjectID parsing — handle hex errors
- ✅ Use `context.Background()` for DB operations (or accept context from request)
- ✅ Follow existing naming: camelCase (TS), snake_case (JSON/BSON), PascalCase (Go exports)
- ✅ Add `json` AND `bson` struct tags to all Go model fields
- ✅ Use TypeScript interfaces matching Go models
- ✅ Test with `go vet ./...` before considering backend work done
- ✅ Create model validation methods for business rules (like `ValidateWeightage`)

### Never
- ❌ Use `interface{}` in new code — use `any` (Go 1.18+)
- ❌ Ignore MongoDB `ObjectIDFromHex` errors (existing code does this — fix when touching)
- ❌ Hardcode API URLs — always use environment variables
- ❌ Use class components in React Native
- ❌ Skip error handling — every DB/API call must handle errors
- ❌ Use `log.Fatal` in handlers — only in startup code
- ❌ Create files without proper file headers / package declarations

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection string |
| `PORT` | `8080` | Backend server port |
| `EXPO_PUBLIC_API_URL` | `http://localhost:8080/api` | Frontend API base URL |

## Commands

```bash
# Backend
cd backend && go run cmd/api/main.go           # Run dev server
cd backend && make swagger                       # Regenerate Swagger docs
cd backend && make build                         # Build binary
cd backend && go vet ./...                       # Lint check
cd backend && go test ./...                      # Run tests

# Frontend
cd frontend && npx expo start                    # Dev server
cd frontend && npx expo start --android          # Android dev
cd frontend && npx expo start --ios              # iOS dev
cd frontend && npm run lint                      # ESLint
```

## Read More

| File | Purpose |
|------|---------|
| [.ai/PROJECT_CONTEXT.md](.ai/PROJECT_CONTEXT.md) | Deep architecture, data flow, patterns |
| [.ai/CODING_STANDARDS.md](.ai/CODING_STANDARDS.md) | Market-standard coding best practices |
| [.ai/ARCHITECTURE_SUGGESTIONS.md](.ai/ARCHITECTURE_SUGGESTIONS.md) | Improvement recommendations |
| [.ai/TECH_STACK_RELEASE_NOTES.md](.ai/TECH_STACK_RELEASE_NOTES.md) | Latest release notes & upgrade paths |
| [.ai/LESSONS_LEARNED.md](.ai/LESSONS_LEARNED.md) | AI lesson tracking log |
| [.agent/workflows/ai-skills.md](.agent/workflows/ai-skills.md) | AI skill maintenance workflow |
