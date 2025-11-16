# CLAUDE.md - LUMA DMX Lighting Controller

> **Last Updated**: 2025-11-16
> **Purpose**: Guide for AI assistants working with the LUMA codebase

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Codebase Structure](#codebase-structure)
4. [Key Conventions](#key-conventions)
5. [Development Workflows](#development-workflows)
6. [Common Tasks](#common-tasks)
7. [Important Patterns](#important-patterns)
8. [Testing & Deployment](#testing--deployment)
9. [Troubleshooting](#troubleshooting)

---

## Project Overview

LUMA is a DMX-512 lighting controller with a Go backend and React frontend. It manages lighting fixtures, presets, and shows while controlling a DMX USB interface for real-time lighting control.

### Core Capabilities

- **Fixture Management**: Define lighting devices with configurable channels (DMX 1-512)
- **Presets**: Save and recall DMX channel states
- **Shows**: Sequential playback of presets with timing and fade control
- **Real-time Control**: Live DMX output via USB with WebSocket updates
- **Embedded UI**: React frontend compiled into Go binary using `go:embed`

### Technology Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Backend | Go | 1.24.3 |
| Web Framework | Fiber | v2.52.8 |
| Frontend | React | 19 |
| Build Tool | Vite | 6 |
| Styling | Tailwind CSS | 4 |
| State Management | TanStack Query | 5.80.6 |
| WebSocket | react-use-websocket | 4.13.0 |
| Persistence | YAML | v3 |
| Serial/DMX | go.bug.st/serial | 1.6.4 |

---

## Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    React Frontend (Vite)                     │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐  │
│  │ Fixtures │ Presets  │  Shows   │  Mixer   │Performance│  │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘  │
│         │ HTTP (REST)       │ WebSocket (Real-time)         │
└─────────┼───────────────────┼───────────────────────────────┘
          │                   │
          ▼                   ▼
┌─────────────────────────────────────────────────────────────┐
│                      Go Backend (Fiber)                      │
│  ┌────────────────┐  ┌──────────────┐  ┌────────────────┐  │
│  │  REST API      │  │  WebSocket   │  │ DMX Controller │  │
│  │  /api/*        │  │  /ws/control │  │  (Serial)      │  │
│  └────────────────┘  └──────────────┘  └────────────────┘  │
│           │                 │                   │            │
│           └────────┬────────┴───────────────────┘            │
│                    ▼                                         │
│         ┌──────────────────────┐                            │
│         │   YAML Storage       │                            │
│         │ .data/project.yaml   │                            │
│         └──────────────────────┘                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  DMX USB Device  │
                    │  (Serial Port)   │
                    └──────────────────┘
```

### Data Flow Patterns

#### REST API Flow (CRUD Operations)
```
UI Component → TanStack Query Hook → ky HTTP Client →
  Fiber Handler → Storage Validation → YAML Write →
  Cache Invalidation → UI Update
```

#### WebSocket Flow (Real-time Control)
```
Performance UI → WebSocket Send → Fiber WS Handler →
  DMX Controller → Serial Port → Broadcast State →
  All Connected Clients
```

#### Show Execution Flow
```
runShow() → WS Message → Backend Goroutine →
  Step Loop: [Blackout → Fade Channels → Wait Duration] →
  Monitor Broadcast (100ms) → UI Updates
```

---

## Codebase Structure

### Directory Layout

```
luma/
├── src/
│   ├── backend/          # Go server code
│   │   ├── main.go       # Entry point, server setup, signal handling
│   │   ├── api/          # REST endpoints (fixtures, presets, shows, projects, usb)
│   │   ├── ws/           # WebSocket handlers and broadcast logic
│   │   ├── dmx/          # DMX controller, serial communication, fade algorithms
│   │   ├── models/       # Data structures (Project, Fixture, Preset, Show)
│   │   ├── storage/      # YAML persistence, backups, validation
│   │   └── utils/        # Environment config, helpers
│   └── frontend/         # React application
│       ├── frontend.go   # Go embed wrapper (//go:embed dist/*)
│       ├── src/
│       │   ├── main.tsx          # React entry, providers setup
│       │   ├── app.tsx           # Tab-based UI layout
│       │   ├── modules/          # Feature modules (fixtures, presets, shows, etc.)
│       │   │   └── [module]/
│       │   │       ├── [module]-api.ts      # HTTP + TanStack Query hooks
│       │   │       ├── [module]-types.ts    # TypeScript interfaces
│       │   │       ├── components/          # UI components
│       │   │       └── views/               # Main tab content
│       │   ├── providers/        # React context providers (WebSocket, theme)
│       │   ├── components/       # Shared UI components (Shadcn)
│       │   ├── integration/      # API clients (ky, TanStack Query config)
│       │   └── lib/              # Utilities
│       ├── package.json
│       └── vite.config.ts
├── .data/                # Runtime data (auto-created)
│   ├── project.yaml      # Main project file
│   └── backups/          # Automatic backups (max 10)
├── dev.sh                # Development script (build frontend + run backend)
├── dockerfile            # Multi-stage Docker build
├── go.mod                # Go dependencies
└── readme.md             # User documentation
```

### Module Organization

Each frontend module follows this structure:

```typescript
modules/[name]/
├── [name]-api.ts         // API client + TanStack Query hooks
├── [name]-types.ts       // TypeScript interfaces
├── components/           // Reusable UI components
│   ├── [name]-table.tsx
│   ├── [name]-form.tsx
│   ├── new-[name]-sheet.tsx
│   └── edit-[name]-sheet.tsx
└── views/
    └── [name]-view.tsx   // Main tab content
```

**Frontend Modules**:
- **fixtures**: Define lighting devices with DMX channels
- **presets**: Save/recall DMX channel states
- **shows**: Sequential preset playback with timing
- **projects**: Root configuration (name, USB port)
- **mixer**: Real-time channel sliders
- **performance**: Live control grid (preset/show buttons)
- **usb**: USB port detection

---

## Key Conventions

### Backend Conventions

#### File Naming
- **Lowercase + underscores**: Package names use single lowercase words
- **Entity-based files**: `fixture.go`, `preset.go`, `show.go`
- **Test files**: `*_test.go` (though none currently exist)

#### Code Style
```go
// Public functions: PascalCase
func CreateFixture(c *fiber.Ctx) error { ... }

// Private functions: camelCase
func validateChannelAddress(addr int) error { ... }

// Constants: All caps with underscores
const DMXChannels = 512
const DMXBaudRate = 250000

// Struct fields: PascalCase with tags
type Fixture struct {
    ID       string    `json:"id" yaml:"id"`
    Name     string    `json:"name" yaml:"name"`
    Channels []Channel `json:"channels" yaml:"channels"`
}
```

#### Error Handling
```go
// Return fiber.Error for HTTP responses
if err := validate(); err != nil {
    return fiber.NewError(fiber.StatusBadRequest, err.Error())
}

// Use context for cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Recover from panics in goroutines
defer func() {
    if r := recover(); r != nil {
        log.Printf("Panic: %v", r)
    }
}()
```

#### Concurrency Patterns
```go
// Use sync.RWMutex for read-heavy operations
store.mu.RLock()
defer store.mu.RUnlock()

// Use sync.Map for concurrent client management
clients sync.Map

// Channels for broadcast
broadcast := make(chan Message, 100)
```

### Frontend Conventions

#### File Naming
- **kebab-case**: All files use kebab-case (`fixture-api.ts`, `preset-table.tsx`)
- **Component files**: End with `.tsx` for JSX components
- **Type files**: `[entity]-types.ts`
- **API files**: `[entity]-api.ts`

#### Code Style
```typescript
// Components: PascalCase
export function FixtureTable() { ... }

// Functions/variables: camelCase
const fetchFixtures = async () => { ... }

// Types/Interfaces: PascalCase
export interface Fixture { ... }
export type ChannelValue = { ... }

// Constants: UPPER_SNAKE_CASE or camelCase for objects
const MAX_DMX_CHANNELS = 512
const queryClient = new QueryClient()
```

#### TanStack Query Patterns

**Query Options** (reusable, exportable):
```typescript
export const fixturesQueryOptions = {
  all: () => ({
    queryKey: ["fixtures"],
    queryFn: getFixtures,
  }),
  detail: (id: string) => ({
    queryKey: ["fixtures", id],
    queryFn: () => getFixture(id),
  }),
}
```

**Custom Hooks**:
```typescript
export function useFixtures() {
  return useQuery(fixturesQueryOptions.all())
}

export function useCreateFixture() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: createFixture,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["fixtures"] })
    },
  })
}
```

**Cache Invalidation Strategy**:
- After create/update: Invalidate list + specific item
- After delete: Remove specific item query
- Use `invalidateQueries` for refetch, `removeQueries` for cleanup

#### WebSocket Patterns

```typescript
// Hook usage in components
const { applyPreset, runShow, stopShow } = useDmxWebSocket()

// Send actions
applyPreset(presetId)
runShow(showId, loop)
updateChannel(dmxAddress, value)

// Read state
const { presetApplied, showStarted, dmxState } = useDmxWebSocket()
```

### Validation Rules

#### DMX Protocol Constraints
```typescript
// Channel addresses: 1-512 (DMX-512 standard)
const isValidAddress = (addr: number) => addr >= 1 && addr <= 512

// Channel values: 0-255 (8-bit)
const isValidValue = (val: number) => val >= 0 && val <= 255

// Timing: Non-negative integers
const isValidDuration = (ms: number) => ms >= 0
```

#### Data Validation Locations
- **Frontend**: Zod schemas in forms (immediate feedback)
- **Backend**: Handler validation before storage (data integrity)
- **Storage**: Load-time validation (corruption detection)

---

## Development Workflows

### Initial Setup

```bash
# Clone repository
git clone <repository-url>
cd luma

# Install Go dependencies (automatic on build)
go mod download

# Install frontend dependencies
cd src/frontend
npm install
cd ../..
```

### Development Mode

**Option 1: Using dev.sh (recommended)**
```bash
# Build frontend + start backend
./dev.sh

# Skip frontend rebuild (faster iteration)
./dev.sh -n
```

**Option 2: Manual**
```bash
# Terminal 1: Frontend dev server with hot reload
cd src/frontend
npm run dev  # Runs on localhost:5173 with proxy to :3000

# Terminal 2: Backend
go run ./src/backend
```

**Accessing the application**:
- Frontend dev server: `http://localhost:5173` (with HMR)
- Backend direct: `http://localhost:3000` (with embedded UI)

### Making Changes

#### Backend Changes (Go)

1. **Adding a new API endpoint**:
   ```go
   // In src/backend/api/[entity].go
   func GetEntity(c *fiber.Ctx) error {
       // 1. Parse request
       id := c.Params("id")

       // 2. Validate input
       if id == "" {
           return fiber.NewError(fiber.StatusBadRequest, "ID required")
       }

       // 3. Fetch from storage
       entity, err := store.GetEntity(id)
       if err != nil {
           return fiber.NewError(fiber.StatusNotFound, err.Error())
       }

       // 4. Return JSON
       return c.JSON(entity)
   }

   // Register in main.go
   api.Get("/entities/:id", api.GetEntity)
   ```

2. **Adding a WebSocket message type**:
   ```go
   // In src/backend/ws/types.go
   const (
       MessageTypeNewAction = "new_action"
   )

   // In src/backend/ws/handlers.go
   func handleNewAction(provider *Provider, payload json.RawMessage) error {
       var req NewActionRequest
       if err := json.Unmarshal(payload, &req); err != nil {
           return err
       }

       // Process action
       // ...

       // Broadcast result
       BroadcastMessage(Message{
           Type: "action_result",
           Payload: result,
       })
       return nil
   }

   // Add to HandleMessage switch
   case MessageTypeNewAction:
       return handleNewAction(provider, msg.Payload)
   ```

3. **Modifying data models**:
   ```go
   // In src/backend/models/[entity].go
   type Entity struct {
       ID      string `json:"id" yaml:"id"`
       Name    string `json:"name" yaml:"name"`
       NewField string `json:"new_field" yaml:"new_field"` // Add field
   }

   // Update validation in storage/store.go
   func (s *Store) Validate() error {
       // Add new validation rules
   }
   ```

#### Frontend Changes (React/TypeScript)

1. **Adding a new form field**:
   ```typescript
   // Update type in [module]-types.ts
   export interface Entity {
     id: string
     name: string
     newField: string  // Add field
   }

   // Update Zod schema in components/[entity]-form.tsx
   const schema = z.object({
     name: z.string().min(1),
     newField: z.string().optional(),  // Add validation
   })

   // Add form field
   <FormField
     control={form.control}
     name="newField"
     render={({ field }) => (
       <FormItem>
         <FormLabel>New Field</FormLabel>
         <FormControl>
           <Input {...field} />
         </FormControl>
       </FormItem>
     )}
   />
   ```

2. **Adding a new API endpoint**:
   ```typescript
   // In modules/[module]/[module]-api.ts
   export async function getEntityDetail(id: string) {
     return await api.get(`entities/${id}`).json<Entity>()
   }

   // Add query option
   export const entityQueryOptions = {
     detail: (id: string) => ({
       queryKey: ["entities", id],
       queryFn: () => getEntityDetail(id),
     }),
   }

   // Create hook
   export function useEntityDetail(id: string) {
     return useQuery(entityQueryOptions.detail(id))
   }
   ```

3. **Adding a new module**:
   ```bash
   # Create module structure
   mkdir -p src/frontend/src/modules/newmodule/{components,views}
   touch src/frontend/src/modules/newmodule/newmodule-api.ts
   touch src/frontend/src/modules/newmodule/newmodule-types.ts
   touch src/frontend/src/modules/newmodule/views/newmodule-view.tsx
   ```

   ```typescript
   // Register in app.tsx
   import { NewModuleView } from "./modules/newmodule/views/newmodule-view"

   const tabs = [
     // ...
     { id: "newmodule", label: "New Module", icon: Icon, component: NewModuleView },
   ]
   ```

### Building for Production

```bash
# Build frontend
cd src/frontend
npm run build

# Build backend with embedded frontend
cd ../..
go build -o luma ./src/backend

# Run production binary
./luma
```

### Docker Build

```bash
# Build image
docker build -t luma:latest .

# Run container
docker run --rm -p 3000:3000 \
  -v $(pwd)/.data:/app/.data \
  --device=/dev/ttyUSB0 \
  luma:latest

# For macOS (USB passthrough requires special handling)
# See docker compose or native execution
```

---

## Common Tasks

### Adding a New Entity Type

**Example: Adding "Scenes"**

1. **Create backend model** (`src/backend/models/scene.go`):
   ```go
   package models

   type Scene struct {
       ID          string   `json:"id" yaml:"id"`
       Name        string   `json:"name" yaml:"name"`
       FixtureIDs  []string `json:"fixture_ids" yaml:"fixture_ids"`
   }
   ```

2. **Update Project model** (`src/backend/models/project.go`):
   ```go
   type Project struct {
       // ...existing fields...
       Scenes []Scene `json:"scenes" yaml:"scenes"`
   }
   ```

3. **Create API handlers** (`src/backend/api/scene.go`):
   ```go
   func GetScenes(c *fiber.Ctx) error { /* ... */ }
   func GetScene(c *fiber.Ctx) error { /* ... */ }
   func CreateScene(c *fiber.Ctx) error { /* ... */ }
   func UpdateScene(c *fiber.Ctx) error { /* ... */ }
   func DeleteScene(c *fiber.Ctx) error { /* ... */ }
   ```

4. **Register routes** (`src/backend/main.go`):
   ```go
   api := app.Group("/api")
   api.Get("/scenes", handlers.GetScenes)
   api.Get("/scenes/:id", handlers.GetScene)
   api.Post("/scenes", handlers.CreateScene)
   api.Put("/scenes/:id", handlers.UpdateScene)
   api.Delete("/scenes/:id", handlers.DeleteScene)
   ```

5. **Create frontend module**:
   ```bash
   mkdir -p src/frontend/src/modules/scenes/{components,views}
   ```

6. **Add TypeScript types** (`scenes-types.ts`):
   ```typescript
   export interface Scene {
     id: string
     name: string
     fixture_ids: string[]
   }
   ```

7. **Create API client** (`scenes-api.ts`):
   ```typescript
   export async function getScenes() {
     return await api.get("scenes").json<Scene[]>()
   }
   // ... other CRUD functions

   export const scenesQueryOptions = { /* ... */ }
   export function useScenes() { /* ... */ }
   export function useCreateScene() { /* ... */ }
   ```

8. **Create UI components** (table, form, view)

9. **Register in app** (`app.tsx`):
   ```typescript
   import { ScenesView } from "./modules/scenes/views/scenes-view"

   const tabs = [
     // ...
     { id: "scenes", label: "Scenes", icon: Layers, component: ScenesView },
   ]
   ```

### Adding a New WebSocket Action

**Example: Adding "pause_show"**

1. **Define message type** (`src/backend/ws/types.go`):
   ```go
   const (
       // ...existing types...
       MessageTypePauseShow = "pause_show"
   )
   ```

2. **Create handler** (`src/backend/ws/handlers.go`):
   ```go
   func handlePauseShow(provider *Provider, payload json.RawMessage) error {
       // Pause logic
       if provider.showController != nil {
           provider.showController.Pause()
       }

       BroadcastMessage(Message{
           Type: "show_paused",
           Payload: json.RawMessage(`{}`),
       })
       return nil
   }
   ```

3. **Register in switch** (`src/backend/ws/handlers.go`):
   ```go
   func HandleMessage(provider *Provider, msg Message) error {
       switch msg.Type {
       // ...existing cases...
       case MessageTypePauseShow:
           return handlePauseShow(provider, msg.Payload)
       }
   }
   ```

4. **Add frontend action** (`src/frontend/src/providers/ws-provider.tsx`):
   ```typescript
   const pauseShow = useCallback(() => {
     sendJsonMessage({
       type: "pause_show",
       payload: {},
     })
   }, [sendJsonMessage])

   // Add to context value
   return (
     <DmxWebSocketContext.Provider value={{
       // ...existing values...
       pauseShow,
     }}>
   ```

5. **Use in component**:
   ```typescript
   const { pauseShow } = useDmxWebSocket()

   <Button onClick={pauseShow}>Pause Show</Button>
   ```

### Debugging WebSocket Issues

1. **Enable browser DevTools**:
   - Network tab → WS filter → Click `/ws/control`
   - View messages in real-time

2. **Backend logging**:
   ```go
   // In ws/handlers.go
   log.Printf("Received message: type=%s, payload=%s", msg.Type, string(msg.Payload))
   ```

3. **Frontend logging**:
   ```typescript
   // In ws-provider.tsx
   useEffect(() => {
     if (lastMessage) {
       console.log("WS Message:", lastMessage)
     }
   }, [lastMessage])
   ```

4. **Check connection status**:
   ```typescript
   const { readyState } = useDmxWebSocket()
   console.log("WebSocket state:", readyState)
   // 0: CONNECTING, 1: OPEN, 2: CLOSING, 3: CLOSED
   ```

### Modifying DMX Output

**Change fade algorithm**:

```go
// In src/backend/dmx/controller.go

// Current: Linear fade
func (c *DMXController) FadeChannels(targets map[int]byte, duration time.Duration) {
    // Change from Linear to Quadratic
    curve := fadeQuadratic  // or fadeCubic, fadeSine, fadeExponential
    // ...
}
```

**Add new fade curve**:

```go
// In src/backend/dmx/controller.go
func fadeCustom(t float64) float64 {
    // t ranges from 0.0 to 1.0
    // Return value should also be 0.0 to 1.0
    return math.Pow(t, 3) * 0.5 + t * 0.5  // Custom curve
}
```

---

## Important Patterns

### Thread-Safe Storage Access

```go
// YAML Store with mutex protection
type Store struct {
    mu       sync.RWMutex
    filePath string
    project  *models.Project
}

// Read operations use RLock
func (s *Store) GetFixtures() []models.Fixture {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return append([]models.Fixture{}, s.project.Fixtures...)  // Return copy
}

// Write operations use Lock
func (s *Store) AddFixture(f models.Fixture) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.project.Fixtures = append(s.project.Fixtures, f)
    return s.save()
}
```

### Atomic File Writes

```go
// Always write to temp file first, then rename
func (s *Store) save() error {
    tmpPath := s.filePath + ".tmp"

    // Write to temp file
    if err := writeYAML(tmpPath, s.project); err != nil {
        return err
    }

    // Atomic rename (filesystem-level)
    return os.Rename(tmpPath, s.filePath)
}
```

### WebSocket Broadcast Pattern

```go
// Single broadcast channel, many clients
var broadcast = make(chan Message, 100)

// Central broadcast loop (single goroutine)
func handleBroadcast() {
    for msg := range broadcast {
        clients.Range(func(key, value interface{}) bool {
            conn := value.(*websocket.Conn)
            conn.WriteJSON(msg)  // Send to each client
            return true
        })
    }
}

// Anywhere in code: non-blocking send
func BroadcastMessage(msg Message) {
    select {
    case broadcast <- msg:
    default:
        log.Println("Broadcast channel full, dropping message")
    }
}
```

### DMX Continuous Sender

```go
// Background goroutine, triggered by changes or timer
func (c *DMXController) startSender() {
    ticker := time.NewTicker(40 * time.Millisecond)  // 25fps
    defer ticker.Stop()

    for {
        select {
        case <-c.sendNow:  // Triggered by SetChannel()
            c.sendFrame()
        case <-ticker.C:   // Regular interval
            c.sendFrame()
        case <-c.stopChan:
            return
        }
    }
}

// Non-blocking trigger
func (c *DMXController) SetChannel(ch int, val byte) {
    c.mu.Lock()
    c.channels[ch-1] = val
    c.mu.Unlock()

    select {
    case c.sendNow <- struct{}{}:  // Signal update
    default:  // Don't block if sender is busy
    }
}
```

### React Query Cache Invalidation

```typescript
// After successful mutation
export function useUpdateFixture() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: updateFixture,
    onSuccess: (data, variables) => {
      // Invalidate list (triggers refetch)
      queryClient.invalidateQueries({ queryKey: ["fixtures"] })

      // Invalidate specific item
      queryClient.invalidateQueries({ queryKey: ["fixtures", variables.id] })

      // Optional: Optimistic update (instant UI)
      queryClient.setQueryData(["fixtures", variables.id], data)
    },
  })
}

// After delete
export function useDeleteFixture() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: deleteFixture,
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ["fixtures"] })
      queryClient.removeQueries({ queryKey: ["fixtures", id] })  // Clean up
    },
  })
}
```

### Form Validation with Zod

```typescript
// Define schema
const fixtureSchema = z.object({
  name: z.string().min(1, "Name is required"),
  channels: z.array(
    z.object({
      name: z.string().min(1),
      dmx_address: z.number().int().min(1).max(512),
      min_value: z.number().int().min(0).max(255),
      max_value: z.number().int().min(0).max(255),
    })
  ).min(1, "At least one channel required"),
}).refine(
  (data) => data.channels.every(ch => ch.min_value <= ch.max_value),
  { message: "Min value must be <= max value" }
)

