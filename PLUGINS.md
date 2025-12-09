# GoAdmin Plugin System

## Overview

Plugins are Go packages compiled into the binary that extend GoAdmin's functionality. They implement a standard interface and interact through a well-defined API.

## Architecture

### Core Components

1. **Plugin Interface** (`app/plugins/plugin.go`)

   - Defines the contract all plugins must implement
   - Lifecycle methods: `Init()`, `Start()`, `Stop()`, `Reload()`
   - Metadata method: `Metadata()` returns plugin information

2. **Plugin Registry** (`app/plugins/registry.go`)

   - Global registry for plugin registration
   - Plugins self-register via `init()` functions
   - Thread-safe operations using `sync.RWMutex`

3. **Plugin Manager** (`app/plugins/registry.go`)

   - Loads all registered plugins on startup
   - Manages plugin lifecycle (start, stop, reload)
   - Tracks plugin state and errors

4. **Plugin Context** (`app/plugins/plugin.go`)
   - Provides plugins access to GoAdmin APIs
   - Includes cancellation context for graceful shutdown
   - 6 API surfaces available to plugins

### Plugin APIs

Plugins have access to these APIs through the `PluginContext`:

#### 1. EventBus API

- **Purpose**: Subscribe to and publish events
- **Methods**:
  - `Subscribe(eventType, handler)` - Listen for events
  - `Unsubscribe(eventType, handler)` - Stop listening
  - `Publish(eventType, data)` - Trigger events

**Available Events**:

- `player.connect` - Player joined server
- `player.disconnect` - Player left server
- `player.banned` - Player was banned
- `player.kicked` - Player was kicked
- `report.created` - Report submitted
- `report.actioned` - Report resolved

#### 2. Command API

- **Purpose**: Register custom in-game commands
- **Methods**:
  - `RegisterCommand(definition)` - Add custom command
  - `UnregisterCommand(name)` - Remove command
  - `ExecuteCommand(playerName, playerGUID, command, args)` - Trigger command programmatically

**Command Definition**:

```go
plugins.CommandDefinition{
    Name:        "hello",           // Command name (used as !hello)
    Usage:       "hello",            // Usage help text
    Description: "Says hello",       // Command description
    MinArgs:     0,                  // Minimum arguments required
    MaxArgs:     0,                  // Maximum arguments (-1 = unlimited)
    MinPower:    0,                  // Minimum group power level
    Permissions: []string{},         // Required permissions
    Handler:     func(playerName, playerGUID string, args []string) error {
        // Command logic here
        return nil
    },
}
```

**Features**:

- Automatic argument validation
- Power level checking (uses in-game group system)
- Permission validation (checks group permissions)
- Error messages sent to player via RCON
- Integrated with command handler (plugin commands checked first)

**Example**: Register `!hello` command that greets players, `!time` shows server time

#### 3. RCON API

- **Purpose**: Execute RCON commands on game servers
- **Methods**:
  - `SendCommand(command)` - Execute raw RCON command
  - `SendCommandWithTimeout(command, timeout)` - Execute with custom timeout
  - `GetStatus()` - Get server status information

**Features**:

- Direct access to RCON client
- Async command execution (non-blocking)
- Error handling and timeout control
- Full access to all RCON commands

**Examples**:

```go
// Send a message to all players
p.ctx.RCONAPI.SendCommand(`say "^2Server restart in 5 minutes"`)

// Send private message to player
p.ctx.RCONAPI.SendCommand(fmt.Sprintf(`tell %s "^2Welcome!"`, playerName))

// Kick a player
p.ctx.RCONAPI.SendCommand(fmt.Sprintf("clientkick %s", playerID))

// Change map
p.ctx.RCONAPI.SendCommand("map mp_crash")
```

#### 4. Database API

- **Purpose**: Direct database access using GORM
- **Methods**:
  - `GetDB()` - Get GORM database instance
  - `Query(sql, args...)` - Execute SELECT query
  - `Exec(sql, args...)` - Execute INSERT/UPDATE/DELETE

**Note**: Not yet implemented. Use standard Go database access for now.

**Example**: Query player statistics, store plugin data

#### 5. Webhook API

