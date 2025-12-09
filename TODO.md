# GoAdmin TODO List

This document tracks major improvements and refactoring tasks for GoAdmin.

## ✅ COMPLETED - Database Schema & Normalization

### Schema Normalization

- [x] ✅ Audit all foreign key relationships and add missing constraints
  - Added `constraint:OnDelete:CASCADE` to Session.UserID
  - Added `constraint:OnDelete:SET NULL` to Report.ReviewedByUserID
  - Added `constraint:OnDelete:SET NULL` to TempBan.BannedByUser
  - Added `constraint:OnDelete:CASCADE` to CommandHistory.UserID
  - Added `constraint:OnDelete:SET NULL` to InGamePlayer.GroupID
- [x] ✅ Normalize command definitions table
  - [x] ✅ Separated command metadata with proper structure
  - [x] ✅ Replaced JSON permissions field with many-to-many relationship
  - [x] ✅ Added FK constraints to custom_commands→permissions via command_permissions junction table
  - [x] ✅ Added CASCADE constraints for referential integrity
  - [x] ✅ Created migration 007 for command permissions normalization
  - [x] ✅ Updated REST API to work with permission IDs
  - [x] ✅ Updated in-game command handler to use Permission objects
  - [x] ✅ Added helper methods: AddPermissionToCommand, RemovePermissionFromCommand, SetCommandPermissions, HasPermission
- [x] ✅ Normalize permission mappings
  - [x] ✅ Ensure all permission relationships have FK constraints
  - [x] ✅ Add cascading rules (ON DELETE CASCADE) to many-to-many relationships
  - [x] ✅ Updated Role model with constraint:OnDelete:CASCADE for user_roles and role_permissions
  - [x] ✅ Updated Permission model with constraint:OnDelete:CASCADE for role_permissions
  - [x] ✅ Updated User model with constraint:OnDelete:CASCADE for user_roles
- [x] ✅ Normalize role mappings
  - [x] ✅ Add FK constraints between users, roles, and permissions
  - [x] ✅ Add unique constraints where needed (already present via uniqueIndex on names)
  - [x] ✅ Ensure proper cascading behavior for role assignments
- [x] ✅ Server instances normalization
  - [x] ✅ Created Server model with proper configuration fields
  - [x] ✅ Linked TempBan, Report, CommandHistory, InGamePlayer, ServerStats to servers
  - [x] ✅ Added ServerID foreign keys with appropriate constraints (SET NULL or CASCADE)
  - [x] ✅ Created migration 008 for server instances
  - [x] ✅ Added server management methods: CreateServer, GetServerByID, GetDefaultServer, etc.
  - [x] ✅ Added auto-initialization of default server from config file
  - [x] ✅ Multi-server foundation ready for future expansion
  - [x] ✅ Created server management REST API (10 endpoints: CRUD, activate/deactivate, set default)
  - [x] ✅ Updated all data routes to accept optional server_id query parameter for filtering
  - [x] ✅ Created ServerProvider context for frontend multi-server management
  - [x] ✅ Created useServers hooks for all server CRUD operations
  - [x] ✅ Created ServerSelector component with dropdown navigation
  - [x] ✅ Restructured frontend with [id] folder for server-scoped routes
  - [x] ✅ Modified custom routing system to support layout.tsx files
  - [x] ✅ Created server management UI at /servers with full CRUD capabilities

### Database Integrity

- [x] ✅ Add database migration versioning system
  - [x] ✅ Migration model with version tracking
  - [x] ✅ Migration history tracking
  - [x] ✅ MigrationRunner with apply/rollback support
  - [x] ✅ Transaction-safe migrations
  - [x] ✅ REST API endpoints for migration management
  - [x] ✅ Frontend UI for migration management
- [x] ✅ Create database integrity validation script
  - [x] ✅ Checks for orphaned records (sessions, roles, permissions, reports, bans, etc.)
  - [x] ✅ Validates FK relationships across all tables
  - [x] ✅ Identifies missing indexes
  - [x] ✅ Reports constraint violations with severity levels (error, warning, info)
  - [x] ✅ JSON output for machine readability
  - [x] ✅ Summary statistics with exit codes
- [x] ✅ Add database backup/restore functionality
  - [x] ✅ Backup script with compression (zip format)
  - [x] ✅ Handles database file, WAL, and SHM files
  - [x] ✅ Automatic cleanup of old backups (keeps last 10)
  - [x] ✅ Restore script with validation
  - [x] ✅ Force flag to overwrite existing database
  - [x] ✅ PowerShell wrapper scripts for easy execution
