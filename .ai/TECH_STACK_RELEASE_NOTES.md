# MasteryPath — Tech Stack Release Notes & Upgrade Paths

> Last updated: February 2026. Reference for AI to understand current versions, latest features, and planned migrations.

---

## Go 1.25.5 (Released Aug 2025, patch 1.25.5)

### Key Features to Leverage
| Feature | Status | How to Use |
|---------|--------|------------|
| Container-Aware `GOMAXPROCS` | ✅ Stable | Auto-detected in Docker/K8s — no code changes needed |
| Experimental `encoding/json/v2` | 🧪 Experimental | Set `GOEXPERIMENT=jsonv2` for 2-3x faster JSON decoding |
| Experimental Green Tea GC | 🧪 Experimental | Set `GOEXPERIMENT=greenteagc` for 10-40% GC overhead reduction |
| `testing/synctest` | ✅ GA | Use for concurrent code testing in isolated "bubbles" |
| `log/slog` | ✅ Stable (since 1.21) | Replace `log.Println` with structured `slog.Info` / `slog.Error` |
| DWARF v5 | ✅ Auto | Smaller binaries, faster linking — automatic |
| `go doc -http` | ✅ Stable | Local documentation server with browser launch |

### Upgrade to Go 1.26 (Expected Feb 2026)
- `new()` function will accept expressions
- Generic types can self-reference in type parameter lists
- **Action**: Update `go.mod` and `Dockerfile` when stable

---

## GoFiber v2.52.10 → v3.0.0 Migration Path

### Current: v2.52.10 (Nov 2025)
- Consider upgrading to **v2.52.11** (Jan 2026) — fixes CVE-2025-66630 (UUID predictability)

### Fiber v3.0.0 (Released Feb 2, 2026)
| Breaking Change | v2 | v3 |
|----------------|-----|-----|
| Context method | `c.JSON()` | `c.JSON()` (mostly compatible) |
| Handler signature | `func(c *fiber.Ctx) error` | `func(c fiber.Ctx) error` (value, not pointer) |
| `interface{}` | Used | Replaced with `any` |
| Unknown methods | Returns 400 | Returns 501 |

### Migration Recommendation
> **Short-term**: Upgrade to v2.52.11 for security patch.
> **Medium-term**: Plan v3 migration after v3.1 stabilizes (1-2 months).

---

## Expo SDK 54 (Released Sep 2025)

### Key Features Active in This Project
| Feature | Status | Notes |
|---------|--------|-------|
| React Native 0.81 | ✅ Active | New Architecture is default |
| React 19.1 | ✅ Active | Concurrent rendering, Suspense |
| Expo Router v6 | ✅ Active | File-based routing with link previews |
| iOS 26 / Liquid Glass | ✅ Available | Premium native iOS design |
| Android 16 (API 36) | ✅ Active | Edge-to-edge default |
| `expo-file-system` (new) | ✅ Stable | Replaces old file system API |
| `expo-app-integrity` | 🆕 Available | DeviceCheck (iOS) + Play Integrity (Android) |
| React Compiler | ✅ Default | Auto-memoization in templates |
| `lightningcss` | ✅ Default | Fast CSS autoprefixing |

### ⚠️ Critical: SDK 55 Migration
> **Expo SDK 54 is the LAST version supporting Legacy Architecture.**
> SDK 55 (beta Jan 2026) requires New Architecture exclusively.
> Ensure all third-party packages support New Architecture before upgrading.

---

## React Native 0.81.5 (Aug 2025)

### Key Changes
- **Android 16 (API 36)** official support
- **`<SafeAreaView>` deprecated** — use `react-native-safe-area-context` instead
- **Precompiled iOS builds** (experimental) — faster CI
- **New Architecture is DEFAULT** since 0.76

### Migration Note for 0.82+
React Native 0.82 runs **entirely** on New Architecture. Legacy bridge is removed.

---

## MongoDB Go Driver v1.17.6 (Oct 2025)

### Features Available
| Feature | Status | Notes |
|---------|--------|-------|
| OIDC Authentication | ✅ Available | For enterprise auth scenarios |
| Queryable Encryption (Range) | ✅ Available | Requires MongoDB Server 8.0+ |
| MongoDB 8.0 Compatibility | ✅ Active | Full support except client bulk write |
| `DropOneWithKey` / `DropWithKey` | ✅ Available | Drop indexes by key spec |

### ⚠️ Critical: v1.17 is Final 1.x Release
> **MongoDB Go Driver 2.0** was released Jan 16, 2025.
> v1.17.x receives only security/bug fixes for 1 year (until Sep 2025).
> **Plan migration to v2.x** — API changes include more idiomatic Go patterns.

---

## TanStack React Query v5.90.20

### Current Best Practices
| Practice | Status in Project | Action |
|---------|------------------|--------|
| `queryOptions` API | ❌ Not used | Adopt for centralized config |
| Query Key Factories | ❌ Not used | Create `queryKeys.ts` |
| Optimistic Updates | ✅ Used (goals toggle) | Extend to other mutations |
| `staleTime` config | ❌ Default (0ms) | Configure per-resource |
| `gcTime` config | ❌ Default (5min) | Fine-tune for UX |
| Error Boundaries | ❌ Not used | Add `throwOnError` + boundaries |
| Suspense integration | ❌ Not used | Enable for premium loading UX |
| DevTools | ❌ Not installed | Add `@tanstack/react-query-devtools` |

---

## Axios 1.13.2

### Current Status
Stable, no critical issues. Consider:
- Adding request/response interceptors for auth tokens
- Adding retry logic with `axios-retry`
- **Future**: Evaluate replacing with `fetch` API (native in React Native)

---

## Version Upgrade Priority Matrix

| Priority | Package | Current | Target | Effort | Risk |
|----------|---------|---------|--------|--------|------|
| 🔴 HIGH | Fiber | 2.52.10 | 2.52.11 | Low | Low (patch) |
| 🟡 MED | MongoDB Driver | 1.17.6 | 2.x | High | Medium |
| 🟡 MED | Dockerfile Go | 1.21 | 1.25 | Low | Low |
| 🟢 LOW | Fiber | 2.52.11 | 3.x | Medium | Medium |
| 🟢 LOW | Expo SDK | 54 | 55 | Medium | Medium |
