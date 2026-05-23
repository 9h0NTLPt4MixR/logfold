# logfold

A structured log aggregator that collapses repetitive log bursts and surfaces anomaly patterns via a terminal dashboard.

---

## Installation

```bash
go install github.com/yourusername/logfold@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logfold.git
cd logfold
go build -o logfold ./cmd/logfold
```

---

## Usage

Pipe any structured (JSON) log stream into `logfold` to start the dashboard:

```bash
tail -f /var/log/app.log | logfold
```

Or point it at a file directly:

```bash
logfold --file /var/log/app.log
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--file` | stdin | Log file to read from |
| `--threshold` | `5` | Burst collapse threshold (repeat count) |
| `--interval` | `5s` | Aggregation window duration |
| `--format` | `json` | Input log format (`json`, `logfmt`) |

Once running, the terminal dashboard displays collapsed log groups, anomaly highlights, and live event rates. Press `q` to quit.

---

## Requirements

- Go 1.21+
- A terminal with 256-color support

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

---

## License

[MIT](LICENSE)