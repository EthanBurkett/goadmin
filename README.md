<div align="center">

# 🎮 GoAdmin

### Modern Administration Platform for Call of Duty 4 Servers

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?logo=typescript)](https://www.typescriptlang.org/)

**A powerful, enterprise-grade web panel for CoD4 servers featuring real-time monitoring, advanced player management, and extensible plugin architecture.**

[Features](#-features) • [Screenshots](#-screenshots) • [Installation](#-installation) • [Documentation](#-documentation) • [Plugins](PLUGINS.md)

</div>

---

## 📸 Screenshots

<div align="center">

### Analytics Dashboard

![Analytics Dashboard](images/showcase-1.png)

### Real-time RCON Console

![RCON Console](images/showcase-2.png)

### Custom Command Builder

![Custom Commands](images/showcase-3.png)

</div>

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎯 Server Management

- **Real-time Monitoring** - Live player counts, server stats, map rotation
- **Multi-Server Support** - Manage multiple CoD4 instances from one dashboard
- **RCON Integration** - Direct server control with permission-based access
- **Log Monitoring** - Real-time games_mp.log parsing and event processing
- **Analytics Dashboard** - Player trends, server uptime, command history

</td>
<td width="50%">

### 👥 Player Management

- **Live Player View** - See who's online with real-time updates
- **Player Statistics** - Track performance, playtime, and history
- **Report System** - In-game player reporting with action dashboard
- **Ban Management** - Temporary and permanent bans with auto-expiration
- **Advanced Search** - Filter and find players by GUID, name, or stats

</td>
</tr>
<tr>
<td width="50%">

### 🛡️ Administration

- **RBAC System** - Role-Based Access Control with granular permissions
- **B3-Style Groups** - Power-based hierarchy (Owner: 100, Admin: 50, VIP: 10)
- **User Approval** - Admin-approved registration system
- **Audit Logging** - Complete trail of all administrative actions
- **Webhook Integration** - External notifications for key events

</td>
<td width="50%">

### ⚡ In-Game Commands

- **Custom Commands** - Build RCON commands with dynamic placeholders
- **Smart Placeholders** - `{arg0}`, `{player}`, `{playerId:arg0}`, `{argsFrom:1}`
- **Permission Checks** - Power level and permission-based validation
- **10+ Built-in Commands** - Ready-to-use admin commands
- **Plugin Commands** - Extend with custom Go-based commands

</td>
</tr>
<tr>
<td width="50%">

### 🔌 Plugin System

- **Event-Driven Architecture** - Subscribe to player events
- **Custom Commands** - Register in-game commands with Go callbacks
- **Full RCON Access** - Execute server commands from plugins
- **Hot Reload** - Start, stop, reload without restart
- **Auto-Import Script** - Automatic plugin discovery and activation

</td>
<td width="50%">

### 🎨 Modern UI/UX

- **React 19** - Latest React with concurrent features
- **shadcn/ui** - Beautiful, accessible component library
- **Dark Mode** - Built-in theme support
- **Responsive Design** - Works on desktop, tablet, mobile
- **Real-time Updates** - TanStack Query for seamless data sync

</td>
</tr>
</table>

---

## 🎮 Built-in Commands

<details>
<summary><b>Click to expand command list</b></summary>

| Command      | Description                             | Example                           |
| ------------ | --------------------------------------- | --------------------------------- |
| `!groups`    | List all available groups               | `!groups`                         |
| `!mygroup`   | Show your current group and permissions | `!mygroup`                        |
| `!putgroup`  | Assign player to a group                | `!putgroup Player1 admin`         |
| `!adminlist` | List all online administrators          | `!adminlist`                      |
| `!help`      | Show paginated help menu                | `!help 2`                         |
| `!report`    | Report a player for admin review        | `!report Player1 cheating`        |
| `!tempban`   | Issue temporary ban                     | `!tempban Player1 2h teamkilling` |
| `!iamgod`    | Claim Owner privileges (first use only) | `!iamgod`                         |

**Ban Duration Formats:** `5m` (minutes), `2h` (hours), `3d` (days), `1M` (months), `2y` (years)

</details>

---

## 🛠️ Tech Stack

<div align="center">

### Backend

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white)

### Frontend

![React](https://img.shields.io/badge/React_19-61DAFB?style=for-the-badge&logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![TanStack Query](https://img.shields.io/badge/TanStack_Query-FF4154?style=for-the-badge&logo=react-query&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-06B6D4?style=for-the-badge&logo=tailwind-css&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-646CFF?style=for-the-badge&logo=vite&logoColor=white)

</div>

---

## 🚀 Installation

### Prerequisites

<table>
<tr>
<td>

**Required:**

- Go 1.21 or higher
- Node.js 18+ with pnpm
- CoD4 server with RCON enabled

</td>
<td>

**Recommended:**

- Linux/Windows/macOS
- 2GB+ RAM
- SSD for database

</td>
</tr>
</table>

### Quick Start

```bash
# Clone repository
git clone https://github.com/ethanburkett/GoAdmin.git
cd GoAdmin

# Configure
cp config.example.json config.json
```

**Edit `config.json`:**

```json
{
  "server": {
    "host": "localhost", // CoD4 server IP
    "port": 28960, // CoD4 server port
    "rcon_password": "your_rcon_password"
  },
  "games_mp_path": "path/to/games_mp.log", // Log file location
  "rest_port": 8080, // API port
  "environment": "development" // development | production
}
```

**Launch:**

```bash
# Install dependencies and start dev servers
pnpm run deps
pnpm dev
```

🌐 **Dashboard:** `http://localhost:5173`  
🔌 **API:** `http://localhost:8080`

### First-Time Setup

1. **Register Account**  
   Navigate to `http://localhost:5173` and create your account

2. **Claim Dashboard Owner**  
   Visit `http://localhost:8080/auth/iamgod` to grant yourself Owner privileges

3. **Claim In-Game Admin**  
   Type `!iamgod` in-game to receive Owner group (100 power)

4. **Approve Users**  
   Manage user access via Dashboard → RBAC

---

## 📚 Documentation

### Custom Command Placeholders

Build dynamic RCON commands with these placeholders:

| Placeholder             | Description                           | Example                             |
| ----------------------- | ------------------------------------- | ----------------------------------- |
| `{arg0}`, `{arg1}`, ... | Command arguments                     | `{arg0}` → first argument           |
| `{player}`              | Player name                           | `{player}` → "Player1"              |
| `{guid}`                | Player GUID                           | `{guid}` → "abc123..."              |
| `{playerId:arg0}`       | Auto-resolve player name to entity ID | `{playerId:arg0}` → "5"             |
| `{argsFrom:1}`          | Join all args from index              | `{argsFrom:1}` → "reason text here" |

**Example Command:**

```
Command: !announce
RCON: say ^1[ADMIN] ^7{argsFrom:0}
Min Power: 50
Result: !announce Server restarting → say ^1[ADMIN] ^7Server restarting
```

### Power Groups

GoAdmin uses B3-style power-based groups:

```
┌─────────────────────────────────┐
│ OWNER (100)  - Full Control     │
├─────────────────────────────────┤
│ ADMIN (50)   - Moderation       │
├─────────────────────────────────┤
│ VIP (10)     - Basic Privileges │
├─────────────────────────────────┤
│ USER (0)     - No Permissions   │
└─────────────────────────────────┘
```

Assign players with `!putgroup <player> <group>` or via dashboard.

---

## 🔌 Plugin Development

GoAdmin features a powerful plugin system for extending functionality.

**Quick Start:**

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
        Author:  "Your Name",
    }
}

