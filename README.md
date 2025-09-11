# s3-oauth2-proxy

[![CI](https://github.com/walnuts1018/s3-oauth2-proxy/actions/workflows/ci.yaml/badge.svg)](https://github.com/walnuts1018/s3-oauth2-proxy/actions/workflows/ci.yaml)
[![Docker](https://github.com/walnuts1018/s3-oauth2-proxy/actions/workflows/docker.yaml/badge.svg)](https://github.com/walnuts1018/s3-oauth2-proxy/actions/workflows/docker.yaml)

s3-oauth2-proxy is a reverse proxy that provides authentication and authorization for S3 buckets.

## Features

- Authenticate with OpenID Connect
- Authorize based on group claims
- Supports AWS S3 and other S3 compatible storages (e.g. MinIO)
- Assumes IAM Role for S3 access (e.g. EKS IAM Roles for Service Accounts, MinIO STS API)

## Configuration

The following environment variables are available for configuration:

| Name | Description | Default |
| --- | --- | --- |
| `OIDC_ISSUER_URL` | OIDC issuer URL | |
| `OIDC_CLIENT_ID` | OIDC client ID | |
| `OIDC_CLIENT_SECRET` | OIDC client secret | |
| `OIDC_REDIRECT_URL` | OIDC redirect URL. Use `/auth/callback` as the path. ||
| `OIDC_GROUP_CLAIM` | Group claim name | `groups` |
| `OIDC_ALLOWED_GROUPS` | Comma separated list of allowed groups/role. | |
| `SESSION_SECRET` | Secret for session | |
| `S3_BUCKET` | S3 bucket name | |
| `S3_USE_PATH_STYLE` | Use path style for S3 access | `false` |
| `LOG_LEVEL` | Log level | `info` |
| `LOG_TYPE` | Log type (json or text) | `json` |

Additionally, s3-oauth2-proxy supports AWS SDK environment variables (<https://docs.aws.amazon.com/sdkref/latest/guide/settings-reference.html#EVarSettings>).

## Example

Here is an example of running s3-oauth2-proxy on Kubernetes.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: s3-oauth2-proxy
  name: s3-oauth2-proxy
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: s3-oauth2-proxy
  template:
    metadata:
      labels:
        app: s3-oauth2-proxy
    spec:
      serviceAccountName: <your-service-account-name>
      containers:
        - name: proxy
          image: ghcr.io/walnuts1018/s3-oauth2-proxy:latest
          env:
            - name: OIDC_ISSUER_URL
              value: <your-oidc-issuer-url>
            - name: OIDC_CLIENT_ID
              valueFrom:
                secretKeyRef:
                  key: client-id
                  name: <your-secret-name>
            - name: OIDC_CLIENT_SECRET
              valueFrom:
                secretKeyRef:
                  key: client-secret
                  name: <your-secret-name>
            - name: OIDC_REDIRECT_URL
              value: <your-redirect-url>
            - name: OIDC_ALLOWED_GROUPS
              value: <your-allowed-groups>
            - name: OIDC_GROUP_CLAIM
              value: <your-group-claim>
            - name: SESSION_SECRET
              valueFrom:
                secretKeyRef:
                  key: session-secret
                  name: <your-secret-name>
            - name: S3_BUCKET
              value: <your-s3-bucket>
            - name: AWS_REGION
              value: <your-aws-region>
            - name: AWS_ROLE_ARN
              value: <your-aws-role-arn>
          livenessProbe:
            httpGet:
              path: /livez
              port: http
          readinessProbe:
            httpGet:
              path: /readyz
              port: http
          ports:
            - containerPort: 8080
              name: http
              protocol: TCP
          resources:
            limits:
              memory: 300Mi
            requests:
              cpu: 10m
              memory: 10Mi
```

Additionally, an example of using the MinIO Operator in an on-premises Kubernetes cluster can be found [here]([examples/minio-operator.yaml](https://github.com/walnuts1018/infra/blob/7642120ecb6f4b5dd415d85ea7bb5099fdcf4725/k8s/apps/ipu/deployment.yaml)).

## Development

### Certificate Issuance

```bash
make cert
```

### Startup

```bash
docker compose watch
```

### Terraform

```bash
terraform -chdir=terraform/local/ init
terraform -chdir=terraform/local/ plan
terraform -chdir=terraform/local/ apply
```

### Editing hosts

Add the following to `/etc/hosts`:

```bash
127.0.0.1 s3-oauth2-proxy.local.walnuts.dev
127.0.0.1 authelia.local.walnuts.dev
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
