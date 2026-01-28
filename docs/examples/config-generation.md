# Config Generation Examples

Generate configuration files from templates and data.

## Application Config

### Template

`config.yaml.tmpl`:
```yaml
app:
  name: {{ .app.name }}
  version: {{ .app.version }}
  debug: {{ .app.debug | default false }}

server:
  host: {{ .server.host | default "0.0.0.0" }}
  port: {{ .server.port | default 8080 }}
  timeout: {{ .server.timeout | default "30s" }}

database:
  driver: {{ .database.driver }}
  host: {{ .database.host }}
  port: {{ .database.port }}
  name: {{ .database.name }}
  {{- if .database.ssl }}
  ssl_mode: require
  {{- end }}

{{- if .features }}
features:
{{- range .features }}
  {{ . }}: true
{{- end }}
{{- end }}

logging:
  level: {{ .logging.level | default "info" }}
  format: {{ .logging.format | default "json" }}
```

### Data

`values.yaml`:
```yaml
app:
  name: myservice
  version: 1.2.0
  debug: false

server:
  port: 3000
  timeout: 60s

database:
  driver: postgres
  host: db.example.com
  port: 5432
  name: myservice_prod
  ssl: true

features:
  - auth
  - metrics
  - tracing

logging:
  level: warn
```

### Command

```bash
render config.yaml.tmpl values.yaml -o config.yaml
```

### Output

`config.yaml`:
```yaml
app:
  name: myservice
  version: 1.2.0
  debug: false

server:
  host: 0.0.0.0
  port: 3000
  timeout: 60s

database:
  driver: postgres
  host: db.example.com
  port: 5432
  name: myservice_prod
  ssl_mode: require

features:
  auth: true
  metrics: true
  tracing: true

logging:
  level: warn
  format: json
```

## Environment-Specific Configs

Generate configs for multiple environments.

### Template

`env-config.yaml.tmpl`:
```yaml
environment: {{ .environment }}
api_url: {{ .api_url }}
log_level: {{ .log_level }}
replicas: {{ .replicas }}

{{- if eq .environment "production" }}
monitoring:
  enabled: true
  endpoint: https://metrics.example.com
{{- end }}
```

### Data Files

`dev.yaml`:
```yaml
environment: development
api_url: http://localhost:3000
log_level: debug
replicas: 1
```

`prod.yaml`:
```yaml
environment: production
api_url: https://api.example.com
log_level: warn
replicas: 3
```

### Commands

```bash
render env-config.yaml.tmpl dev.yaml -o config/development.yaml
render env-config.yaml.tmpl prod.yaml -o config/production.yaml
```

## Docker Compose

### Template

`docker-compose.yaml.tmpl`:
```yaml
version: "3.8"

services:
{{- range .services }}
  {{ .name }}:
    image: {{ .image }}:{{ .tag | default "latest" }}
    {{- if .ports }}
    ports:
      {{- range .ports }}
      - "{{ . }}"
      {{- end }}
    {{- end }}
    {{- if .environment }}
    environment:
      {{- range $key, $value := .environment }}
      {{ $key }}: {{ $value | quote }}
      {{- end }}
    {{- end }}
    {{- if .volumes }}
    volumes:
      {{- range .volumes }}
      - {{ . }}
      {{- end }}
    {{- end }}
    {{- if .depends_on }}
    depends_on:
      {{- range .depends_on }}
      - {{ . }}
      {{- end }}
    {{- end }}
{{ end }}
{{- if .networks }}
networks:
{{- range $name, $config := .networks }}
  {{ $name }}:
    driver: {{ $config.driver | default "bridge" }}
{{- end }}
{{- end }}
```

### Data

`services.yaml`:
```yaml
services:
  - name: web
    image: nginx
    tag: alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - api

  - name: api
    image: myapp/api
    tag: v1.2.0
    ports:
      - "3000:3000"
    environment:
      DATABASE_URL: postgres://db:5432/myapp
      REDIS_URL: redis://cache:6379
    depends_on:
      - db
      - cache

  - name: db
    image: postgres
    tag: "15"
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: app
      POSTGRES_PASSWORD: secret
    volumes:
      - pgdata:/var/lib/postgresql/data

  - name: cache
    image: redis
    tag: "7-alpine"

networks:
  default:
    driver: bridge
```

### Command

```bash
render docker-compose.yaml.tmpl services.yaml -o docker-compose.yaml
```

## Nginx Config

### Template

`nginx.conf.tmpl`:
```nginx
worker_processes auto;

events {
    worker_connections {{ .worker_connections | default 1024 }};
}

http {
    include mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent"';

    {{- range .upstreams }}
    upstream {{ .name }} {
        {{- range .servers }}
        server {{ .host }}:{{ .port }} weight={{ .weight | default 1 }};
        {{- end }}
    }
    {{- end }}

    {{- range .servers }}
    server {
        listen {{ .listen | default 80 }};
        server_name {{ .server_name }};

        {{- if .ssl }}
        listen 443 ssl;
        ssl_certificate {{ .ssl.cert }};
        ssl_certificate_key {{ .ssl.key }};
        {{- end }}

        {{- range .locations }}
        location {{ .path }} {
            {{- if .proxy_pass }}
            proxy_pass {{ .proxy_pass }};
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            {{- else if .root }}
            root {{ .root }};
            index index.html;
            {{- end }}
        }
        {{- end }}
    }
    {{- end }}
}
```

### Data

`nginx.yaml`:
```yaml
worker_connections: 2048

upstreams:
  - name: api_backend
    servers:
      - host: 10.0.0.1
        port: 3000
        weight: 2
      - host: 10.0.0.2
        port: 3000
        weight: 1

servers:
  - server_name: example.com
    ssl:
      cert: /etc/ssl/certs/example.com.crt
      key: /etc/ssl/private/example.com.key
    locations:
      - path: /
        root: /var/www/html
      - path: /api
        proxy_pass: http://api_backend
```

### Command

```bash
render nginx.conf.tmpl nginx.yaml -o nginx.conf
```

## Kubernetes ConfigMap

### Template

`configmap.yaml.tmpl`:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .name }}
  namespace: {{ .namespace | default "default" }}
  labels:
    app: {{ .app }}
    {{- range $key, $value := .labels }}
    {{ $key }}: {{ $value }}
    {{- end }}
data:
  {{- range $key, $value := .data }}
  {{ $key }}: |
{{ $value | indent 4 }}
  {{- end }}
```

### Data

`cm-data.yaml`:
```yaml
name: app-config
namespace: production
app: myapp
labels:
  version: v1
  team: backend
data:
  config.json: |
    {
      "debug": false,
      "port": 8080
    }
  settings.ini: |
    [database]
    host = db.example.com
    port = 5432
```

### Command

```bash
render configmap.yaml.tmpl cm-data.yaml -o k8s/configmap.yaml
```
