# Birthdays API

## Key Points

- Built using GoLang
- Deployed using Docker + Kubernetes
- Throughput throttled and protected from DoS/DDoS to some extent through the use of HAProxy
- Autoscaling setup using Kubernetes HPA's
- All data stored in Redis
- Full test suite tetsing both the individual handlers as well as the HTTP server
- OpenAPI Specification
- Kubernetes overlay for AWS and ALB
  - Various configuration changes made to work better in an AWS environment
  - Zero downtime deployments and scaling
  - Will instruct Kubernretes to make a best attempt at providing a highly available setup
    - Pod affinity is configured to attempt to put pods for the application on separate nodes when possible
    - TopologySpreadConstraints are setup to make a best attempt at spreading the application across zones where possible
- Minimal Docker image built on Scratch to ensure security and footprint
- API response time < 3ms on average
- In-Memory based caching to ease load on backend storage service

## Deployment

A docker image is already build and hosted using GitHub packages (ghcr.io).

A setup has been provided for a Kubernetes cluster without needing any modification. This setup can be used by deploying the `aws` overlay.

Be sure to edit `./k8s/overlays/aws/ingress.yaml` to use a domain you have access to.

```shell
kubectl -n <my_namespace> -k ./k8s/overlays/aws
```

This makes the following assumptions:
- You have a Kubernetes Cluster in AWS (EKS or other) running Kubernetes 1.24
- You have ExternalDNS running and configured to use Route53
- You have AWS Load Balancer Controller running and configured
- You have at least 1 node group (or alternative if using Karpenter or Fargate)