- **Purpose**: Dispatch custom webhook events
- **Methods**:
  - `Dispatch(eventType, data)` - Trigger webhook
  - `RegisterEvent(eventType)` - Register new event type

**Note**: Not yet implemented.

**Example**: Notify external services of plugin events

#### 6. Config API

- **Purpose**: Persistent plugin configuration storage
- **Methods**:
  - `Get(key)` - Retrieve config value
  - `Set(key, value)` - Store config value
  - `GetString(key)`, `GetInt(key)`, `GetBool(key)` - Typed getters

**Note**: Not yet implemented. Use environment variables or config files for now.

**Example**: Store greeting messages, thresholds, toggles

## Plugin Lifecycle

```
1. REGISTER → Plugin calls Registry.Register() in init()
2. IMPORT   → Plugin package imported in main.go
3. LOAD     → Manager.LoadAll() discovers registered plugins
4. INIT     → Init() called with PluginContext
5. START    → Start() called, subscriptions/registrations occur
6. RUNNING  → Plugin handles events and commands
7. RELOAD   → Reload() called to refresh config
8. STOP     → Stop() called, cleanup subscriptions
```

## Quick Start Guide

1. **Create your plugin file** (`plugins/myplugin/myplugin.go`):

   ```go
   package myplugin

   import "github.com/ethanburkett/goadmin/app/plugins"

   type MyPlugin struct {
       ctx *plugins.PluginContext
   }

   func (p *MyPlugin) Metadata() plugins.PluginMetadata {
       return plugins.PluginMetadata{
           ID:      "my-plugin",
           Name:    "My Plugin",
           Version: "1.0.0",
       }
   }

   func (p *MyPlugin) Init(ctx *plugins.PluginContext) error {
       p.ctx = ctx
       return nil
   }

   func (p *MyPlugin) Start() error {
       // Your plugin logic here
       return nil
   }

   func (p *MyPlugin) Stop() error { return nil }
   func (p *MyPlugin) Reload() error { return nil }

   func init() {
       plugins.Registry.Register(&MyPlugin{})
   }
   ```

2. **Run the auto-import script**:

   ```powershell
   .\scripts\build_plugins.ps1
   ```

3. **Rebuild and restart**:

   ```bash
   go build -o goadmin.exe ./app
   # Restart GoAdmin
   ```

   > Or if you're using the `pnpm dev` script, this will auto-rebuild.

4. **Verify in UI**: Navigate to **Plugins** page to see your plugin running!

## Creating a Plugin

### 1. Plugin Structure

```go
package myplugin

import (
    "fmt"
    "github.com/ethanburkett/goadmin/app/plugins"
)

type MyPlugin struct {
    ctx *plugins.PluginContext
}

func (p *MyPlugin) Metadata() plugins.PluginMetadata {
    return plugins.PluginMetadata{
        ID:          "my-plugin",
        Name:        "My Plugin",
        Version:     "1.0.0",
        Author:      "Your Name",
        Description: "Plugin description",
        Permissions: []string{"required.permission"},
    }
}

func (p *MyPlugin) Init(ctx *plugins.PluginContext) error {
    p.ctx = ctx
    // Initialize plugin
    return nil
}

func (p *MyPlugin) Start() error {
    // Subscribe to events, register commands
    fmt.Println("[MyPlugin] Started")
    return nil
}

func (p *MyPlugin) Stop() error {
    // Cleanup subscriptions
    fmt.Println("[MyPlugin] Stopped")
    return nil
}

func (p *MyPlugin) Reload() error {
    // Reload configuration
    fmt.Println("[MyPlugin] Reloaded")
    return nil
}

// Register the plugin on import
func init() {
    plugins.Registry.Register(&MyPlugin{})
}
```

### 2. Adding to GoAdmin

**Recommended: Use the Auto-Import Script**

The easiest way to add plugins is to use the built-in script that automatically discovers and imports all plugins:

```powershell
.\scripts\build_plugins.ps1
```

This script will:

- ✅ Scan the entire `plugins/` directory for plugins
- ✅ Automatically add imports to `app/main.go`
- ✅ Show plugin status (active/inactive)
- ✅ Sort and deduplicate imports

After running the script:

