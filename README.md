# Civo Kubernetes Cluster with Argo CD

This Pulumi project creates a Kubernetes cluster on Civo and installs Argo CD for GitOps-based continuous delivery.

## Prerequisites

- [Pulumi](https://www.pulumi.com/docs/get-started/install/) installed
- [Go](https://golang.org/doc/install) 1.23 or later
- [Civo CLI](https://github.com/civo/cli) installed and configured
- A Civo account with API access

## Project Structure

```
.
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── main.go          # Main Pulumi program
├── Pulumi.yaml      # Pulumi project configuration
└── Pulumi.dev.yaml  # Pulumi stack configuration
```

## Infrastructure Components

This project creates the following infrastructure:

1. **Civo Firewall**
   - Allows inbound traffic on ports 80, 443, and 6443
   - Configured in the FRA1 region

2. **Civo Kubernetes Cluster**
   - Kubernetes version: 1.30.5-k3s1
   - CNI: Cilium
   - Node size: g4s.kube.xsmall
   - Node count: 1
   - Region: FRA1
   - Pre-installed applications: metrics-server

3. **Argo CD**
   - Version: 7.8.24
   - Installed via Helm
   - Service type: NodePort
   - Namespace: argocd

## Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd pulumi-civo
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure Pulumi:
   ```bash
   pulumi stack init dev
   ```

4. Deploy the infrastructure:
   ```bash
   pulumi up
   ```

## Accessing Argo CD

After deployment, you can access the Argo CD UI through the NodePort service. To get the admin password:

```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

## Cleanup

To destroy all resources:

```bash
pulumi destroy
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 