
# 389 Directory Server Prometheus Exporter

This exporter collects and exposes metrics from a [3️⃣8️⃣9️⃣ Directory Server](https://directory.fedoraproject.org/) instance, making it possible to monitor the performance and health of your directory server using Prometheus.

## Features

- Easy as possible to configure, deploy and modify
- Very quickly collects a variety of metrics from the 389 Directory Server Exposes metrics in a format compatible with Prometheus
- Groups metrics by label for easier use in processing and take care of cardinality

## Metrics

The exporter collects various metrics from the 389 Directory Server, including:

| Metric                                                                         | Type    | Description                                                                                                                |
|--------------------------------------------------------------------------------|---------|----------------------------------------------------------------------------------------------------------------------------|
| ds389_scape_duration_millis                                                    | GAUGE   | Duration of the request to obtain metrics in milliseconds                                                                  |
| ds389_version{version="..."}                                                   | GAUGE   | Version 389 Directory Server                                                                                               |
| ds389_uptime                                                                   | GAUGE   | Uptime in seconds                                                                                                          |
| ds389_startup_timestamp                                                        | GAUGE   | Server startup unix timestamp                                                                                              |
| ds389_current_connections                                                      | GAUGE   | Number of active connections to the server                                                                                 |
| ds389_total_connections                                                        | COUNTER | Total number of connections made to the server                                                                             |
| ds389_max_threads_per_conn_hits                                                | COUNTER | Total number of times the maximum threads per connection limit has been reached                                            |
| ds389_connections_in_max_threads                                               | COUNTER | Number of connections currently utilizing the maximum allowed threads                                                      |
| ds389_connections_max_threads                                                  | COUNTER | Total number of concurrent connections that the LDAP server can service simultaneously using the maximum number of threads |
| ds389_current_connections_at_max_threads                                       | GAUGE   | Number of connections currently utilizing the maximum allowed threads per connection                                       |
| ds389_read_waiters                                                             | GAUGE   | Number of threads waiting for read operations to complete                                                                  |
| ds389_threads                                                                  | GAUGE   | Total number of threads currently active in the server                                                                     |
| ds389_operations{type="ADD / SEARCH / MODIFY / MODIFY_RND / COMPARE / DELETE"} | COUNTER | Total number of entry operations by type processed by the server                                                           |
| ds389_bind_operations{type="ANONYMOUS / UNAUTHORIZED / SIMPLE / STRONG"}       | COUNTER | Total number of bind operations performed on the server                                                                    |
| ds389_bind_errors                                                              | COUNTER | Total number of bind operations rejected due to security-related errors                                                    |
| ds389_rx_bytes                                                                 | COUNTER | Total number of bytes received by the server                                                                               |
| ds389_tx_bytes                                                                 | COUNTER | Total number of bytes transferred by the server                                                                            |
| ds389_returned_entries                                                         | COUNTER | Total number of entries returned by search operations on the server                                                        |
| ds389_completed_operations                                                     | COUNTER | Total number of directory operations completed by the server                                                               |
| ds389_initiated_operations                                                     | COUNTER | Total number of directory operations initiated by clients and processed by the server                                      |
| ds389_cache_entries                                                            | GAUGE   | Number of entries currently cached in the server's cache                                                                   |
| ds389_cache_hits_count                                                         | COUNTER | Total number of times entries have been found in the cache and returned without accessing the backend storage              |
| ds389_dtable_size                                                              | GAUGE   | Size of the Directory Server descriptor table                                                                              |
| ds389_search_operations_level{type="ONE / SUBTREE"}                            | COUNTER | Total number of search operations by level executed by the server                                                          |
| ds389_copy_entries                                                             | COUNTER | Total number of operations to copy or move entries between different containers or subsections                             |
| ds389_errors                                                                   | COUNTER | Total number of errors that occurred on the server                                                                         |
| ds389_number_backends                                                          | GAUGE   | Number of alternative backends or data stores in use that are integrated and used by the LDAP server                       |
| ds389_returned_referrals                                                       | COUNTER | Total number of referrals returned to the client in response to its request to the LDAP server                             |
| ds389_security_errors                                                          | COUNTER | Total number of security errors that occurred on the server                                                                |
| ds389_supplier_entries                                                         | COUNTER | Total number of entries that were received from the data provider                                                          |
| ds389_consumer_hits                                                            | COUNTER | Total number of times entries from the supplier's database have been accessed by consumer servers                          |

## Easy start

1️⃣ Create a configuration file named `config.yml` in the root directory of the exporter. The file should contain the following settings:

```yaml
server: ipa0.home.local
port: 389
bind_dn: uid=monitor,cn=users,cn=compat,dc=home,dc=local
bind_password: MyStrongPassword
```

- `server`: The hostname or IP address of your 389 Directory Server.
- `port`: The port on which your 389 Directory Server is running.
- `bind_dn`: The Distinguished Name (DN) to bind to the directory server.
- `bind_password`: The password for the bind DN.

2️⃣ Run the exporter using the following command:

```sh
./ds389-exporter 
```

#### Available parameters

| Parameter            | Default value | Description                                                                                   |
|----------------------|---------------|-----------------------------------------------------------------------------------------------|
| --config.file        | ./config.yml  | Path to configuration file                                                                    |
| --log.level          | info          | Only log messages with the given severity or above. One of: [debug, info, warn, error, fatal] |
| --web.listen.address | 0.0.0.0:9389  | Address to listen on for get telemetry                                                        |
| --web.telemetry-path | /metrics      | Path under which to expose metrics                                                            |

3️⃣ If the exporter launched successfully, you will see a log like this

```shell
time="2024-07-16T21:47:15+03:00" level=info msg="Starting ds389-exporter"
time="2024-07-16T21:47:15+03:00" level=info msg="Using config file: ./config.yml"
time="2024-07-16T21:47:15+03:00" level=info msg="Metrics are available at http://0.0.0.0:9389/metrics"
```

4️⃣ Visit page http://ipa0.home.local:9389/metrics for get metrics. if you have configured everything correctly you will see the [metrics](https://github.com/rorudzhov/ds389-exporter/blob/dev/metrics.prom)

5️⃣ I recommend using a [systemd unit](ds389-exporter.service) for autonomous operation and ease of management
```unit file (systemd)
[Unit]
Description=389 Directory Server Exporter
After=network.target

[Service]
Type=simple
User=prometheus
ExecStart=/opt/ds389-exporter/ds389-exporter \
                --config.file /opt/ds389-exporter/config.yml \
                --log.level info \
                --web.listen.address 0.0.0.0:9389 \
                --web.telemetry-path /metrics
Restart=on-failure
RestartSec=15s

[Install]
WantedBy=multi-user.target
```

## Ansible
For convenient mass deployment, you can use a ready-made [ansible role](ansible%2Froles%2Fds389-exporter). Just run it
```shell
ansible-playbook -i ipa0.home.local, playbook.yml -e "role=ds389-exporter" -t install
```
To uninstall use the `-t uninstall` tag

## Build

If the release does not have a version for your bundle, you can build it yourself, use this script. The parameters GOOS and GOARCH can be found in the official documentation [Golang specification](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63)
```sh
https://github.com/rorudzhov/ds389-exporter.git
cd 389-ds-exporter

GOOS=linux # REQUIRED REPLACE
GOARCH=amd64 # REQUIRED REPLACE

export GOARCH=$GOARCH
export GOOS=$GOOS
file="ds389-exporter"
go build -ldflags "-s -w" -o $file main.go
```

## Prometheus Configuration

Add the following job to your Prometheus configuration file (`prometheus.yml`):

```yaml
scrape_configs:
  - job_name: 'ds389'
    metrics_path: "/metrics"
    static_configs:
      - targets: ['ipa0.home.local:9389']
```

Adjust the target to match the hostname and port where your exporter is running.

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes.
