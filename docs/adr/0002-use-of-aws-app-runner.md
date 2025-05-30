# 2. Use of AWS App Runner

Date: 2025-04-04

## Status

Superceded by [4. Use of AWS Elastic Container Service](0004-use-of-aws-elastic-container-service.md)

## Context

We need a straightforward and fully managed way to deploy containerised web applications without requiring extensive infrastructure management.

As a fully FLOSS organisation, we recognise the long-term benefits of adopting a portable, open, and community-supported deployment platform. In the fullness of time, we intend to move towards a golden path built on Kubernetes. However, adopting Kubernetes today would be premature and overly complex for our current scale and needs. We need something simple, reliable, and cost-effective that gets out of the way and allows us to focus on delivering features.

## Decision

We will use AWS App Runner to deploy and manage our containerised web applications. This service simplifies deployment by automatically building and running applications from source code or a container registry, and scales based on demand with minimal configuration. It requires no infrastructure management and integrates easily into our CI/CD processes.

## Consequences

Deployment complexity will be significantly reduced, especially for small services and proof-of-concept applications. Operational burden around scaling, health checks, and networking will also decrease, allowing the team to focus more on development. However, we will need to monitor the limitations of App Runner in terms of customisability, such as finer-grained scaling controls and access to underlying infrastructure. There may also be a cost trade-off for certain workloads depending on traffic patterns or specific performance needs.

We will continue to evaluate the maturity of our internal platform needs, and revisit a transition to Kubernetes once the complexity it introduces is justified by operational scale and developer experience gains.
