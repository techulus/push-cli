# Push CLI

A command-line tool for sending push notifications via [Push by Techulus](https://push.techulus.com).

## Install

### Homebrew

```bash
brew install techulus/tap/push
```

### From source

```bash
go install github.com/techulus/push-cli@latest
```

### Binary releases

Download pre-built binaries from [GitHub Releases](https://github.com/techulus/push-cli/releases).

## Setup

Get your API key from [Push by Techulus](https://push.techulus.com) and configure the CLI:

```bash
push config set-key <your-api-key>
```

Verify your configuration:

```bash
push config show
```

## Usage

### Send a notification

```bash
push notify --title "Deploy Complete" --body "Production v2.1.0 is live"
```

### With optional flags

```bash
push notify \
  --title "Alert" \
  --body "CPU usage above 90%" \
  --sound default \
  --channel monitoring \
  --link "https://grafana.example.com/dashboard" \
  --image "https://example.com/chart.png" \
  --time-sensitive
```

### Pipe body from stdin

```bash
echo "Build failed on main" | push notify --title "CI Alert"
```

```bash
cat error.log | push notify --title "Error Log"
```

```bash
push notify --title "Disk Usage" --body - <<< "$(df -h /)"
```

### Send async

```bash
push notify-async --title "Queued" --body "This is processed asynchronously"
```

### Send to a group

```bash
push notify-group my-team --title "Standup" --body "Daily standup in 5 minutes"
```

### Available sounds

`default`, `arcade`, `correct`, `fail`, `harp`, `reveal`, `bubble`, `doorbell`, `flute`, `money`, `scifi`, `clear`, `elevator`, `guitar`, `pop`

## Configuration

The API key is stored in `~/.push/config.yaml`.

## License

MIT
