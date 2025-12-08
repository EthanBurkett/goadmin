# GoAdmin TODO List

This document tracks major improvements and refactoring tasks for GoAdmin.

## 🔴 Critical - Database Schema & Normalization

### Schema Normalization

- [ ] Audit all foreign key relationships and add missing constraints
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

## 🟡 High Priority - Audit Logging System

### Core Audit Infrastructure

- [ ] Create `audit_logs` table with proper schema
  - [ ] Timestamp (with timezone)
  - [ ] User ID (who performed action)
  - [ ] Action type (enum: ban, kick, command, role_change, etc.)
  - [ ] Target entity (player ID, user ID, command ID, etc.)
  - [ ] Source (web_ui, in_game, api)
  - [ ] IP address
  - [ ] Metadata (JSON for additional context)
  - [ ] Result (success/failure)

### Audit Event Types

- [ ] Ban actions (temp/permanent)
  - [ ] Who issued the ban
  - [ ] Who was banned
  - [ ] Duration and reason
  - [ ] Source (web/in-game)
- [ ] Kick actions
- [ ] RCON command execution
  - [ ] Raw command
  - [ ] Arguments
  - [ ] Result/output
- [ ] Role/permission changes
  - [ ] Who changed what
  - [ ] Before/after state
- [ ] Group assignments
- [ ] Custom command creation/modification/deletion
- [ ] User approval/rejection
- [ ] Login/logout events
- [ ] Report submissions and actions

### Audit UI & Reporting

- [ ] Create audit log viewer in web dashboard
  - [ ] Filter by user, action type, date range
  - [ ] Search functionality
  - [ ] Export to CSV/JSON
- [ ] Real-time audit log streaming (optional WebSocket)
- [ ] Audit log retention policy configuration
- [ ] Audit log archiving system

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

## 📊 Priority Matrix

| Priority    | Category                 | Estimated Effort |
| ----------- | ------------------------ | ---------------- |
| 🔴 Critical | Schema Normalization     | 2-3 days         |
| 🔴 Critical | Audit Logging            | 3-4 days         |
| 🟠 High     | Rate Limiting & Security | 2-3 days         |
| 🟢 Medium   | Plugin System (Basic)    | 5-7 days         |
| 🟢 Medium   | Webhook System           | 2-3 days         |
| 🔵 Low      | Additional Improvements  | Ongoing          |

## Implementation Order Recommendation

1. **Phase 1: Foundation** (Week 1)

   - Database schema normalization
   - Basic audit logging infrastructure
   - Critical security fixes

2. **Phase 2: Security** (Week 2)

   - Rate limiting system
   - Command sandboxing
   - Advanced permissions

3. **Phase 3: Extensibility** (Week 3-4)

   - Event system
   - Webhook system
   - Basic plugin architecture

4. **Phase 4: Polish** (Ongoing)
   - Testing
   - Documentation
   - Performance optimization
   - Monitoring

---

_Last Updated: December 8, 2025_