- [x] ✅ Implement transaction safety for critical operations
  - [x] ✅ Role assignments (AddRoleToUser, RemoveRoleFromUser)
  - [x] ✅ Permission assignments (AddPermissionToRole, RemovePermissionFromRole)
  - [x] ✅ Group deletion (DeleteGroup with cascading player updates)
  - [x] ✅ Temp ban creation (CreateTempBan)
  - [x] ✅ All operations use DB.Transaction with automatic rollback on error
- [ ] Add database constraint violation handling

**Migration System Files:**

- `app/models/Migration.go` - Migration tracking models
- `app/database/migrations.go` - MigrationRunner implementation
- `app/main.go` - Migration definitions and system integration (7 migrations)
- `app/rest/migrations.go` - REST API endpoints
- `frontend/src/hooks/useMigrations.ts` - React hooks for migrations
- `frontend/src/pages/migrations.tsx` - Migration management UI

**Command Permissions Normalization Files:**

- `app/models/CustomCommand.go` - Updated to use many-to-many relationship with permissions
  - Replaced Permissions string (JSON) with []Permission slice
  - Added command_permissions junction table with CASCADE constraints
  - Added helper methods: AddPermissionToCommand, RemovePermissionFromCommand, SetCommandPermissions, HasPermission
  - Updated all query methods to Preload permissions
- `app/rest/commands.go` - Updated to convert permission names to IDs
  - CreateCommand converts permission names to IDs before saving
  - UpdateCommand uses SetCommandPermissions for atomic updates
  - Removed JSON marshaling/unmarshaling
- `app/commands/handler.go` - Updated in-game command permission checking
  - Updated to work with Permission objects instead of JSON strings
  - Removed hasRequiredPermissions function (JSON-based)
  - Direct permission name comparison with Permission.Name field
- `app/commands/admin.go` - Updated admin list permission checking
  - Works with Permission objects for command filtering
- `app/main.go` - Updated default command initialization
  - Added migration 007 for command permissions normalization
  - Converts permission names to IDs when creating default commands

**Server Instances Normalization Files:**

- `app/models/Server.go` - New server instances model (147 lines)
  - Stores server configuration: host, port, RCON password, games_mp.log path
  - Support for multiple servers with default server selection
  - Management methods: CreateServer, GetServerByID, GetDefaultServer, GetActiveServers
  - Helper methods: SetAsDefault, Activate, Deactivate
  - Relationships to TempBan, Report, CommandHistory, InGamePlayer, ServerStats
- `app/rest/servers.go` - Server management REST API (395 lines)
  - 10 endpoints: GET /servers, GET /servers/active, GET /servers/default, POST /servers, GET /servers/:id, PUT /servers/:id, DELETE /servers/:id, POST /servers/:id/default, POST /servers/:id/activate, POST /servers/:id/deactivate
  - Requires servers.manage permission
  - Full CRUD with validation and audit logging
- `app/models/TempBan.go` - Added ServerID field with FK constraint
  - Links temp bans to specific servers
  - Updated CreateTempBan to accept serverID parameter
  - Updated query methods with optional serverID filtering
- `app/models/Report.go` - Added ServerID field with FK constraint
  - Links reports to specific servers
  - Updated CreateReport to accept serverID parameter
  - Updated query methods with optional serverID filtering
- `app/models/CommandHistory.go` - Added ServerID field with FK constraint
  - Tracks which server commands were executed on
  - Updated CreateCommandHistory to accept serverID parameter
  - Updated query methods with optional serverID filtering
- `app/models/Group.go` (InGamePlayer) - Added ServerID field with FK constraint
  - Links in-game players to specific servers (for multi-server setups)
  - Updated GetAllInGamePlayers with optional serverID filtering
- `app/models/ServerStats.go` - Added ServerID field with FK constraint
  - Links server statistics to specific server instances
  - Updated GetServerStatsRange with optional serverID filtering
- `app/rest/reports.go` - Updated to support server_id query parameter
  - getAllReports, getPendingReports, getAllTempBans, getActiveTempBans accept optional server_id
- `app/rest/rcon.go` - Updated to support server_id query parameter
  - getServerStats, getCommandHistory accept optional server_id
- `app/rest/groups.go` - Updated to support server_id query parameter
  - getAllInGamePlayers accepts optional server_id
- `frontend/src/providers/ServerProvider.tsx` - Server context for multi-server management (125 lines)
  - Auto-detects server from URL params (:id)
  - Redirects to default server if none specified (unless disableRedirect prop is set)
  - Provides currentServer, servers list, switchServer, refreshServers
  - useServerContext hook for consuming context
- `frontend/src/hooks/useServers.ts` - React Query hooks for server management (185 lines)
  - useServers, useActiveServers, useDefaultServer, useServer
  - useCreateServer, useUpdateServer, useDeleteServer
  - useSetDefaultServer, useActivateServer, useDeactivateServer
  - All hooks use API generics pattern (api.get<Server[]>())
