# Multi-File Output Examples

Generate multiple output files from a single data source using each mode.

## User Profile Pages

Generate one HTML page per user.

### Template

`profile.html.tmpl`:
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .name }} - Profile</title>
</head>
<body>
    <h1>{{ .name }}</h1>
    <dl>
        <dt>Email</dt>
        <dd>{{ .email }}</dd>
        <dt>Role</dt>
        <dd>{{ .role | title }}</dd>
        <dt>Joined</dt>
        <dd>{{ .joinedAt }}</dd>
    </dl>
    {{- if .bio }}
    <section>
        <h2>About</h2>
        <p>{{ .bio }}</p>
    </section>
    {{- end }}
</body>
</html>
```

### Data

`users.json`:
```json
{
  "users": [
    {
      "username": "alice",
      "name": "Alice Smith",
      "email": "alice@example.com",
      "role": "admin",
      "joinedAt": "2023-01-15",
      "bio": "Software engineer with 10 years of experience."
    },
    {
      "username": "bob",
      "name": "Bob Jones",
      "email": "bob@example.com",
      "role": "developer",
      "joinedAt": "2023-03-22"
    },
    {
      "username": "charlie",
      "name": "Charlie Brown",
      "email": "charlie@example.com",
      "role": "designer",
      "joinedAt": "2023-06-01",
      "bio": "UX designer passionate about accessibility."
    }
  ]
}
```

### Command

```bash
render profile.html.tmpl users.json \
  --item-query '.users[]' \
  -o 'profiles/{{ .username }}.html'
```

### Output

```
profiles/
  alice.html
  bob.html
  charlie.html
```

## API Documentation

Generate documentation for each API endpoint.

### Template

`endpoint.md.tmpl`:
```markdown
# {{ .method }} {{ .path }}

{{ .description }}

## Request

**Method:** `{{ .method }}`
**Path:** `{{ .path }}`

{{- if .parameters }}

### Parameters

| Name | In | Type | Required | Description |
|------|-----|------|----------|-------------|
{{- range .parameters }}
| {{ .name }} | {{ .in }} | {{ .type }} | {{ if .required }}Yes{{ else }}No{{ end }} | {{ .description }} |
{{- end }}
{{- end }}

{{- if .requestBody }}

### Request Body

```json
{{ .requestBody | toPrettyJson }}
```
{{- end }}

## Response

**Status:** {{ .response.status }}

{{- if .response.body }}

```json
{{ .response.body | toPrettyJson }}
```
{{- end }}

{{- if .errors }}

## Errors

| Status | Description |
|--------|-------------|
{{- range .errors }}
| {{ .status }} | {{ .description }} |
{{- end }}
{{- end }}
```

### Data

`api.yaml`:
```yaml
endpoints:
  - method: GET
    path: /users
    description: List all users with optional filtering.
    parameters:
      - name: limit
        in: query
        type: integer
        required: false
        description: Maximum number of users to return
      - name: offset
        in: query
        type: integer
        required: false
        description: Number of users to skip
    response:
      status: 200
      body:
        users:
          - id: "123"
            name: "Alice"
        total: 1
    errors:
      - status: 401
        description: Unauthorized

  - method: POST
    path: /users
    description: Create a new user.
    requestBody:
      name: "string"
      email: "string"
    response:
      status: 201
      body:
        id: "456"
        name: "New User"
    errors:
      - status: 400
        description: Invalid request body
      - status: 409
        description: Email already exists

  - method: GET
    path: /users/{id}
    description: Get a user by ID.
    parameters:
      - name: id
        in: path
        type: string
        required: true
        description: User ID
    response:
      status: 200
      body:
        id: "123"
        name: "Alice"
        email: "alice@example.com"
    errors:
      - status: 404
        description: User not found
```

### Command

```bash
render endpoint.md.tmpl api.yaml \
  --item-query '.endpoints[]' \
  -o 'docs/api/{{ .method | lower }}-{{ .path | replace "/" "-" | trim "-" }}.md'
```

### Output

```
docs/api/
  get-users.md
  post-users.md
  get-users-id.md
```

## Environment Config Files

Generate config files for each environment.

### Template

`config.yaml.tmpl`:
```yaml
# Configuration for {{ .name }} environment
environment: {{ .name }}

