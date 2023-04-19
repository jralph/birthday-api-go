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
- Minimal Docker image built on Scratch to ensure security and footprint
- API response time < 3ms on average