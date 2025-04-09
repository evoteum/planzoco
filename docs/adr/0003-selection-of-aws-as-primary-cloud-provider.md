# 3. Selection of AWS as Primary Cloud Provider

Date: 2025-04-08

## Status

Accepted

## Context

The project requires a reliable, scalable, and widely supported infrastructure platform to deploy its initial services. There are several major infrastructure models to consider:
- Public cloud providers (e.g., AWS, Azure, GCP)
- Self-hosted infrastructure (“own metal”)
- Open infrastructure platforms (e.g., OpenStack providers)

Each model comes with distinct trade-offs regarding cost, operational complexity, scalability, and alignment with our organisational values.

## Decision

We have chosen Amazon Web Services (AWS) as our primary cloud provider for this project.

Rationale:
1.	Existing Expertise:<br />
Our infrastructure development expertise already lies within AWS. The team is deeply familiar with AWS tooling, APIs, and best practices. This significantly reduces delivery risk and enables us to move at speed without incurring a skills gap.
2.	Time to Deployment:<br />
Given project timelines, our priority is to deliver a working system as efficiently as possible. AWS provides mature tooling, excellent documentation, and well-understood patterns within the team, which accelerates initial deployment and ongoing maintenance.
3.	Planned Future Migration to Open Infrastructure:<br />
As a Free/Libre and Open Source Software (FLOSS) organisation, we hold a strong belief in using and supporting open ecosystems. While AWS offers the most expedient route to launch, we consider it proper and aligned with our values to migrate to an OpenStack-based provider in the future, once the project has gained sufficient maturity and stability. Our architectural decisions are being made with this eventual transition in mind.
4.	Consideration of Self-Hosted Infrastructure<br />
We also evaluated the option of running our own infrastructure on physical hardware. While several high-profile organisations have recently migrated away from public cloud and back to on-premises or colocated environments, this is not currently feasible for us. We do not have the required physical space, staffing capacity, or up-front capital to support such a deployment. Our infrastructure needs to scale without introducing significant operational burden, and cloud services provide this elasticity without long-term hardware commitments.	
5. Comparison with Azure and GCP:
   - AWS is the clear market leader among public cloud providers, with the largest share of cloud infrastructure spend globally. This leadership translates into a mature and well-documented ecosystem, strong community support, a wide range of third-party integrations, and long-term platform stability. These attributes reduce risk and improve sustainability.
   - Azure is a platform we are technically familiar with; however, we do not depend on any other Microsoft services. In our experience, Azure is typically favoured by organisations that are already tightly integrated with the Microsoft ecosystem (e.g., Active Directory ("Entra"), Office365). As we have no such dependency, Azure does not offer material advantages for our use case.
   - GCP provides a modern developer experience and strong data analytics services, but its service model and ecosystem are less familiar to the team, and support resources are not as extensive.


## Consequences

- The infrastructure will be provisioned, automated, and maintained using AWS-native services and tooling.
- For now, we accept the trade-off of potential vendor lock-in in exchange for rapid initial progress.
- Migration to other providers in the future will require abstraction layers or infrastructure refactoring, which will be periodically reviewed as part of our platform evolution planning.
- Documentation and infrastructure-as-code will aim to preserve cloud-agnostic practices where feasible, particularly in OpenTofu module design and CI/CD workflows.
