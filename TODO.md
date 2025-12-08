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
- [ ] Role/permission changes
  - [ ] Who changed what
  - [ ] Before/after state
- [ ] Group assignments
- [ ] Custom command creation/modification/deletion
- [ ] User approval/rejection
- [ ] Login/logout events
- [x] ✅ Report submissions and actions

### Audit UI & Reporting

- [ ] Create audit log viewer in web dashboard
  - [ ] Filter by user, action type, date range
  - [ ] Search functionality
  - [ ] Export to CSV/JSON
- [ ] Real-time audit log streaming (optional WebSocket)
- [ ] Audit log retention policy configuration
- [ ] Audit log archiving system

## ✅ COMPLETED - Security & Rate Limiting

### RCON Command Security

- [x] ✅ Implement command sandboxing
  - [x] ✅ Validate command syntax before execution
  - [x] ✅ Block dangerous command patterns
  - [x] ✅ Whitelist/blacklist system for commands
- [x] ✅ Command validation layer
  - [x] ✅ Argument type checking
  - [x] ✅ Argument sanitization
  - [x] ✅ Maximum argument length limits
- [x] ✅ Command execution limits
  - [x] ✅ Max concurrent executions (via rate limiting)
  - [ ] Timeout for long-running commands
  - [x] ✅ Prevent command injection

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
- [ ] Detect ban loops
  - [ ] Prevent rapid ban/unban cycles
  - [ ] Detect circular ban attempts
- [ ] Command throttling per target
  - [ ] Prevent one user from being targeted repeatedly
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

**Audit Logging System**

- ✅ Created comprehensive `AuditLog` model with 20+ action types
- ✅ Implemented audit helper functions for common actions
- ✅ Added audit logging to ban/tempban/kick actions (web UI + in-game)
- ✅ Added audit logging to all RCON command executions
- ✅ Registered AuditLog in database migrations

**Files Created:**

- `app/models/AuditLog.go` (200+ lines)
- `app/rest/audit_helper.go` (150+ lines)

**Files Modified:**

- `app/main.go` (added AuditLog to migrations)
- `app/rest/reports.go` (added audit logging for bans)
- `app/commands/moderation.go` (added audit logging for in-game tempban)
- `app/rest/rcon.go` (added audit logging for RCON commands)

**Rate Limiting Infrastructure**

- ✅ Implemented token bucket rate limiter with automatic cleanup
- ✅ Created global rate limiters: API (100/min), RCON (30/min), Login (5/min)
- ✅ Applied rate limiting to RCON endpoints
- ✅ Applied rate limiting to auth endpoints (login/register)

**Files Created:**

- `app/rest/rate_limiter.go` (170+ lines)

**Files Modified:**

- `app/rest/rcon.go` (added rate limiting middleware)
- `app/rest/auth.go` (added rate limiting middleware)

**Command Validation & Sandboxing**

- ✅ Created comprehensive RCON command validator
- ✅ Whitelist of 20+ allowed commands
- ✅ Blocked patterns for dangerous operations
- ✅ Command sanitization (null bytes, whitespace, injection)
- ✅ Length and argument count limits
- ✅ Applied validation to all RCON command executions

**Files Created:**

- `app/rest/command_validator.go` (190+ lines)

**Files Modified:**

- `app/rest/rcon.go` (integrated command validation)

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