func (p *MyPlugin) Init(ctx *plugins.PluginContext) error {
    p.ctx = ctx
    return nil
}

func (p *MyPlugin) Start() error {
    // Subscribe to events
    p.ctx.EventBus.Subscribe("player.connect", func(data interface{}) {
        // Handle player join
    })

    // Register custom command
    p.ctx.CommandAPI.RegisterCommand(plugins.CommandDefinition{
        Name: "hello",
        Handler: func(player, guid string, args []string) error {
            p.ctx.RCONAPI.SendCommand(`say "Hello World!"`)
            return nil
        },
    })

    return nil
}

func (p *MyPlugin) Stop() error { return nil }
func (p *MyPlugin) Reload() error { return nil }

func init() {
    plugins.Registry.Register(&MyPlugin{})
}
```

**Auto-Import and Build:**

```bash
.\scripts\build_plugins.ps1
go build -o goadmin app/main.go
```

📖 **Full Documentation:** [PLUGINS.md](PLUGINS.md)

---

## 🏗️ Development

### Project Structure

```
GoAdmin/
├── app/
│   ├── commands/        # In-game command handlers
│   ├── config/          # Configuration management
│   ├── database/        # Database models and migrations
│   ├── logger/          # Logging utilities
│   ├── models/          # Data models
│   ├── parser/          # Log file parser
│   ├── plugins/         # Plugin system core
│   ├── rcon/            # RCON client
│   ├── rest/            # REST API endpoints
│   ├── watcher/         # Log file watcher
│   └── webhook/         # Webhook dispatcher
├── frontend/
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── hooks/       # Custom React hooks
│   │   ├── pages/       # Route pages
│   │   ├── providers/   # Context providers
│   │   └── lib/         # Utilities
│   └── public/          # Static assets
├── plugins/
│   └── examples/        # Example plugins
└── scripts/             # Build and utility scripts
```

### Build for Production

```bash
# Backend
go build -o goadmin app/main.go

# Frontend
cd frontend
pnpm build
```

### Running Tests

```bash
# Backend tests
go test ./...

# Frontend tests
cd frontend
pnpm test
```

---

## 🔒 Security

GoAdmin implements enterprise-grade security practices:

- 🔐 **bcrypt Password Hashing** - Industry-standard password security
- 🎟️ **Session Tokens** - Secure, expiring authentication tokens
- 🛡️ **RBAC** - Role-Based Access Control with granular permissions
- ✅ **User Approval** - Admin-approved registration system
- 🔍 **Command Validation** - Power level and permission checks
- 📝 **Audit Logging** - Complete action history tracking
- 🚫 **SQL Injection Protection** - Parameterized queries via GORM
- 🔒 **XSS Prevention** - React automatic escaping

---

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines

- Follow Go and TypeScript best practices
- Write tests for new features
- Update documentation
- Keep commits atomic and descriptive

---

## 📝 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 💝 Acknowledgments

- **B3 (BigBrotherBot)** - Inspiration for the command and group system
- **CoD4 Community** - Continued support and server hosting
- **Contributors** - Everyone who has helped improve GoAdmin

---

<div align="center">

### ⭐ Star us on GitHub!

Made with ❤️ for the Call of Duty 4 community

[Report Bug](https://github.com/ethanburkett/GoAdmin/issues) • [Request Feature](https://github.com/ethanburkett/GoAdmin/issues) • [Documentation](PLUGINS.md)

</div>