- `frontend/src/components/ServerSelector.tsx` - Server dropdown component (67 lines)
  - Dropdown menu with server list
  - Shows current server with checkmark
  - Positioned side="right" to stay within sidebar bounds
  - Uses useServerContext for current server and switching
- `frontend/src/components/DashboardLayout.tsx` - Updated with ServerSelector and server-aware navigation
  - Added ServerSelector to sidebar header
  - buildHref() adds server ID to all navigation paths (/:id/analytics, etc.)
  - isActive() checks routes accounting for server ID
- `frontend/src/pages/[id]/layout.tsx` - Layout wrapper for server-scoped routes
  - Wraps children in ServerProvider → DashboardLayout → padding wrapper
  - Provides centralized layout structure for all /:id routes
- `frontend/src/pages/servers.tsx` - Server management UI (700+ lines)
  - Full CRUD interface for server instances
  - Create/edit server dialogs with form validation
  - Delete, activate/deactivate, set default actions
  - Table view with server status badges
  - Sidebar with ServerSelector and user controls (no per-server navigation)
  - Uses ServerProvider with disableRedirect to prevent auto-navigation
- `frontend/startup/routes.ts` - Modified custom routing generator to support layouts
  - Detects layout.tsx files and creates parent/child route structure
  - Groups pages by directory for proper layout nesting
  - Handles dynamic route params ([id] → :id)
- `app/main.go` - Added migration 008 and initializeDefaultServer function
  - Auto-creates default server from config file on first run
  - Sets up server infrastructure for multi-server support
  - Registered servers.manage permission for super_admin role

**Database Integrity & Transaction Safety Files:**

- `scripts/validate_db.go` - Comprehensive database integrity validation script
  - Checks for orphaned sessions, roles, permissions, reports, bans, command history, players, webhooks, audit logs
  - Validates foreign key relationships
  - Identifies missing indexes
  - JSON output with severity levels and summary statistics
- `scripts/validate_db.ps1` - PowerShell wrapper for validation script
- `scripts/backup_db.go` - Database backup with compression and rotation
  - Compresses database, WAL, and SHM files to zip archive
  - Timestamp-based filenames (backup_YYYY-MM-DD_HH-MM-SS.zip)
  - Automatic cleanup keeps last 10 backups
- `scripts/backup_db.ps1` - PowerShell wrapper for backup script
- `scripts/restore_db.go` - Database restore from compressed backups
  - Extracts zip archive to database location
  - Validates backup file exists
  - Force flag to overwrite existing database
- `scripts/restore_db.ps1` - PowerShell wrapper for restore script
- `app/models/User.go` - Updated with transactional role assignment operations
- `app/models/Role.go` - Updated with CASCADE constraints and transactional permission operations
- `app/models/Permission.go` - Updated with CASCADE constraints
- `app/models/Group.go` - Updated with transactional delete operation
- `app/models/TempBan.go` - Updated with transactional ban creation
- `app/main.go` - Added migration 005 for permission constraint updates

**Performance Optimization Files:**

- `app/models/Report.go` - Added index to ReviewedByUserID
- `app/models/TempBan.go` - Added index to BannedByUser
- `app/models/Group.go` - Added index to InGamePlayer.GroupID
- `app/database/database.go` - Added connection pool configuration (25 max open, 10 max idle, 1hr lifetime)
- `app/main.go` - Added migration 006 for performance indexes

**Health & Monitoring Files:**

- `app/rest/health.go` - Health check endpoints
  - GET /health - Comprehensive health status with DB and RCON checks
  - GET /health/ready - Readiness probe for Kubernetes
  - GET /health/live - Liveness probe for Kubernetes
  - Connection pool statistics in health response
- `app/rest/main.go` - Registered health routes

## ✅ COMPLETED - Audit Logging System

### Core Audit Infrastructure

- [x] ✅ Create `audit_logs` table with proper schema
  - [x] ✅ Timestamp (with timezone)
  - [x] ✅ User ID (who performed action)
  - [x] ✅ Action type (enum: ban, kick, command, role_change, etc.)
  - [x] ✅ Target entity (player ID, user ID, command ID, etc.)
  - [x] ✅ Source (web_ui, in_game, api)
  - [x] ✅ IP address
  - [x] ✅ Metadata (JSON for additional context)
  - [x] ✅ Result (success/failure)
- [x] ✅ Created `AuditLog` model in `app/models/AuditLog.go`
- [x] ✅ Created audit helper functions in `app/rest/audit_helper.go`
- [x] ✅ Registered AuditLog model in database migrations

