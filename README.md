# gekka-dashboard

Web-based operational console for [gekka](https://github.com/sopranoworks/gekka) clusters.
Provides real-time monitoring, management actions, and configurable notifications.

Derived from [gekka-metrics](https://github.com/sopranoworks/gekka-metrics) — includes
all metrics export and notification capabilities.

## Installation

```bash
go install github.com/sopranoworks/gekka-dashboard@latest
```

Or build from source:

```bash
make build
```

## Running

```bash
gekka-dashboard --config cluster.conf [--listen :9000] [--otlp http://otel-collector:4318]
```

| Flag | Default | Description |
|---|---|---|
| `--config FILE` | *(required)* | Path to a HOCON cluster config |
| `--listen ADDR` | `:9000` | Dashboard HTTP listen address |
| `--otlp ENDPOINT` | *(empty)* | OTLP/HTTP endpoint for metrics export |
| `--headless` | false | Disable UI, run as metrics-only exporter with notifications |

**Minimal HOCON config:**

```hocon
pekko {
  remote.artery.canonical {
    hostname = "127.0.0.1"
    port     = 2560
  }
  cluster.seed-nodes = ["pekko://ClusterSystem@127.0.0.1:2552"]
}

gekka.dashboard {
  listen = ":9000"   # HTTP listen address (default ":9000")
}

gekka.telemetry.exporter.otlp {
  endpoint = "http://otel-collector:4318"
}
```

The `dashboard` role is injected automatically so sharding allocators
and singleton managers exclude this node from hosting production workloads.

Open `http://localhost:9000` in your browser after starting.

## Headless Mode

With `--headless`, the dashboard disables the web UI and runs as a metrics exporter
with notification support — effectively replacing gekka-metrics.

## Notifications

Configure notification rules in your HOCON config:

```hocon
gekka.notifications {
  rules {
    critical-down {
      events = ["node.unreachable", "node.downed"]
      roles = ["cart", "payment"]
      channels = ["slack", "email"]
      throttle = 5m
    }
  }
  channels {
    slack { webhook-url = "https://hooks.slack.com/services/..." }
    email {
      smtp-host = "smtp.example.com"
      smtp-port = 587
      from = "alerts@example.com"
      to = ["ops@example.com"]
    }
  }
}
```

## Development

```bash
# Start frontend dev server (hot reload)
make dev

# In another terminal, run the Go backend
go run . --config your-cluster.conf --listen :9000
```

## License

This project is licensed under the MIT License — see [LICENSE](LICENSE) for details.

Derived from [gekka](https://github.com/sopranoworks/gekka) by Sopranoworks, Osamu Takahashi.