// Use in form
const form = useForm({
  resolver: zodResolver(fixtureSchema),
  defaultValues: {
    name: "",
    channels: [],
  },
})
```

### Environment Configuration

```go
// In src/backend/utils/env.go
type Config struct {
    ServerPort  string
    DataFile    string
    EnableDMX   bool
}

func LoadConfig() Config {
    return Config{
        ServerPort: getEnv("SERVER_PORT", ":3000"),
        DataFile:   getEnv("DATA_FILE", ".data/project.yaml"),
        EnableDMX:  getBoolEnv("ENABLE_DMX", true),
    }
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}
```

---

## Testing & Deployment

### Running Tests

**Backend**:
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package
go test ./src/backend/models

# Verbose output
go test -v ./...
```

**Frontend**:
```bash
cd src/frontend

# Run tests (if configured)
npm test

# Run with coverage
npm test -- --coverage
```

### Building Production Binary

```bash
# Build frontend
cd src/frontend
npm ci
npm run build

# Build backend with embedded UI
cd ../..
go build -o luma ./src/backend

# Run
./luma
```

### Docker Deployment

**Build**:
```bash
docker build -t luma:latest .
```

**Run**:
```bash
docker run -d \
  --name luma \
  -p 3000:3000 \
  -v $(pwd)/.data:/app/.data \
  --device=/dev/ttyUSB0:/dev/ttyUSB0 \
  -e ENABLE_DMX=true \
  luma:latest
```