### Audit Event Types

- [x] ✅ Ban actions (temp/permanent)
  - [x] ✅ Who issued the ban
  - [x] ✅ Who was banned
  - [x] ✅ Duration and reason
  - [x] ✅ Source (web/in-game)
- [x] ✅ Kick actions
- [x] ✅ RCON command execution
  - [x] ✅ Raw command
  - [x] ✅ Arguments
  - [x] ✅ Result/output
- [x] ✅ Role/permission changes
  - [x] ✅ Role assignments/removals
  - [x] ✅ User approval/rejection
- [x] ✅ Group assignments
  - [x] ✅ Group creation with permissions/power
  - [x] ✅ Group updates with metadata tracking
  - [x] ✅ Group deletion
  - [x] ✅ Player-to-group assignments
  - [x] ✅ Player removal from groups
- [x] ✅ Custom command creation/modification/deletion
  - [x] ✅ Command creation with permissions
  - [x] ✅ Command updates with change tracking
  - [x] ✅ Command deletion with security checks
  - [x] ✅ Built-in command protection logging
- [x] ✅ User approval/rejection
- [x] ✅ Login/logout events
  - [x] ✅ Successful logins
  - [x] ✅ Failed login attempts
  - [x] ✅ Logout events
- [x] ✅ Report submissions and actions
- [x] ✅ Security violations
  - [x] ✅ Invalid command attempts
  - [x] ✅ Restricted command attempts
  - [x] ✅ Command injection attempts

### Audit UI & Reporting

- [x] ✅ Create audit log viewer in web dashboard
  - [x] ✅ Filter by user, action type, date range, source, success status
  - [x] ✅ Search functionality
  - [x] ✅ Export to CSV/JSON
  - [x] ✅ Pagination support
- [x] ✅ Audit log API endpoints
  - [x] ✅ GET /audit/logs with filters
  - [x] ✅ GET /audit/logs/recent
  - [x] ✅ GET /audit/logs/user/:userId
  - [x] ✅ GET /audit/logs/action/:action
- [ ] Real-time audit log streaming (optional WebSocket)
- [ ] Audit log retention policy configuration
- [ ] Audit log archiving system

**Audit Logging Implementation Files:**

- `app/rest/groups.go` - Group operation audit logging
  - Group creation with permissions/power metadata
  - Group updates with change tracking
  - Group deletion with power metadata
  - Player-to-group assignments with group name
  - Player removal from groups
  - Security violations for failed operations
- `app/rest/commands.go` - Command operation audit logging
  - Command creation with RCON command and permissions
  - Command updates with change tracking metadata
  - Command deletion with command details
  - Built-in command protection logging
  - Security violations for failed operations
- `app/models/Group.go` - Added `GetInGamePlayerByID` function for audit trail

## ✅ COMPLETED - Security & Rate Limiting

### RCON Command Security

- [x] ✅ Implement command sandboxing
  - [x] ✅ Validate command syntax before execution
  - [x] ✅ Block dangerous command patterns
  - [x] ✅ Disallow list system for commands (killserver, quit, plugins, etc.)
- [x] ✅ Command validation layer
  - [x] ✅ Argument type checking
  - [x] ✅ Argument sanitization
  - [x] ✅ Maximum argument length limits (500 chars)
  - [x] ✅ Maximum argument count limits (20 args)
- [x] ✅ Command execution limits
  - [x] ✅ Max concurrent executions (via rate limiting)
  - [x] ✅ Timeout for long-running commands (5s default, configurable, context-aware)
  - [x] ✅ Prevent command injection (blocked patterns, metacharacter filtering)

### Rate Limiting System

- [x] ✅ Global rate limiting
  - [x] ✅ Per-user rate limits
  - [x] ✅ Per-IP rate limits
  - [x] ✅ Per-endpoint rate limits
- [x] ✅ RCON-specific rate limiting
  - [x] ✅ Commands per minute per user (30/min with 10 burst)
  - [x] ✅ Commands per minute per server
  - [x] ✅ Custom command execution limits
- [x] ✅ Rate limit storage (in-memory with cleanup)
- [x] ✅ Rate limit exceeded handling
  - [x] ✅ Cooldown periods (token bucket refill)
  - [ ] Auto-ban for abuse
  - [ ] Alert admins of rate limit violations
- [x] ✅ Created `app/rest/rate_limiter.go` with token bucket implementation
- [x] ✅ Applied rate limiting to RCON endpoints (30 req/min)
- [x] ✅ Applied rate limiting to auth endpoints (5 req/min for login/register)

### Command Abuse Prevention

