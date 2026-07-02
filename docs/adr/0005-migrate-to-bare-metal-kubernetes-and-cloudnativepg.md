# 5. Migrate to Bare Metal Kubernetes and CloudNativePG

Date: 2026-07-02

## Status

Accepted

Supercedes [3. Selection of AWS as Primary Cloud Provider](0003-selection-of-aws-as-primary-cloud-provider.md)

Supercedes [4. Use of AWS Elastic Container Service](0004-use-of-aws-elastic-container-service.md)

## Context

[ADR 0003](0003-selection-of-aws-as-primary-cloud-provider.md) chose AWS on the explicit basis that self-hosted infrastructure was "not currently feasible" due to lack of physical space, staffing capacity, and up-front capital — with a stated intent to migrate to open, self-operated infrastructure once the project matured. [ADR 0004](0004-use-of-aws-elastic-container-service.md) reaffirmed ECS as "a stable interim solution until we reach the level of scale that justifies the move to Kubernetes."

We now operate a bare metal Kubernetes cluster with the operational building blocks that migration depended on: GitOps deployment via ArgoCD, Rook-Ceph for persistent storage, the CloudNativePG operator for managing Postgres, and Gateway API-based ingress. This removes the original blockers to self-hosting and lets us retire the AWS dependency for this application.

Planzoco's data model (events, questions, options with vote counts) does not require DynamoDB's scale characteristics — it is a small relational dataset with clear foreign-key relationships between entities. Postgres, managed in-cluster by CloudNativePG, is a better fit and removes the single-table DynamoDB design (manual PK/SK composition, GSIs standing in for what are really joins) that added complexity without a corresponding benefit at this scale.

## Decision

We are migrating planzoco off AWS App Runner/ECS and DynamoDB onto our own bare metal Kubernetes cluster:

- The application is redeployed as a standard Kubernetes `Deployment`/`Service`, packaged as a Helm chart in this repository's `chart/` directory and registered with ArgoCD, following the same pattern already established for other in-house services.
- DynamoDB is replaced with a CloudNativePG-managed Postgres `Cluster`, provisioned outside this repository (in `kubernetes-lab-services`, alongside this app's other cluster-side resources) rather than an external managed database service.
- The Go application connects via a standard `DATABASE_URL`, sourced from the Secret CloudNativePG generates automatically, instead of the AWS SDK.
- Container images continue to be published to Quay via the existing CI pipeline; only the deployment target changes.

AWS infrastructure (`tofu/production`) is left in place for now rather than torn down in the same change, so the two deployments can run side by side until the Kubernetes deployment is verified in production.

## Consequences

- We give up AWS's managed elasticity (DynamoDB's on-demand scaling, App Runner/ECS's managed compute) in exchange for operating our own compute, storage, and database — a trade-off we're now equipped to make per the building blocks described above. In return, we remove the risk of a surprise AWS bill from unexpected traffic, a runaway process, or a misconfiguration — bare metal capacity is fixed and already paid for, so there's no equivalent failure mode.
- At planzoco's current scale, the loss of elasticity isn't a real constraint: the chart's `HorizontalPodAutoscaler` can still scale the `Deployment` up to 10 replicas across the existing bare metal nodes, which is comfortably more headroom than this app needs for the foreseeable future.
- Backups, failover, and scaling for Postgres become our own responsibility via CloudNativePG, rather than AWS's.
- The AWS tofu module, DynamoDB table, and associated Cloudflare DNS records remain live until we explicitly decommission them in a follow-up change; this ADR does not itself constitute that cutover.
- This establishes the pattern other AWS-hosted services in the estate are expected to follow as they're migrated.