**Docker Compose**:
```yaml
version: '3.8'
services:
  luma:
    build: .
    ports:
      - "3000:3000"
    volumes:
      - ./.data:/app/.data
    devices:
      - /dev/ttyUSB0:/dev/ttyUSB0
    environment:
      - SERVER_PORT=:3000
      - ENABLE_DMX=true
    restart: unless-stopped
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `:3000` | HTTP server listen address |
| `DATA_FILE` | `.data/project.yaml` | Project storage file path |
| `ENABLE_DMX` | `true` | Enable/disable DMX hardware output |
| `CORS_ORIGINS` | `*` | CORS allowed origins (comma-separated) |

### Deployment Checklist

- [ ] Frontend built with `npm run build`
- [ ] Production build tested locally
- [ ] Environment variables configured
- [ ] DMX USB device accessible (check permissions)
- [ ] `.data` directory mounted/persistent
- [ ] Backups directory configured
- [ ] Firewall allows port 3000
- [ ] CORS origins set (if not using `*`)
- [ ] Graceful shutdown tested (SIGTERM/SIGINT)

---

## Troubleshooting

### Common Issues

#### Issue: DMX not outputting

**Symptoms**: UI works, but lights don't respond

**Checks**:
1. Verify USB device connected:
   ```bash
   # Linux
   ls -l /dev/ttyUSB* /dev/ttyACM*

   # macOS
   ls -l /dev/cu.usbserial*

   # Windows
   # Check Device Manager → Ports (COM & LPT)
   ```

2. Check permissions:
   ```bash
   # Linux: Add user to dialout group
   sudo usermod -a -G dialout $USER
   # Log out and back in

   # Or run with sudo (not recommended for production)
   sudo ./luma
   ```

3. Verify `ENABLE_DMX` environment variable:
   ```bash
   export ENABLE_DMX=true
   ./luma
   ```

4. Check backend logs for serial port errors:
   ```
   Failed to open serial port: permission denied
   DMX controller disabled
   ```

5. Test with different USB port or cable

#### Issue: Frontend not loading

**Symptoms**: Blank page or 404 errors

**Checks**:
1. Verify frontend was built:
   ```bash
   ls -l src/frontend/dist/
   # Should contain index.html, assets/, etc.
   ```

2. Rebuild frontend:
   ```bash
   cd src/frontend
   rm -rf dist node_modules
   npm install
   npm run build
   ```

3. Check Go embed:
   ```bash
   # Rebuild backend after frontend changes
   go build -o luma ./src/backend
   ```

4. Verify server started:
   ```
   Server started on :3000
   ```

5. Check browser console for errors

#### Issue: WebSocket disconnects frequently

**Symptoms**: "Connection lost" messages, WsStatus shows red

**Checks**:
1. Check firewall/proxy settings (WebSocket upgrade required)

2. Verify browser DevTools → Network → WS shows upgrade:
   ```
   Status: 101 Switching Protocols
   Upgrade: websocket
   ```

3. Check backend logs for connection errors

4. Reduce message frequency if network is slow

5. Adjust ping/pong timeout in `ws/routes.go`:
   ```go
   ReadTimeout:  60 * time.Second,  // Increase if needed
   WriteTimeout: 10 * time.Second,
   ```

#### Issue: YAML file corrupted

**Symptoms**: Server fails to start with validation errors

**Fixes**:
1. Check backups:
   ```bash
   ls -lt .data/backups/
   # Copy most recent backup
   cp .data/backups/project_*.yaml.backup .data/project.yaml
   ```

2. Validate YAML manually:
   ```bash
   # Check syntax
   yamllint .data/project.yaml

   # Or use Go to validate
   go run ./src/backend -validate
   ```

3. Start with fresh project:
   ```bash
   mv .data/project.yaml .data/project.yaml.bad
   ./luma  # Will create new default project
   ```

#### Issue: Show playback stutters

**Symptoms**: Lights flicker, timing inconsistent

**Checks**:
1. System load (CPU usage):
   ```bash
   top  # or htop
   # Check if luma process is being throttled
   ```

2. Serial port buffer overflow (backend logs):
   ```
   DMX send error: resource temporarily unavailable
   ```

3. Reduce show complexity (fewer steps, longer durations)

4. Check fade algorithm (simpler = faster):
   ```go
   // Use Linear instead of Exponential
   curve := fadeLinear
   ```

5. Verify USB cable quality (avoid long/cheap cables)

#### Issue: Build fails

**Symptoms**: `go build` or `npm build` errors

**Go build errors**:
```bash
# Update dependencies
go mod tidy
go mod download

