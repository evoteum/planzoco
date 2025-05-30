# 4. Use of AWS Elastic Container Service

Date: 2025-04-21

## Status

Accepted

Supercedes [2. Use of AWS App Runner](0002-use-of-aws-app-runner.md)

## Context

AWS App Runner was initially selected for its promise of simplicity and low operational overhead. The goal was to deploy a containerised Go application with minimal configuration, ideally using the source deployment option. However, despite following AWS documentation and best practices, we encountered persistent issues that consumed significant development time.

The service repeatedly failed at either the build or start stages, with inconsistent and often misleading error messages. We explored various debugging strategies, including restructuring the repository to match expected patterns, using container builds instead of source code builds, adjusting IAM policies, and simplifying the codebase layout. Even with a working container locally and in other environments, App Runner would either fail silently or exit with a zero exit code, offering no actionable insights.

In parallel, we observed that App Runner has limited adoption in enterprise environments and lacks the transparency and control necessary for production-grade deployments. Logs are often unhelpful, and its internal behaviour during builds and deployment is opaque.

As noted in [ADR 0002](0002-use-of-aws-app-runner.md), our long-term plan is still to migrate to Kubernetes once the scale and operational complexity of the system warrant it. App Runner was intended as a stepping stone, but has proven too unreliable for even early-stage deployments.

## Decision

We have decided to switch to AWS Elastic Container Service (ECS), using Fargate for simplified container orchestration without managing EC2 instances. ECS provides full control over container runtime, IAM permissions, networking, and logging. The service is more commonly used in enterprise contexts, with broader community support, better documentation, and predictable behaviour when pulling and running containers.

We will publish and manage our containers through `quay.io`, and ECS will pull from this registry during deployments.

## Consequences

This decision increases the initial setup complexity slightly, as ECS requires additional resources like task definitions, services, and roles. However, this is mitigated by the use of infrastructure as code and the ability to fully define and version control all components. In return, we gain reliability, observability, and compatibility with industry standards.

Abandoning App Runner avoids further unproductive debugging time and aligns our infrastructure with commonly used patterns in production and enterprise environments. ECS serves as a stable interim solution until we reach the level of scale that justifies the move to Kubernetes.