- [x] ✅ Detect spam patterns (via rate limiting)
  - [x] ✅ Token bucket algorithm prevents identical/similar commands in quick succession
- [x] ✅ Command deduplication
  - [x] ✅ Prevent duplicate command execution from CoD4's dual log entries (say/sayteam)
  - [x] ✅ 2-second deduplication window per player
- [x] ✅ Detect ban loops
  - [x] ✅ Prevent rapid ban/unban cycles (5 bans in 15 min threshold)
  - [x] ✅ Detect circular ban attempts (admin repeatedly banning same player)
  - [x] ✅ Track ban pattern statistics (suspicion scoring)
  - [x] ✅ Log security violations for ban loop abuse
- [x] ✅ Command throttling per target
  - [x] ✅ Prevent one admin from targeting same player too frequently (30s cooldown)
  - [x] ✅ Track target statistics per admin
- [ ] Emergency shutdown triggers
  - [ ] Auto-disable commands on abuse detection
  - [ ] Alert super admins

## 🟢 Medium Priority - Plugin/Extension System

### ✅ COMPLETED - Plugin Architecture Design

- [x] ✅ Design plugin interface/contract
  - [x] ✅ Define plugin lifecycle (init, start, stop, reload)
  - [x] ✅ Define plugin metadata structure (ID, Name, Version, Author, Description, Website, Dependencies, Permissions)
  - [x] ✅ Define plugin API surface (6 APIs: EventBus, Command, RCON, Database, Webhook, Config)
- [x] ✅ Create plugin loader system
  - [x] ✅ Go native plugin loading (.so files)
  - [x] ✅ Thread-safe plugin manager with lifecycle control
  - [x] ✅ Plugin state tracking (loaded, started, stopped, error)
  - [x] ✅ Context-aware cancellation for graceful shutdown
- [x] ✅ Plugin REST API
  - [x] ✅ GET /api/plugins - List all plugins
  - [x] ✅ GET /api/plugins/:id - Get plugin status
  - [x] ✅ POST /api/plugins/:id/start - Start plugin
  - [x] ✅ POST /api/plugins/:id/stop - Stop plugin
  - [x] ✅ POST /api/plugins/:id/reload - Reload plugin
- [x] ✅ Plugin management UI
  - [x] ✅ List installed plugins with status badges
  - [x] ✅ Start/stop/reload controls
  - [x] ✅ View plugin metadata (name, version, author, description)
  - [x] ✅ View plugin dependencies
  - [x] ✅ View plugin permissions
- [x] ✅ Example plugin implementation
  - [x] ✅ Event subscriptions (player connect/disconnect)
  - [x] ✅ Custom command registration (!hello)
  - [x] ✅ Configuration storage
  - [x] ✅ RCON command execution
  - [x] ✅ Webhook dispatching

**Files Created:**

- `app/plugins/plugin.go` (166 lines) - Plugin interface, PluginMetadata, PluginContext, API interfaces (EventBus, Command, RCON, Database, Webhook, Config)
- `app/plugins/manager.go` (265 lines) - Manager with LoadAll/StartAll/StopAll, individual Start/Stop/Reload, thread-safe with sync.RWMutex
- `app/rest/plugins.go` (120 lines) - REST API endpoints with permission checks (plugins.view, plugins.manage)
- `frontend/src/hooks/usePlugins.ts` (110 lines) - React hooks for plugin management
- `frontend/src/pages/plugins.tsx` (340 lines) - Plugin management UI with status display and controls
- `plugins/examples/example/example.go` (160 lines) - Example plugin demonstrating all APIs
- `plugins/examples/example/README.md` - Build and installation instructions

**Files Modified:**

- `app/rest/main.go` - RegisterPluginRoutes
- `frontend/routes.tsx` - Added plugins route
- `frontend/src/components/DashboardLayout.tsx` - Added Plugins navigation item

### Plugin Types & Capabilities (Future Enhancements)

- [ ] Hot-reload support (currently requires stop/start)
- [ ] Plugin dependency management (validation)
- [ ] Plugin versioning (compatibility checks)
- [ ] Plugin sandbox/isolation
  - [ ] Resource limits (CPU, memory)
  - [ ] API access controls beyond permissions
- [ ] Advanced command plugins
  - [ ] Command hooks/middleware
- [ ] Advanced event listener plugins
  - [ ] Kill/death events
  - [ ] Chat message events
  - [ ] Server state change events
- [ ] UI plugins
  - [ ] Custom dashboard widgets
  - [ ] Custom pages/routes
- [ ] Integration plugins
  - [ ] Discord webhooks (can use WebhookAPI)
  - [ ] Slack notifications
  - [ ] External API integrations

### ✅ COMPLETED - Webhook System