# Clear cache
go clean -cache -modcache
go build -o luma ./src/backend
```

**npm build errors**:
```bash
cd src/frontend

# Clear cache
rm -rf node_modules package-lock.json
npm install

# Check Node version (should be recent LTS)
node --version  # 18.x or newer recommended

# Build with verbose output
npm run build -- --debug
```

### Debugging Techniques

#### Backend Debugging

1. **Add logging**:
   ```go
   import "log"

   log.Printf("Debug: fixture=%+v", fixture)
   log.Printf("Channel %d set to %d", ch, val)
   ```

2. **Use delve debugger**:
   ```bash
   go install github.com/go-delve/delve/cmd/dlv@latest
   dlv debug ./src/backend

   # Set breakpoint
   (dlv) break main.main
   (dlv) continue
   ```

3. **Inspect WebSocket messages**:
   ```go
   // In ws/handlers.go
   func HandleMessage(provider *Provider, msg Message) error {
       log.Printf("WS Message: type=%s payload=%s", msg.Type, msg.Payload)
       // ...
   }
   ```

#### Frontend Debugging

1. **React DevTools**: Install browser extension, inspect component tree

2. **TanStack Query DevTools**: Already enabled in dev mode (bottom-left icon)

3. **Console logging**:
   ```typescript
   console.log("Fixture:", fixture)
   console.table(fixtures)  // Nice table view
   ```

4. **Network inspection**: DevTools → Network → filter by `/api` or `WS`

5. **WebSocket message logging**:
   ```typescript
   // In ws-provider.tsx
   useEffect(() => {
     if (lastMessage) {
       console.log("WS Received:", JSON.parse(lastMessage.data))
     }
   }, [lastMessage])
   ```

---

## Git Workflow

### Branch Strategy

This project uses feature branches:

```bash
# Create feature branch
git checkout -b feature/new-fade-algorithm