```bash
go build -o goadmin.exe ./app
```

**Manual Import (Alternative)**

If you prefer manual control:

1. Create your plugin anywhere (e.g., `plugins/myplugin/myplugin.go`)
2. Add the import to `app/main.go`:
   ```go
   import (
       // ... other imports
       _ "github.com/ethanburkett/goadmin/plugins/myplugin"
   )
   ```
3. Rebuild GoAdmin:
   ```bash
   go build -o goadmin.exe ./app
   ```

**Note**: Plugins can be located anywhere in the `plugins/` directory - the location doesn't matter, only the import path.

### 3. Managing Plugins

Use the plugin script to manage imports:

```powershell
.\scripts\build_plugins.ps1
```

**Output Example:**

```
Found 2 plugin(s):

  📦 auto-messages
     Status: + Will be imported (ACTIVATING)

  📦 examples/example
     Status: ✓ Already imported (ACTIVE)

Auto-importing 1 plugin(s)...
  ✓ Added: github.com/ethanburkett/goadmin/plugins/auto-messages
```

## Plugin Management

### REST API

- `GET /plugins` - List all plugins
- `GET /plugins/:id` - Get plugin status
- `POST /plugins/:id/start` - Start plugin
- `POST /plugins/:id/stop` - Stop plugin
- `POST /plugins/:id/reload` - Reload plugin config

### Web UI

Navigate to **Plugins** in the sidebar to:

- View all registered plugins
- See plugin status (Running, Stopped, Error)
- Start/Stop/Reload plugins
- View plugin metadata and permissions
- Monitor plugin dependencies

### Permissions

- `plugins.view` - View plugin list and status
- `plugins.manage` - Start, stop, reload plugins

## Example Plugins

### 1. Example Plugin (`plugins/examples/example/example.go`)

Demonstrates core functionality:

- **Event subscriptions**: Player connect/disconnect
- **Custom commands**: `!hello`, `!time`, `!echo <message>`
- **RCON integration**: Send welcome messages, reply to commands
- **Lifecycle management**: Proper Init, Start, Stop, Reload

### 2. Auto Messages Plugin (`plugins/auto-messages/auto-messages.go`)

Production-ready plugin that:

- **Periodic broadcasts**: Send messages every 5 minutes
- **Custom command**: `!nextmsg` to preview next message
- **RCON usage**: Uses `say` command for server-wide messages
- **Configurable messages**: Easy to modify message list

Both plugins showcase best practices for plugin development.

## Best Practices

1. Use `.\scripts\build_plugins.ps1` for auto-import
2. Cleanup in `Stop()` - unsubscribe events, unregister commands
3. Return errors from lifecycle methods
4. Call `plugins.Registry.Register()` in `init()`
5. Don't block in event handlers (they run in goroutines)
6. Request required permissions in metadata
7. Use semantic versioning
8. List dependencies in metadata
9. Use lowercase package names
10. Recreate channels/tickers in `Start()` for proper restart behavior

## Troubleshooting

### Plugin not appearing

1. Run `.\scripts\build_plugins.ps1`
2. Verify `init()` calls `plugins.Registry.Register()`
3. Rebuild: `go build -o goadmin.exe ./app`
4. Restart GoAdmin
5. Check logs for errors

### Plugin crashes on start

- Verify required permissions are granted
- Check dependencies are available
- Review `Init()` and `Start()` for errors
- Ensure EventBus/APIs are initialized

### Events not firing

- Verify subscription in `Start()`
- Check event type spelling (case-sensitive)
- Ensure plugin state is "Running"

## Architecture Notes

**Registry-Based Approach:**

- Plugins compile into the binary (single executable)
- Cross-platform compatible
- Requires rebuild to add/remove plugins
- Runtime lifecycle management (start/stop/reload)

## Future Enhancements

- Dependency validation
- Version compatibility checks
- Resource monitoring
- UI extension points
- Additional event types
- Configuration UI
- Plugin marketplace

## Security

- Plugins have full database and RCON access
- Run with GoAdmin's process permissions
- **Only enable trusted plugins** - review code before importing
- Use permission system to restrict capabilities
- Monitor resource usage in production
- Audit actions via audit log
