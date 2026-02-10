# MasteryPath — Coding Standards & Best Practices

> Market-standard patterns derived from: Go official style guide, Uber Go Style, Fiber best practices, React Native community guidelines, Expo documentation, and TanStack Query patterns.

---

## Go Backend Standards

### File Organization
```
internal/
  config/       # One file per config source
  database/     # Connection + migration logic
  handlers/     # One handler struct per resource
  models/       # One file per domain entity
  routes/       # Route registration
  services/     # Business logic (to be added)
  middleware/   # HTTP middleware (to be added)
  apperror/     # Typed error definitions (to be added)
```

### Naming Conventions
| Scope | Convention | Example |
|-------|-----------|---------|
| Exported types | PascalCase | `GoalHandler`, `SkillPopulated` |
| Unexported | camelCase | `isValidCategory` |
| Files | snake_case | `goal_handler.go` |
| JSON tags | snake_case | `json:"created_at"` |
| BSON tags | snake_case | `bson:"parent_id"` |
| Constants | PascalCase | `MaxRetries` |
| Interfaces | -er suffix or descriptive | `GoalService`, `Validator` |
| Test files | `_test.go` suffix | `goal_handler_test.go` |

### Handler Pattern (MUST FOLLOW)
```go
// 1. Struct with dependencies
type ResourceHandler struct {
    collection *mongo.Collection
}

// 2. Constructor
func NewResourceHandler(col *mongo.Collection) *ResourceHandler {
    return &ResourceHandler{collection: col}
}

// 3. Handler methods — always return error
// 4. Swagger annotations REQUIRED
// @Summary Short description
// @Description Detailed description
// @Tags resource
// @Accept json
// @Produce json
func (h *ResourceHandler) GetAll(c *fiber.Ctx) error {
    // Use c.UserContext() for DB operations
    // Return empty slice, not nil
    // Use consistent error responses
}
```

### Error Handling Rules
```go
// ✅ DO: Always check errors explicitly
id := c.Params("id")
objectID, err := primitive.ObjectIDFromHex(id)
if err != nil {
    return c.Status(400).JSON(fiber.Map{"error": "Invalid ID format"})
}

// ❌ DON'T: Ignore errors
objID, _ := primitive.ObjectIDFromHex(id) // NEVER DO THIS
```

### Go Code Quality Checklist
- [ ] Run `go vet ./...` before every commit
- [ ] Run `gofmt` / `goimports` for formatting
- [ ] No `interface{}` — use `any`
- [ ] No `log.Fatal` in handlers
- [ ] All public functions have doc comments
- [ ] All handlers have Swagger annotations
- [ ] Empty slices initialized, not nil

---

## Frontend (React Native / Expo) Standards

### File Organization
```
frontend/
  app/              # Expo Router screens (file-based)
    (tabs)/          # Tab group
    _layout.tsx      # Root layout
  components/        # Reusable UI components
    ui/              # Primitive UI components
  constants/         # Theme, config values
  hooks/             # Custom React hooks
  services/          # API service layer (to be added)
  types/             # Shared TypeScript types (to be added)
  utils/             # Utility functions (to be added)
```

### Naming Conventions
| Scope | Convention | Example |
|-------|-----------|---------|
| Components | PascalCase | `GoalCard.tsx` |
| Hooks | camelCase with `use` prefix | `useGoals.ts` |
| Services | camelCase | `goalService.ts` |
| Types/Interfaces | PascalCase | `interface Goal {}` |
| Constants | UPPER_SNAKE or PascalCase | `BASE_API_URL`, `Colors` |
| Files | kebab-case | `goal-card.tsx` |
| Styles | camelCase | `styles.goalItem` |

### Component Pattern (MUST FOLLOW)
```tsx
// 1. Imports
import React, { useState } from 'react';
import { View, Text, StyleSheet } from 'react-native';

// 2. TypeScript interfaces
interface Props {
  title: string;
  onPress: () => void;
}

// 3. Functional component (NEVER class)
export default function GoalCard({ title, onPress }: Props) {
  // 4. Hooks first
  const [isActive, setIsActive] = useState(false);

  // 5. Handlers
  const handlePress = () => { /* ... */ };

  // 6. Render
  return (
    <View style={styles.container}>
      <Text>{title}</Text>
    </View>
  );
}

// 7. Styles at bottom
const styles = StyleSheet.create({
  container: { /* ... */ },
});
```

### TanStack Query Pattern (MUST FOLLOW)
```tsx
// ✅ DO: Use query options with key factory
const goalKeys = {
  all: ['goals'] as const,
  detail: (id: string) => ['goals', id] as const,
};

// ✅ DO: Invalidate queries on mutation success
const createMutation = useMutation({
  mutationFn: createGoalApi,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: goalKeys.all });
  },
});

// ✅ DO: Optimistic updates for toggle-like operations
// ✅ DO: Configure staleTime for data that doesn't change often
// ❌ DON'T: Use useEffect for data fetching
// ❌ DON'T: Store server data in useState
```

### TypeScript Rules
- All component props MUST have TypeScript interfaces
- API response types MUST mirror Go models
- Use `unknown` instead of `any` where possible
- Enable strict mode in `tsconfig.json`
- Use discriminated unions for state machines

---

## API Design Standards

### REST Conventions
| Operation | Method | URL Pattern | Status Codes |
|-----------|--------|-------------|-------------|
| List all | GET | `/api/resources` | 200, 500 |
| Get one | GET | `/api/resources/:id` | 200, 400, 404, 500 |
| Create | POST | `/api/resources` | 201, 400, 500 |
| Update | PUT | `/api/resources/:id` | 200, 400, 404, 500 |
| Delete | DELETE | `/api/resources/:id` | 200, 400, 404, 500 |

### Response Format
```json
// Success (list)
[{ "id": "...", "title": "..." }]

// Success (single)
{ "id": "...", "title": "..." }

// Success (action)
{ "message": "Resource deleted successfully" }

// Error
{ "error": "Descriptive error message" }
```

### JSON Naming
- All fields: `snake_case` (matching Go's BSON tags)
- IDs: Use `id` (mapped from MongoDB `_id`)
- Timestamps: ISO 8601 format via `time.Time`

---

## Git Workflow

### Commit Messages (Conventional Commits)
```
feat: add skill analytics dashboard
fix: handle null parent_id in skill update
refactor: extract API calls to service layer
docs: update CLAUDE.md with new patterns
test: add unit tests for goal handler
chore: upgrade Fiber to v2.52.11
```

### Branch Naming
```
feature/skill-analytics-dashboard
fix/null-parent-id-crash
refactor/api-service-layer
docs/ai-context-files
```