# Make changes, commit
git add .
git commit -m "Add exponential fade algorithm"

# Push to remote
git push -u origin feature/new-fade-algorithm

# Create pull request (via GitHub UI or gh CLI)
gh pr create --title "Add exponential fade" --body "..."
```

### Commit Message Convention

```bash
# Format: [type] - Description
# Types: feat, fix, refacto, docs, test, chore

git commit -m "[feat] - Add exponential fade to DMX controller"
git commit -m "[fix] - Fixed default value for show loop"
git commit -m "[refacto] - Renamed delay into duration"
```

**Recent examples**:
- `[fix] - Removed dmx env variable`
- `[refacto] - Renamed driver to dmx`
- `[feat] - Added fade to show`

### Pre-Commit Checks

Before committing:

1. **Backend**: Ensure builds without errors
   ```bash
   go build ./src/backend
   ```

2. **Frontend**: Ensure builds without errors
   ```bash
   cd src/frontend && npm run build
   ```

3. **Formatting**:
   ```bash
   # Go: Use gofmt
   gofmt -w .

   # TypeScript: Use Prettier (if configured)
   npm run format
   ```

---

## Additional Resources

### External Documentation

- **Fiber Framework**: https://docs.gofiber.io/
- **TanStack Query**: https://tanstack.com/query/latest
- **Radix UI**: https://www.radix-ui.com/
- **Tailwind CSS**: https://tailwindcss.com/
- **DMX-512 Protocol**: https://en.wikipedia.org/wiki/DMX512

### Project-Specific Files

- `readme.md` - User-facing documentation
- `LICENSE` - MIT License
- `dockerfile` - Container build instructions
- `dev.sh` - Development helper script

### Key Files to Review

When starting work on this codebase, review these files first:

1. **Backend entry point**: `src/backend/main.go`
2. **Data models**: `src/backend/models/project.go`
3. **Storage logic**: `src/backend/storage/store.go`
4. **WebSocket handling**: `src/backend/ws/handlers.go`
5. **DMX controller**: `src/backend/dmx/controller.go`
6. **Frontend entry**: `src/frontend/src/main.tsx`
7. **App structure**: `src/frontend/src/app.tsx`
8. **WebSocket provider**: `src/frontend/src/providers/ws-provider.tsx`

---

## AI Assistant Guidelines

### When Working on This Codebase

1. **Always check existing patterns** before adding new code
2. **Follow the module structure** for frontend features
3. **Use TanStack Query hooks** for all API calls (don't use fetch directly)
4. **Validate user input** on both frontend (Zod) and backend (handlers)
5. **Use WebSocket for real-time updates**, REST for CRUD
6. **Mutex-protect** all shared state in Go
7. **Return copies** from storage getters (don't expose internal state)
8. **Broadcast changes** to WebSocket clients after mutations
9. **Handle errors gracefully** with user-friendly messages
10. **Test both frontend and backend** after changes

### Code Review Checklist

When reviewing or writing code, check:

- [ ] Thread safety (mutexes, channels)
- [ ] Input validation (DMX ranges, required fields)
- [ ] Error handling (don't panic in handlers)
- [ ] Memory leaks (goroutines cleaned up)
- [ ] WebSocket message types documented
- [ ] TanStack Query cache invalidation after mutations
- [ ] TypeScript types match backend structs
- [ ] Form validation with Zod schemas
- [ ] Accessibility (ARIA labels, keyboard navigation)
- [ ] Responsive design (mobile-friendly)

### Making Pull Requests

1. Create feature branch from `main`
2. Make focused, single-purpose changes
3. Write descriptive commit messages
4. Test thoroughly (both dev and production builds)
5. Update CLAUDE.md if adding new patterns
6. Add comments for complex logic
7. Ensure no linter errors

---

**Questions or Issues?**

If you encounter problems not covered in this guide:
1. Check recent commits for similar changes
2. Review git history for the file you're modifying
3. Search GitHub issues for similar problems
4. Ask for clarification from the project maintainer

---

*This document is maintained as the single source of truth for AI assistants working with the LUMA codebase. Please keep it updated as the project evolves.*