server:
  host: {{ .server.host }}
  port: {{ .server.port }}

database:
  host: {{ .database.host }}
  port: {{ .database.port }}
  name: {{ .database.name }}
  pool_size: {{ .database.poolSize }}

logging:
  level: {{ .logging.level }}
  format: {{ .logging.format }}

features:
{{- range $feature, $enabled := .features }}
  {{ $feature }}: {{ $enabled }}
{{- end }}
```

### Data

`environments.yaml`:
```yaml
environments:
  - name: development
    server:
      host: localhost
      port: 3000
    database:
      host: localhost
      port: 5432
      name: app_dev
      poolSize: 5
    logging:
      level: debug
      format: text
    features:
      debug_mode: true
      hot_reload: true
      mock_services: true

  - name: staging
    server:
      host: 0.0.0.0
      port: 8080
    database:
      host: staging-db.example.com
      port: 5432
      name: app_staging
      poolSize: 10
    logging:
      level: info
      format: json
    features:
      debug_mode: false
      hot_reload: false
      mock_services: false

  - name: production
    server:
      host: 0.0.0.0
      port: 8080
    database:
      host: prod-db.example.com
      port: 5432
      name: app_prod
      poolSize: 50
    logging:
      level: warn
      format: json
    features:
      debug_mode: false
      hot_reload: false
      mock_services: false
```

### Command

```bash
render config.yaml.tmpl environments.yaml \
  --item-query '.environments[]' \
  -o 'config/{{ .name }}.yaml'
```

### Output

```
config/
  development.yaml
  staging.yaml
  production.yaml
```

## Kubernetes Manifests

Generate K8s resources for multiple services.

### Template

`deployment.yaml.tmpl`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .name }}
  namespace: {{ .namespace | default "default" }}
spec:
  replicas: {{ .replicas | default 1 }}
  selector:
    matchLabels:
      app: {{ .name }}
  template:
    metadata:
      labels:
        app: {{ .name }}
    spec:
      containers:
        - name: {{ .name }}
          image: {{ .image }}:{{ .tag | default "latest" }}
          ports:
            - containerPort: {{ .port }}
          {{- if .env }}
          env:
            {{- range $key, $value := .env }}
            - name: {{ $key }}
              value: {{ $value | quote }}
            {{- end }}
          {{- end }}
          resources:
            requests:
              cpu: {{ .resources.cpu | default "100m" }}
              memory: {{ .resources.memory | default "128Mi" }}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .name }}
  namespace: {{ .namespace | default "default" }}
spec:
  selector:
    app: {{ .name }}
  ports:
    - port: {{ .servicePort | default .port }}
      targetPort: {{ .port }}
```

### Data

`services.yaml`:
```yaml
namespace: production
services:
  - name: api-gateway
    image: myorg/api-gateway
    tag: v1.2.0
    port: 8080
    replicas: 3
    env:
      LOG_LEVEL: info
      RATE_LIMIT: "1000"
    resources:
      cpu: 500m
      memory: 512Mi

  - name: user-service
    image: myorg/user-service
    tag: v2.0.1
    port: 3000
    replicas: 2
    env:
      DATABASE_URL: postgres://db:5432/users
    resources:
      cpu: 200m
      memory: 256Mi

  - name: notification-service
    image: myorg/notification-service
    tag: v1.0.0
    port: 3001
    replicas: 1
    env:
      SMTP_HOST: smtp.example.com
```

### Command

```bash
render deployment.yaml.tmpl services.yaml \
  --query '. as $root | .services[] | . + {namespace: $root.namespace}' \
  --item-query '.' \
  -o 'k8s/{{ .name }}.yaml'
```

### Output

```
k8s/
  api-gateway.yaml
  user-service.yaml
  notification-service.yaml
```

## Filtering Items

Generate files only for items matching criteria.

### Only Active Users

```bash
render user.tmpl data.json \
  --item-query '.users[] | select(.active)' \
  -o '{{ .username }}.txt'
```

### Only Production Services

```bash
render service.tmpl data.json \
  --item-query '.services[] | select(.env == "production")' \
  -o '{{ .name }}.yaml'
```

### Sorted Output

```bash
render item.tmpl data.json \
  --item-query '[.items[] | select(.priority > 5)] | sort_by(.name) | .[]' \
  -o '{{ .name }}.txt'
```