- [x] ✅ Create webhook configuration table
- [x] ✅ Webhook event triggers
  - [x] ✅ Player banned (in-game & web)
  - [x] ✅ Report submitted (in-game)
  - [x] ✅ Report actioned (web)
  - [ ] Server status change
- [x] ✅ Webhook delivery system
  - [x] ✅ Retry logic with exponential backoff
  - [x] ✅ Delivery status tracking
  - [x] ✅ Webhook HMAC SHA256 signing for security
- [x] ✅ Webhook REST API
  - [x] ✅ Create/edit/delete webhooks
  - [x] ✅ Test webhook delivery
  - [x] ✅ View delivery logs
- [x] ✅ Webhook management UI (frontend)
  - [x] ✅ Create/edit/delete webhooks
  - [x] ✅ Test webhook delivery
  - [x] ✅ View delivery logs

**Files Created:**

- `app/models/Webhook.go` (180 lines) - Webhook & WebhookDelivery models
- `app/webhook/dispatcher.go` (255 lines) - Dispatcher with retry logic
- `app/rest/webhooks.go` (276 lines) - REST API endpoints
- `frontend/src/hooks/useWebhooks.ts` (130 lines) - React hooks for webhook CRUD
- `frontend/src/pages/webhooks.tsx` (420 lines) - Webhook management UI

**Files Modified:**

- `app/main.go` - Migrations & retry worker startup
- `app/rest/main.go` - Route registration
- `app/commands/moderation.go` - Dispatch ban/report events
- `app/rest/reports.go` - Dispatch web ban/report events

### Event System

- [x] ✅ Core event dispatcher (webhook.GlobalDispatcher)
- [x] ✅ Standard event types (10 defined)
  - [x] ✅ player.banned, player.unbanned, player.kicked
  - [x] ✅ report.created, report.actioned
  - [x] ✅ user.approved, user.rejected
  - [x] ✅ server.online, server.offline (integrated in stats collector)
  - [x] ✅ security.alert
- [ ] Event middleware/filtering
- [ ] Event persistence (optional)
- [ ] Event replay capability

**Server Status Event Integration:**

- `app/watcher/stats.go` - Server status tracking and webhook dispatch
  - Added `lastOnline` field to track server state
  - Added webhook dispatcher to stats collector
  - Dispatches `server.online` when server comes online
  - Dispatches `server.offline` when server goes offline
  - Status changes detected during stat collection cycle

## ✅ COMPLETED - Permission System Refactoring

### Granular Permissions

- [x] ✅ Added specific permissions to replace generic `rbac.manage`
  - [x] ✅ `audit.view` - View audit logs
  - [x] ✅ `webhooks.manage` - Manage webhook configurations
  - [x] ✅ `migrations.manage` - Manage database migrations
  - [x] ✅ `groups.manage` - Manage in-game groups
  - [x] ✅ `commands.manage` - Manage custom commands
- [x] ✅ Updated all REST API routes to use specific permissions
- [x] ✅ Updated frontend sidebar navigation with granular permissions
- [x] ✅ Registered new permissions in super admin role initialization

**Updated Files:**

- `app/main.go` - Added 5 new permission definitions
- `app/rest/audit.go` - Changed to `audit.view`
- `app/rest/webhooks.go` - Changed to `webhooks.manage`
- `app/rest/migrations.go` - Changed to `migrations.manage`
- `app/rest/groups.go` - Changed to `groups.manage`
- `app/rest/commands.go` - Changed to `commands.manage`
- `frontend/src/components/DashboardLayout.tsx` - Updated sidebar permissions

## 🟠 High Priority - Security & Rate Limiting

### RCON Command Security

- [ ] Implement command sandboxing
  - [ ] Validate command syntax before execution
  - [ ] Block dangerous command patterns
  - [ ] Whitelist/blacklist system for commands
- [ ] Command validation layer
  - [ ] Argument type checking
  - [ ] Argument sanitization
  - [ ] Maximum argument length limits
- [ ] Command execution limits
  - [ ] Max concurrent executions
  - [ ] Timeout for long-running commands
  - [ ] Prevent command injection

### Advanced Permission System

- [x] ✅ Granular command permissions
  - [x] ✅ Per-command permission requirements (rcon.command, rcon.kick, rcon.ban, etc.)
  - [x] ✅ Command execution context (web vs in-game)
  - [x] ✅ Specific permissions for admin features (audit.view, webhooks.manage, etc.)
- [x] ✅ Permission audit trail
  - [x] ✅ Track permission grants/revokes via audit logs
  - [x] ✅ Track permission usage via audit logs

### Command Abuse Prevention

- [ ] Detect spam patterns
  - [ ] Identical commands in quick succession
  - [ ] Similar commands with minor variations
