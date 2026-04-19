# gekka-metrics

A full [gekka](https://github.com/sopranoworks/gekka) cluster node that joins
the ring with the `metrics-exporter` role and exports cluster metrics via the
OpenTelemetry Protocol (OTLP/HTTP).

Because `gekka-metrics` is a real cluster node it sees membership changes in
real time via gossip, giving sub-second lag between a node going unreachable
and the metric updating.

## Installation

```bash
go install github.com/sopranoworks/gekka-metrics@latest
```

## Running

```bash
gekka-metrics --config cluster.conf [--otlp http://otel-collector:4318]
```

| Flag | Default | Description |
|---|---|---|
| `--config FILE` | *(required)* | Path to a HOCON cluster config |
| `--otlp ENDPOINT` | *(empty -- local only)* | OTLP/HTTP endpoint to push metrics to |

**Minimal HOCON config:**

```hocon
pekko {
  remote.artery.canonical {
    hostname = "127.0.0.1"
    port     = 2560
  }
  cluster.seed-nodes = ["pekko://ClusterSystem@127.0.0.1:2552"]
}

gekka.telemetry.exporter.otlp {
  endpoint = "http://otel-collector:4318"
}
```

The `metrics-exporter` role is injected automatically so sharding allocators
and singleton managers exclude this node from hosting production workloads.

When no OTLP endpoint is configured the process still joins the cluster and
displays a live TUI view of membership state.

**Log Verbosity:**
Set `gekka.logging.level = "DEBUG"` in your HOCON configuration to enable
detailed protocol tracing.

## Integration with Prometheus / Grafana

Use the OpenTelemetry Collector as the bridge:

```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: "0.0.0.0:4318"

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
  debug:
    verbosity: basic

service:
  pipelines:
    metrics:
      receivers: [otlp]
      exporters: [prometheus, debug]
```

Then add the Prometheus scrape target:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: gekka
    static_configs:
      - targets: ['otel-collector:8889']
```

## Key Metrics

| Metric | Type | Unit | Attributes | Description |
|---|---|---|---|---|
| `gekka.cluster.members` | ObservableGauge | `{members}` | `status`, `dc` | Members grouped by status and data-center |

**Attribute values for `status`:** `up`, `joining`, `leaving`, `exiting`, `down`, `weakly-up`, `removed`

## License

MIT
