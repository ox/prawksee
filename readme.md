# Prawksee

Proxies requests to backend servers.

Configuration looks like:

```toml
[servers.hello]
  bind = 8001
  network = 'tcp'
  address = 'localhost:9001'

[servers.world]
  bind = 8002
  network = 'tcp'
  address = 'localhost:9002'
```

Run it:

```
% prawksee -config path-to-config.toml
```