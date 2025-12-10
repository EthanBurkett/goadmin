# Advanced Example Plugin

This plugin demonstrates advanced GoAdmin plugin features including:

- **Semantic Versioning**: Uses API version constraints (1.0.0 - 2.0.0)
- **Resource Limits**: Configured limits for memory, CPU, and goroutines
- **Dependency Management**: Can specify required plugins
- **Hot-Reload Support**: Properly implements Reload() for configuration changes
- **Event Handling**: Subscribes to player connect/disconnect events
- **Custom Commands**: Registers !status and !info commands
- **Webhook Integration**: Dispatches plugin lifecycle events

## Features

### Resource Limits

The plugin is configured with the following limits:

- **Max Memory**: 100MB
- **Max CPU**: 50%
- **Max Goroutines**: 50
- **Operation Timeout**: 30 seconds

These limits are monitored by the plugin manager and violations are logged.

### Commands

- `!status` - Shows plugin status and monitoring info
- `!info` - Shows detailed server information

### Events

Listens to:

- `player.connect` - Sends personalized welcome message
- `player.disconnect` - Logs player departure

### Webhooks

Dispatches:

- `plugin.started` - When plugin starts
- `plugin.stopped` - When plugin stops

## Building

### Linux/Mac

```bash
cd plugins/examples/advanced-example
go build -buildmode=plugin -o advanced-example.so advanced.go
```

### Windows

**Note**: Go plugins are not supported on Windows. Use Linux or WSL.

## Installation

1. Build the plugin as shown above
2. Copy `advanced-example.so` to the `plugins/` directory
3. Restart GoAdmin or use the hot-reload endpoint

## API Version Compatibility

This plugin requires:

- **Minimum API Version**: 1.0.0
- **Maximum API Version**: 2.0.0

If the GoAdmin API version is outside this range, the plugin will fail to load.

## Resource Monitoring

View plugin metrics via the API:

```bash
curl http://localhost:8080/plugins/advanced-example/metrics
```

Response:

```json
{
  "pluginId": "advanced-example",
  "memoryUsageMB": 45.2,
  "goroutineCount": 12,
  "lastChecked": "2025-12-09T10:30:00Z",
  "violationCount": 0,
  "throttled": false
}
```

## Hot Reload

To reload the plugin without stopping:

```bash
curl -X POST http://localhost:8080/plugins/advanced-example/hot-reload
```

This will:

1. Stop the plugin
2. Call Reload() to refresh configuration
3. Restart the plugin

## Dependency Example

To require the base example plugin, update the Dependencies field:

```go
Dependencies: []string{"example-plugin"},
```

The plugin manager will ensure dependencies are loaded first.

## Configuration

Plugin configuration can be stored/retrieved via the ConfigAPI:

```go
p.ctx.ConfigAPI.Set("welcome_enabled", true)
enabled := p.ctx.ConfigAPI.GetBool("welcome_enabled", true)
```

## Troubleshooting

### Plugin Won't Load

Check API version compatibility:

```bash
curl http://localhost:8080/plugins/advanced-example
```

### Resource Limit Violations

View metrics to see current usage:

```bash
curl http://localhost:8080/plugins/advanced-example/metrics
```

If violations occur, adjust limits in the Metadata() function or optimize your plugin.

### Dependency Errors

View dependency tree:

```bash
curl http://localhost:8080/plugins/advanced-example/dependencies
```

Ensure all required plugins are installed and started.
