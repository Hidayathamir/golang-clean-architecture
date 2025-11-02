mkdir -p signoz/deploy/common/clickhouse
mkdir -p signoz/deploy/common/signoz
mkdir -p signoz/deploy/docker

curl -fsSL https://raw.githubusercontent.com/SigNoz/signoz/refs/heads/main/deploy/common/clickhouse/config.xml -o signoz/deploy/common/clickhouse/config.xml
curl -fsSL https://raw.githubusercontent.com/SigNoz/signoz/refs/heads/main/deploy/common/clickhouse/users.xml -o signoz/deploy/common/clickhouse/users.xml
curl -fsSL https://raw.githubusercontent.com/SigNoz/signoz/refs/heads/main/deploy/common/clickhouse/custom-function.xml -o signoz/deploy/common/clickhouse/custom-function.xml
curl -fsSL https://raw.githubusercontent.com/SigNoz/signoz/refs/heads/main/deploy/common/clickhouse/cluster.xml -o signoz/deploy/common/clickhouse/cluster.xml
curl -fsSL https://raw.githubusercontent.com/SigNoz/signoz/refs/heads/main/deploy/common/signoz/prometheus.yml -o signoz/deploy/common/signoz/prometheus.yml
curl -fsSL https://raw.githubusercontent.com/SigNoz/signoz/refs/heads/main/deploy/docker/otel-collector-config.yaml -o signoz/deploy/docker/otel-collector-config.yaml
curl -fsSL https://raw.githubusercontent.com/SigNoz/signoz/refs/heads/main/deploy/common/signoz/otel-collector-opamp-config.yaml -o signoz/deploy/common/signoz/otel-collector-opamp-config.yaml
