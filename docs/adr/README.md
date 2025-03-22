# Architecture Decision Records (ADRs)

This directory contains **Architecture Decision Records (ADRs)**, which document important technical and architectural decisions made in this project.

## What is an ADR?

An **Architecture Decision Record (ADR)** is a document that captures key technical decisions and their rationale. Each ADR follows a consistent format to ensure clarity and traceability.

## ADR Format

Every ADR should include the following sections:

1. **Context** – What problem or situation led to this decision?
2. **Decision** – What was decided and why?
3. **Status** – The current state of this decision (e.g., Proposed, Accepted, Deprecated, Superseded).
4. **Consequences** – What trade-offs and impacts does this decision have?

## Why Use ADRs?

- **Avoid repeating past mistakes** – Future contributors can see why certain options were rejected.
- **Aids reassessment** - Allows stakeholders to more easily identify if the situation is still relevant.
- **Improve knowledge sharing** – New team members can quickly understand past decisions.
- **Enable informed decisions** – Understanding past choices helps guide future architecture.

## How to Create a New ADR

### Prerequisites

Install [`adr-tools`](https://github.com/npryce/adr-tools/blob/master/INSTALL.md)

### ADR Creation
- New ADR 
  ```shell
  adr new Amazing Decision
  ```
- ADR that supersedes a previous one (ADR 9, for example)
  ```shell
   adr new -s 9 Use Rust for performance-critical functionality
  ```
  
Then just push your ADR.

---

> [!WARNING]  
> This file is controlled by OpenTofu in the [estate-repos](https://github.com/evoteum/estate-repos) repository.  
>  
> Manual changes will be **overwritten** within 24 hours.