- [ ] Detect ban loops
  - [ ] Prevent rapid ban/unban cycles
  - [ ] Detect circular ban attempts
- [ ] Command throttling per target
  - [ ] Prevent one user from being targeted repeatedly
- [ ] Emergency shutdown triggers
  - [ ] Auto-disable commands on abuse detection
  - [ ] Alert super admins

## 🔵 Additional Improvements

### Performance

- [x] ✅ Add database query optimization
  - [x] ✅ Index analysis and optimization
    - [x] ✅ Added index to Report.ReviewedByUserID
    - [x] ✅ Added index to TempBan.BannedByUser
    - [x] ✅ Added index to InGamePlayer.GroupID
    - [x] ✅ Created migration 006 for performance indexes
  - [ ] Query caching for common operations (user sessions, role/permission lookups, server status)
  - [x] ✅ Connection pooling tuning
    - [x] ✅ Set MaxOpenConns to 25
    - [x] ✅ Set MaxIdleConns to 10
    - [x] ✅ Set ConnMaxLifetime to 1 hour
    - [x] ✅ Added connection pool metrics logging

### Testing

- [ ] Unit tests for core functionality (models, validators, rate limiters, ban loop detector, command throttler)
- [ ] Integration tests for RCON communication
- [ ] E2E tests for critical user flows
- [ ] Load testing for rate limiting
- [ ] Security testing for command injection

### Documentation

- [ ] Plugin development guide
- [ ] API documentation
- [ ] Security best practices
- [ ] Deployment guide
- [ ] Troubleshooting guide

### Monitoring & Observability

- [ ] Prometheus metrics export
- [x] ✅ Health check endpoints
  - [x] ✅ GET /health - Comprehensive health status with database and RCON checks
  - [x] ✅ GET /health/ready - Kubernetes readiness probe endpoint
  - [x] ✅ GET /health/live - Kubernetes liveness probe endpoint
  - [x] ✅ Connection pool statistics in health response
  - [x] ✅ Status codes: 200 (healthy), 503 (unhealthy/degraded)
- [ ] Performance monitoring
- [ ] Error tracking (Sentry integration?)
- [ ] Server metrics dashboard

---

## 📋 Implementation Summary (December 8, 2025)

### ✅ Phase 1: Foundation - COMPLETED

**Database Schema & Foreign Key Constraints**

- ✅ Added CASCADE constraint to `Session.UserID`
- ✅ Added SET NULL constraint to `Report.ReviewedByUserID`
- ✅ Added SET NULL constraint to `TempBan.BannedByUser`
- ✅ Added CASCADE constraint to `CommandHistory.UserID`
- ✅ Added SET NULL constraint to `InGamePlayer.GroupID`

**Files Modified:**

- `app/models/Session.go`
- `app/models/Report.go`
- `app/models/TempBan.go`
- `app/models/CommandHistory.go`
- `app/models/Group.go`

### ✅ Phase 2: Audit Logging System - COMPLETED

**Audit Logging System**

- ✅ Created comprehensive `AuditLog` model with 22 action types
- ✅ Implemented audit helper functions for common actions
- ✅ Added audit logging to ban/tempban/kick actions (web UI + in-game)
- ✅ Added audit logging to all RCON command executions
- ✅ Added audit logging to authentication events (login/logout/failures)
- ✅ Added audit logging to RBAC changes (role assignment/removal)
- ✅ Added audit logging to user approval/rejection
- ✅ Added audit logging to security violations
- ✅ Registered AuditLog in database migrations
- ✅ Created audit log API endpoints with filtering
- ✅ Created audit log viewer UI in dashboard

**Files Created:**

- `app/models/AuditLog.go` (194+ lines)
- `app/rest/audit_helper.go` (187+ lines)
- `app/rest/audit.go` (190+ lines)
- `frontend/src/hooks/useAudit.ts` (130+ lines)
- `frontend/src/pages/audit.tsx` (350+ lines)

**Files Modified:**

- `app/main.go` (added AuditLog to migrations, registered audit routes)
- `app/rest/reports.go` (added audit logging for bans)
- `app/commands/moderation.go` (added audit logging for in-game tempban)
- `app/rest/rcon.go` (added audit logging for RCON commands and security violations)
- `app/rest/auth.go` (added audit logging for authentication events)
- `app/rest/rbac.go` (added audit logging for RBAC changes)
- `frontend/src/components/DashboardLayout.tsx` (added audit logs nav item)
- `frontend/routes.tsx` (audit route auto-generated)

### ✅ Phase 3: Security & Rate Limiting - COMPLETED

**Rate Limiting Infrastructure**

