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
- [ ] Normalize command definitions table
  - [ ] Separate command metadata from execution logic
  - [ ] Add proper FK constraints to roles/permissions
- [ ] Normalize permission mappings
  - [ ] Ensure all permission relationships have FK constraints
  - [ ] Add cascading rules (ON DELETE CASCADE/RESTRICT)
- [ ] Normalize role mappings
  - [ ] Add FK constraints between users, roles, and permissions
  - [ ] Add unique constraints where needed
- [ ] Server instances normalization
  - [ ] Create proper server configuration table
  - [ ] Link commands/groups/bans to specific server instances
  - [ ] Add server-level isolation for multi-server setups

### Database Integrity

- [ ] Add database migration versioning system
- [ ] Create database integrity validation script
- [ ] Add database backup/restore functionality
- [ ] Implement transaction safety for critical operations
- [ ] Add database constraint violation handling

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
- [ ] Group assignments
- [ ] Custom command creation/modification/deletion
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

### Plugin Architecture Design

- [ ] Design plugin interface/contract
  - [ ] Define plugin lifecycle (init, start, stop, reload)
  - [ ] Define plugin metadata structure
  - [ ] Define plugin API surface
- [ ] Create plugin loader system
  - [ ] Hot-reload support
  - [ ] Plugin dependency management
  - [ ] Plugin versioning
- [ ] Plugin sandbox/isolation
  - [ ] Resource limits (CPU, memory)
  - [ ] Permission system for plugins
  - [ ] API access controls

### Plugin Types & Capabilities

- [ ] Command plugins
  - [ ] Custom in-game commands
  - [ ] Command hooks/middleware
- [ ] Event listener plugins
  - [ ] Player join/leave events
  - [ ] Kill/death events
  - [ ] Chat message events
  - [ ] Server state change events
- [ ] UI plugins
  - [ ] Custom dashboard widgets
  - [ ] Custom pages/routes
- [ ] Integration plugins
  - [ ] Discord webhooks
  - [ ] Slack notifications
  - [ ] External API integrations

### Webhook System

- [ ] Create webhook configuration table
- [ ] Webhook event triggers
  - [ ] Player banned
  - [ ] Report submitted
  - [ ] Admin action taken
  - [ ] Server status change
- [ ] Webhook delivery system
  - [ ] Retry logic with exponential backoff
  - [ ] Delivery status tracking
  - [ ] Webhook signing for security
- [ ] Webhook management UI
  - [ ] Create/edit/delete webhooks
  - [ ] Test webhook delivery
  - [ ] View delivery logs

### Event System

- [ ] Create core event bus/dispatcher
- [ ] Define standard event types
- [ ] Event middleware/filtering
- [ ] Event persistence (optional)
- [ ] Event replay capability

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

### Rate Limiting System

- [ ] Global rate limiting
  - [ ] Per-user rate limits
  - [ ] Per-IP rate limits
  - [ ] Per-endpoint rate limits
- [ ] RCON-specific rate limiting
  - [ ] Commands per minute per user
  - [ ] Commands per minute per server
  - [ ] Custom command execution limits
- [ ] Rate limit storage (Redis recommended)
- [ ] Rate limit exceeded handling
  - [ ] Cooldown periods
  - [ ] Auto-ban for abuse
  - [ ] Alert admins of rate limit violations

### Advanced Permission System

- [ ] Granular command permissions
  - [ ] Per-command permission requirements
  - [ ] Command execution context (web vs in-game)
  - [ ] Time-based permissions (only during certain hours)
- [ ] Permission inheritance
  - [ ] Role hierarchy
  - [ ] Permission delegation
- [ ] Temporary permissions
  - [ ] Time-limited admin access
  - [ ] Scheduled permission changes
- [ ] Permission audit trail
  - [ ] Track permission grants/revokes
  - [ ] Track permission usage

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

- [ ] Add database query optimization
  - [ ] Index analysis and optimization
  - [ ] Query caching for common operations
  - [ ] Connection pooling tuning
- [ ] Add Redis caching layer
  - [ ] Cache user sessions
  - [ ] Cache role/permission lookups
  - [ ] Cache server status
- [ ] Background job processing
  - [ ] Async ban processing
  - [ ] Batch operations
  - [ ] Scheduled tasks (temp ban expiry, cleanup)

### Testing

- [ ] Unit tests for all core functionality
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
- [ ] Health check endpoints
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
| 🟢 Medium   | Plugin System (Basic)    | 5-7 days         | 📋 Planned  |
| 🟢 Medium   | Webhook System           | 2-3 days         | 📋 Planned  |
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

3. **📋 Phase 3: Extensibility** (Week 3-4) - PLANNED

   - Event system
   - Webhook system
   - Basic plugin architecture

4. **📋 Phase 4: Polish** (Ongoing) - PLANNED
   - Testing
   - Documentation
   - Performance optimization
   - Monitoring

---

_Last Updated: December 8, 2025_
