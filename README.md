# GoAdmin

Modern web-based administration panel for Call of Duty 4 servers with RBAC and in-game command system.

## Screenshots

![Analytics](images/showcase-1.png)
![RCON Console](images/showcase-2.png)
![Custom Commands](images/showcase-3.png)

## Features

- Real-time server status, player management, analytics
- B3-style groups (Owner: 100, Admin: 50, VIP: 10)
- Custom RCON commands with placeholders (`{arg0}`, `{player}`, `{playerId:arg0}`, `{argsFrom:1}`)
- Report system with temp/permanent bans (5m, 2h, 3d, 1M, 2y)
- Plugin system with event subscriptions, custom commands, RCON access
- RBAC with granular permissions

Commands: `!groups`, `!mygroup`, `!putgroup`, `!adminlist`, `!help`, `!report`, `!tempban`, `!iamgod`

See [PLUGINS.md](PLUGINS.md) for plugin development.

## Tech Stack

**Backend:** Go, Gin, GORM, SQLite, RCON  
**Frontend:** React 19, TypeScript, TanStack Query v5, shadcn/ui, Tailwind, Vite

## Installation

**Prerequisites:** Go 1.21+, Node.js 18+, pnpm, CoD4 server with RCON

```bash
git clone https://github.com/ethanburkett/GoAdmin.git
cd GoAdmin
cp config.example.json config.json
```

Edit `config.json` with server details, then:

```bash
pnpm run deps
pnpm dev
```

Dashboard: `http://localhost:5173`

**First-time:** Register → `http://localhost:8080/auth/iamgod` → `!iamgod` in-game

## Development

```bash
# Build
go build -o goadmin app/main.go
cd frontend && pnpm build

# Plugins
.\scripts\build_plugins.ps1
go build -o goadmin app/main.go
```

## License

MIT License - Inspired by B3 (BigBrotherBot)