- ✅ Implemented token bucket rate limiter with automatic cleanup
- ✅ Created global rate limiters: API (100/min), RCON (30/min), Login (5/min)
- ✅ Applied rate limiting to RCON endpoints
- ✅ Applied rate limiting to auth endpoints (login/register)

**Files Created:**

- `app/rest/rate_limiter.go` (170+ lines)

**Command Validation & Sandboxing**

- ✅ Created comprehensive RCON command validator
- ✅ Changed from allowlist to disallowlist (blocks: quit, killserver, plugins, devmap, etc.)
- ✅ Blocked patterns for dangerous operations (command injection, password exposure)
- ✅ Command sanitization (null bytes, whitespace, injection)
- ✅ Length and argument count limits (500 chars, 20 args max)
- ✅ Applied validation to all RCON command executions
- ✅ Security violations logged to audit trail

**Files Created:**

- `app/rest/command_validator.go` (125+ lines)

**Files Modified:**

- `app/rest/rcon.go` (integrated command validation, rate limiting, audit logging)
- `app/rest/auth.go` (added rate limiting middleware)

### ✅ Phase 4: Command Abuse Prevention - COMPLETED

**Command Deduplication**

- ✅ Prevents duplicate command execution from CoD4's dual log entries (say/sayteam)
- ✅ 2-second deduplication window per player
- ✅ Thread-safe with automatic cleanup

**Files Modified:**

- `app/commands/handler.go` (added deduplication logic, recent command tracking)

**Ban Loop Detection**

- ✅ Detects rapid ban/unban cycles (5 bans in 15 min threshold)
- ✅ Detects circular ban attempts (admin repeatedly banning same player)
- ✅ Tracks ban pattern statistics with suspicion scoring
- ✅ Logs security violations for ban loop abuse
- ✅ Provides detailed ban history and statistics

**Files Created:**

- `app/models/BanLoopDetector.go` (200+ lines)

**Files Modified:**

- `app/commands/moderation.go` (added ban loop detection)
- `app/rest/reports.go` (added ban loop detection to web UI tempban)

**Command Throttling**

- ✅ Prevents admins from targeting same player too frequently (30s cooldown)
- ✅ Tracks target statistics per admin
- ✅ Thread-safe with automatic cleanup

**Files Created:**

- `app/models/CommandThrottler.go` (105+ lines)

**Files Modified:**

- `app/commands/moderation.go` (added command throttling)

**Command Timeout Handling**

- ✅ Default 5-second timeout for RCON commands
- ✅ Configurable timeout via `SendCommandWithTimeout`
- ✅ Context-aware cancellation via `SendCommandWithContext`

**Files Modified:**

- `app/rcon/index.go` (added timeout methods)

### 🎯 Next Steps

**High Priority:**

- Create audit log viewer UI in dashboard
- Add audit logging for role/permission changes
- Add audit logging for user approval/rejection
- Add audit logging for login/logout events
- Implement ban loop detection
- Add emergency shutdown triggers for abuse

**Medium Priority:**

- Design plugin architecture
- Create webhook system
- Implement event bus/dispatcher

**Low Priority:**

- Database migration versioning
- Redis caching layer
- Comprehensive testing suite
- Performance monitoring

---

## 📊 Priority Matrix

| Priority    | Category                 | Estimated Effort | Status      |
| ----------- | ------------------------ | ---------------- | ----------- |
| 🔴 Critical | Schema Normalization     | 2-3 days         | ✅ Complete |
| 🔴 Critical | Audit Logging            | 3-4 days         | ✅ Complete |
| 🟠 High     | Rate Limiting & Security | 2-3 days         | ✅ Complete |
| 🟢 Medium   | Webhook System           | 2-3 days         | ✅ Complete |
| 🟢 Medium   | Plugin System (Basic)    | 5-7 days         | 📋 Planned  |
| 🔵 Low      | Additional Improvements  | Ongoing          | 📋 Planned  |

## Implementation Order

1. **✅ Phase 1: Foundation** (Week 1) - COMPLETED

   - ✅ Database schema normalization
   - ✅ Basic audit logging infrastructure
   - ✅ Critical security fixes

2. **✅ Phase 2: Security** (Week 2) - COMPLETED

   - ✅ Rate limiting system
   - ✅ Command sandboxing
   - ✅ Command validation

3. **✅ Phase 3: Extensibility** (Week 3-4) - COMPLETED

   - ✅ Webhook system with retry logic
   - ✅ Event dispatcher system
   - ✅ HMAC webhook signing

4. **📋 Phase 4: Polish** (Ongoing) - PLANNED
   - Testing
   - Documentation
   - Performance optimization
   - Monitoring

---

_Last Updated: December 8, 2025_
