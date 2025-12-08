// Call of Duty 4 RCON Commands
// Comprehensive list of available server commands

export const cod4Commands = [
  // Core Admin & Server Management
  {
    command: "login <password>",
    description: "Authenticate with RCON password",
  },
  { command: "say <message>", description: "Send message to all players" },
  {
    command: "serverinfo",
    description: "Display server settings/configuration",
  },
  { command: "systeminfo", description: "Display server machine/system info" },
  {
    command: "status",
    description: "List connected players with ID, ping, GUID, IP",
  },
  { command: "exec <filename>", description: "Execute a server config file" },
  {
    command: "writeconfig <filename>",
    description: "Save server configuration to file",
  },
  { command: "quit", description: "Shutdown server gracefully" },
  { command: "killserver", description: "Stop the server" },

  // Map / Game-Mode / Round Management
  { command: "map <mapname>", description: "Change to specified map" },
  { command: "map_rotate", description: "Load next map in rotation" },
  { command: "map_restart", description: "Restart current map/round" },
  {
    command: "fast_restart",
    description: "Fast map/round restart without reload",
  },
  {
    command: "g_gametype <type>",
    description: "Change game mode (dm, war, dom, sd, koth, sab)",
  },
  { command: "set <cvar> <value>", description: "Set server variable value" },

  // Player Management (Kick / Ban / Messaging)
  { command: "kick <name>", description: "Kick player by name" },
  { command: "clientkick <id>", description: "Kick player by ID/slot number" },
  { command: "kickall", description: "Kick all players (empty server)" },
  { command: "banUser <name>", description: "Permanently ban player by name" },
  { command: "banClient <id>", description: "Permanently ban player by ID" },
  {
    command: "tempBanUser <name>",
    description: "Temporarily ban player by name",
  },
  {
    command: "tempBanClient <id>",
    description: "Temporarily ban player by ID",
  },
  { command: "unbanUser <name>", description: "Remove ban by player name" },
  {
    command: "dumpuser <name>",
    description: "Show detailed player info (GUID, IP, etc)",
  },
  {
    command: "tell <id> <message>",
    description: "Send private message to player",
  },

  // Extended / CoD4X Admin Commands
  {
    command: "AdminAddAdmin <name>",
    description: "Add admin privileges to player",
  },
  {
    command: "AdminRemoveAdmin <name>",
    description: "Remove admin privileges",
  },
  {
    command: "AdminChangeCommandPower",
    description: "Change command permission level",
  },
  { command: "AdminListAdmins", description: "List all admins" },
  {
    command: "AdminListCommands",
    description: "List available admin commands",
  },

  // Ban List Management
  { command: "dumpbanlist", description: "Display all banned players" },
  { command: "unban <id>", description: "Remove ban by ID" },
  { command: "permban <id>", description: "Permanently ban player" },
  {
    command: "tempban <id> <duration>",
    description: "Temporarily ban with duration",
  },

  // Server Execution & Plugin Handling
  { command: "loadPlugin <name>", description: "Load server plugin" },
  { command: "unloadPlugin <name>", description: "Unload server plugin" },
  { command: "pluginInfo <name>", description: "Show plugin information" },
  { command: "plugins", description: "List all loaded plugins" },

  // Server Stats & Diagnostics
  { command: "meminfo", description: "Display memory information" },
  { command: "zonememinfo", description: "Display zone memory info" },
  { command: "net_restart", description: "Restart network connection" },
  { command: "info", description: "Display server info" },
  { command: "getmodules", description: "List loaded server modules" },

  // Messaging & Console Utilities
  { command: "screensay <message>", description: "Display message on screen" },
  {
    command: "screentell <id> <msg>",
    description: "Display message to specific player",
  },
  { command: "echo <text>", description: "Echo text to console" },
  { command: "which <command>", description: "Show command information" },
  { command: "vstr <variable>", description: "Execute variable string" },

  // Game Settings & DVARs
  { command: "sv_maxclients <number>", description: "Set max players" },
  { command: "g_speed <number>", description: "Set player movement speed" },
  { command: "g_gravity <number>", description: "Set gravity (default 800)" },
  { command: "g_knockback <number>", description: "Set knockback amount" },
  { command: "g_password <password>", description: "Set server password" },
  {
    command: "sv_privatePassword <pw>",
    description: "Set private slot password",
  },
  {
    command: "scr_game_matchstarttime <s>",
    description: "Match start countdown",
  },

  // Team Settings
  { command: "g_teamname_allies <name>", description: "Set allies team name" },
  { command: "g_teamname_axis <name>", description: "Set axis team name" },
  { command: "g_teamcolor_allies <color>", description: "Set allies color" },
  { command: "g_teamcolor_axis <color>", description: "Set axis color" },

  // Gameplay DVARs
  { command: "scr_dm_scorelimit <score>", description: "DM score limit" },
  { command: "scr_dm_timelimit <minutes>", description: "DM time limit" },
  { command: "scr_war_scorelimit <score>", description: "TDM score limit" },
  { command: "scr_war_timelimit <minutes>", description: "TDM time limit" },
  { command: "scr_sd_scorelimit <rounds>", description: "S&D round limit" },
  { command: "scr_sd_timelimit <minutes>", description: "S&D round time" },
  {
    command: "scr_dom_scorelimit <score>",
    description: "Domination score limit",
  },
  {
    command: "scr_koth_scorelimit <score>",
    description: "King of the Hill score limit",
  },
  {
    command: "scr_sab_scorelimit <score>",
    description: "Sabotage score limit",
  },

  // Weapon Settings
  { command: "scr_weapon_allowc4 <0|1>", description: "Allow C4" },
  { command: "scr_weapon_allowclaymore <0|1>", description: "Allow claymores" },
  { command: "scr_weapon_allowfrag <0|1>", description: "Allow frag grenades" },
  { command: "scr_weapon_allowflash <0|1>", description: "Allow flashbangs" },
  {
    command: "scr_weapon_allowsmoke <0|1>",
    description: "Allow smoke grenades",
  },
  { command: "scr_weapon_allowrpg <0|1>", description: "Allow RPGs" },

  // Killstreak Settings
  { command: "scr_team_killstreak <0|1>", description: "Enable killstreaks" },
  { command: "scr_hardpoint_allowuav <0|1>", description: "Allow UAV" },
  {
    command: "scr_hardpoint_allowartillery <0|1>",
    description: "Allow airstrike",
  },
  {
    command: "scr_hardpoint_allowhelicopter <0|1>",
    description: "Allow helicopter",
  },

  // HUD Settings
  {
    command: "scr_hud_showobjicons <0|1>",
    description: "Show objective icons",
  },
  { command: "scr_drawfriend <0|1>", description: "Show friendly indicators" },
  { command: "scr_game_allowkillcam <0|1>", description: "Allow killcam" },
  {
    command: "scr_game_spectatetype <0-2>",
    description: "Spectator mode (0=off, 1=team, 2=free)",
  },

  // Server DVARs
  { command: "sv_hostname <name>", description: "Set server name" },
  { command: "sv_maxRate <rate>", description: "Max data rate" },
  { command: "sv_fps <fps>", description: "Server frame rate" },
  { command: "sv_pure <0|1>", description: "Pure server mode" },
  { command: "sv_voice <0|1>", description: "Enable voice chat" },
  { command: "sv_allowdownload <0|1>", description: "Allow file downloads" },

  // Network Settings
  { command: "net_ip <ip>", description: "Set server IP" },
  { command: "net_port <port>", description: "Set server port" },

  // Logging
  { command: "g_log <filename>", description: "Set log file name" },
  { command: "g_logsync <0|1|2>", description: "Log sync mode" },
  { command: "logfile <0|1|2>", description: "Enable logging" },

  // Admin / Debug
  { command: "sv_cheats <0|1>", description: "Enable/disable cheats" },
  { command: "developer <0|1|2>", description: "Set developer mode" },
  { command: "devmap <mapname>", description: "Load map with cheats enabled" },

  // Punkbuster
  { command: "pb_sv_enable", description: "Enable Punkbuster" },
  { command: "pb_sv_kick <id> <reason>", description: "PB kick player" },
  { command: "pb_sv_ban <id> <reason>", description: "PB ban player" },

  // Console Utilities
  { command: "cmdlist", description: "List all commands" },
  { command: "cvarlist", description: "List all cvars" },
  { command: "dir <path>", description: "List directory contents" },
];

export const commandSuggestions = cod4Commands.map((cmd) => cmd.command);

export const getCommandDescription = (command: string): string | undefined => {
  const cmd = cod4Commands.find((c) =>
    c.command.toLowerCase().startsWith(command.toLowerCase())
  );
  return cmd?.description;
};

export const getMatchingCommands = (partial: string): typeof cod4Commands => {
  if (!partial) return cod4Commands;
  const lower = partial.toLowerCase();
  return cod4Commands.filter(
    (cmd) =>
      cmd.command.toLowerCase().includes(lower) ||
      cmd.description.toLowerCase().includes(lower)
  );
};
