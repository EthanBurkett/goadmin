# Example Plugin

This is a simple example plugin demonstrating the GoAdmin plugin system.

## Building

To build this plugin:

```bash
go build -buildmode=plugin -o example.so example.go
```

## Installing

1. Build the plugin using the command above
2. Copy `example.so` to the `plugins/` directory in your GoAdmin installation
3. Restart GoAdmin or reload plugins through the UI

## Features

This example plugin demonstrates:

- Event subscriptions (player connect/disconnect)
- Custom command registration (!hello)
- Configuration storage
- RCON command execution
- Webhook dispatching

## Configuration

The plugin stores configuration using the ConfigAPI. Example:

- `greeting_message`: The message shown when players use !hello
