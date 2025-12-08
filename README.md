# GoAdmin

A modern web-based administration panel for Call of Duty 4 (CoD4) game servers with advanced player management, RBAC (Role-Based Access Control), and in-game command system.

## Screenshots

![Analytics](images/showcase-1.png)
![RCON Console](images/showcase-2.png)
![Custom Commands](images/showcase-3.png)

## Features

### Web Dashboard

- **Modern React UI** - Built with React 19, TypeScript, and shadcn/ui components
- **Real-time Server Status** - Live player counts, map info, and server statistics
- **Player Management** - View online players, track statistics, and manage player data
- **RBAC System** - Granular permission control with roles and permissions
- **User Approval System** - Admin-approved user registration

### In-Game Administration

- **B3-Style Groups** - Power-based group hierarchy (Owner: 100, Admin: 50, VIP: 10)
- **Custom Commands** - Create custom RCON commands with placeholder support
  - `{arg0}`, `{arg1}` - Command arguments
  - `{player}`, `{guid}` - Player information
  - `{playerId:arg0}` - Auto-resolve player names to entity IDs
  - `{argsFrom:1}` - Join remaining arguments (for reasons, messages)
- **Built-in Commands** - Pre-configured admin commands
  - `!groups` - List all available groups
  - `!mygroup` - Show your current group and permissions
  - `!putgroup <player> <group>` - Assign players to groups
  - `!adminlist` - List online admins
  - `!help [page]` - Paginated command help
  - `!report <player> <reason>` - Report players for admin review
  - `!tempban <player> <duration> <reason>` - Temporarily ban players (5m, 2h, 3d, 1M, 2y)
  - `!iamgod` - First-use only: claim Owner privileges

### Player Moderation

- **Report System** - Players can report others in-game
- **Temporary Bans** - Time-based bans with automatic expiration
  - Auto-kick banned players on join
  - Duration formats: minutes (m), hours (h), days (d), months (M), years (y)
- **Action Dashboard** - Review and action reports from web UI
  - Dismiss reports
  - Issue temporary bans
  - Issue permanent bans

### RCON Integration

- **Direct RCON Access** - Execute any RCON command from the dashboard
- **Permission-Based Commands** - Secure command execution with role checking
- **Quick Actions** - Pre-configured buttons for common tasks
  - Kick players
  - Ban players
  - Change maps
  - Restart map
  - Send messages
  - Manage game settings

### Analytics

- **Server Statistics** - Track player counts, map playtime, and server uptime
- **Player Analytics** - Historical player data and trends
- **Command History** - Audit trail of all executed commands

## Tech Stack

### Backend

- **Go 1.x** - High-performance backend
- **Gin** - HTTP web framework
- **GORM** - ORM for SQLite
- **SQLite** - Embedded database
- **RCON** - Direct server communication

### Frontend

- **React 19** - Modern UI library
- **TypeScript** - Type-safe development
- **TanStack Query v5** - Data fetching and caching
- **TanStack Router** - Type-safe routing
- **shadcn/ui** - Beautiful UI components
- **Tailwind CSS** - Utility-first styling
- **Vite** - Fast build tool

## Installation

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and pnpm
- CoD4 server with RCON enabled

### Backend Setup

1. Clone the repository:

```bash
git clone https://github.com/ethanburkett/GoAdmin.git
cd GoAdmin
```

2. Configure the application:

```bash
cp config.example.json config.json
```

Edit `config.json`:

```json
{
  "server": {
    "host": "localhost",
    "port": 28960,
    "rcon_password": "your_rcon_password_here"
  },
  "games_mp_path": "...\\Call of Duty 4\\Mods\\your_mod\\games_mp.log",
  "rest_port": 8080,
  "environment": "development | production"
}
```

3. Install dependencies:

```bash
pnpm run deps
```

3b. Then run:

```bash
pnpm dev
```

The dashboard will be available at `http://localhost:5173`

## Configuration

### RCON Setup

Ensure your CoD4 server has RCON enabled in `server.cfg`:

```
set rcon_password "your_secure_password"
```

### Log File Watching

GoAdmin monitors the `games_mp.log` file for real-time events. Ensure the path in `config.json` points to your server's log file.

### First-Time Setup

1. Start both backend and frontend
2. Navigate to `http://localhost:5173`
3. Register your account
4. Head to `http://localhost:8080/auth/iamgod` in a separate tab, or if you changed the rest port, use that port. This will allow the first user to claim Owner privileges for the dashboard.
5. The first user should use `!iamgod` in-game to claim Owner privileges
6. Approve other users through the RBAC management panel

## Usage

### Creating Custom Commands

1. Navigate to Commands in the dashboard
2. Click "Create Command"
3. Configure:
   - **Name**: Command name (used as `!commandname`)
   - **Usage**: Help text shown to users
   - **Description**: Detailed explanation
   - **RCON Command**: Template with placeholders
   - **Min Power**: Minimum group power required
   - **Permissions**: Required permission array
   - **Requirement Type**: `power`, `permission`, or `both`

Example - Custom say command:

- Name: `announce`
- RCON: `say ^1[SERVER] ^7{argsFrom:0}`
- Min Power: 50
- Permissions: `["say"]`

### Managing Groups

Groups control in-game player privileges:

- **Owner (100)**: Full server control
- **Admin (50)**: Moderation and game management
- **VIP (10)**: Basic privileges

Assign players to groups using `!putgroup <player> <group>` in-game or through the dashboard.

### Reviewing Reports

1. Navigate to Reports in dashboard
2. Review pending reports
3. Take action:
   - **Dismiss**: Close the report
   - **Temp Ban**: Issue time-limited ban
   - **Permanent Ban**: Issue permanent RCON ban

## Development

### Building for Production

Backend:

```bash
go build -o goadmin app/main.go
```

Frontend:

```bash
cd frontend
pnpm build
```

## Security

- **Password Hashing**: bcrypt with salt rounds
- **Session Tokens**: Secure random tokens with expiration
- **RBAC**: Fine-grained permission system
- **User Approval**: Admin-approved registration
- **Command Validation**: Power and permission checks

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - See LICENSE file for details

## Support

For issues, questions, or feature requests, please open an issue on GitHub.

## Credits

Built with inspiration from B3 (BigBrotherBot) admin system for CoD4 servers.
