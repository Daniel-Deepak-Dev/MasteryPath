# MasteryPath — AI Lessons Learned Log

> **Purpose**: Track mistakes, discoveries, and patterns learned during AI-assisted development. Each entry has a category, date, lesson, and resolution. New lessons should be PREPENDED to the top.
>
> **For AI**: Before making changes, scan this log for relevant past lessons to avoid repeating mistakes.

---

## How to Add a Lesson

```markdown
### [YYYY-MM-DD] Category: Brief Title
**Context**: What were you doing?
**Problem**: What went wrong or what was discovered?
**Resolution**: How was it fixed?
**Prevention**: How to avoid this in the future?
**Tags**: `backend`, `frontend`, `database`, `config`, `deployment`, `testing`
```

---

## Lessons

### [2026-02-10] Architecture: Silently Ignored ObjectID Parsing Errors
**Context**: Reviewing `skill_handler.go` and `progress_handler.go` for code quality.
**Problem**: Multiple handlers use `objID, _ := primitive.ObjectIDFromHex(id)` — silently ignoring parse errors. If a user sends a non-hex string, `objID` becomes a zero-value ObjectID, leading to unexpected query behavior instead of a clear 400 error.
**Resolution**: Always check the error: `if err != nil { return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"}) }`.
**Prevention**: Add a lint rule or standardized `parseID(c)` helper that all handlers use.
**Tags**: `backend`, `error-handling`

---

### [2026-02-10] Architecture: Dockerfile Go Version Mismatch
**Context**: Reviewing build configuration.
**Problem**: `Dockerfile` uses `FROM golang:1.21-alpine` but `go.mod` specifies `go 1.25.5`. This can cause build failures or silently use older Go features.
**Resolution**: Update Dockerfile to `FROM golang:1.25-alpine AS builder`.
**Prevention**: When upgrading Go version in `go.mod`, always update `Dockerfile` to match.
**Tags**: `deployment`, `config`

---

### [2026-02-10] Data Pattern: SafeAreaView Deprecation in React Native 0.81
**Context**: Frontend uses `SafeAreaView` from `react-native` directly.
**Problem**: `<SafeAreaView>` from `react-native` is deprecated in RN 0.81+. Use `react-native-safe-area-context` instead for proper edge-to-edge rendering.
**Resolution**: Replace `import { SafeAreaView } from 'react-native'` with `import { SafeAreaView } from 'react-native-safe-area-context'`.
**Prevention**: Check deprecation warnings in React Native release notes before each SDK bump.
**Tags**: `frontend`, `deprecation`

---

### [2026-02-10] Pattern: MongoDB v1.17 is Final 1.x Release
**Context**: Reviewing dependency versions.
**Problem**: MongoDB Go Driver v1.17.x is the final 1.x series release — only receives security fixes until ~Sep 2025 (already expired). v2.0 released Jan 2025 with breaking API changes.
**Resolution**: Plan migration to v2.x driver with updated import paths and API patterns.
**Prevention**: Track driver EOL dates; add to upgrade priority matrix.
**Tags**: `backend`, `database`, `dependency`

---

### [2026-02-10] Pattern: Expo SDK 54 Legacy Architecture Deadline
**Context**: Reviewing Expo SDK release notes.
**Problem**: SDK 54 is the LAST Expo version that supports React Native Legacy Architecture. SDK 55+ will REQUIRE New Architecture exclusively.
**Resolution**: Ensure all third-party packages support New Architecture before upgrading to SDK 55.
**Prevention**: Before any SDK upgrade, audit all dependencies for New Architecture compatibility.
**Tags**: `frontend`, `dependency`, `migration`

---

## Statistics

| Category | Count |
|----------|-------|
| Backend | 3 |
| Frontend | 2 |
| Database | 1 |
| Deployment | 1 |
| Config | 1 |
| **Total** | **5** |
