[//]: # (STANDARD README)
[//]: # (https://github.com/RichardLitt/standard-readme)
[//]: # (----------------------------------------------)
[//]: # (Uncomment optional sections as required)
[//]: # (----------------------------------------------)

[//]: # (Title)
[//]: # (Match repository name)
[//]: # (REQUIRED)

# planzoco

[//]: # (Banner)
[//]: # (OPTIONAL)
[//]: # (Must not have its own title)
[//]: # (Must link to local image in current repository)


[//]: # (Badges)
[//]: # (OPTIONAL)
[//]: # (Must not have its own title)


[//]: # (Short description)
[//]: # (REQUIRED)
[//]: # (An overview of the intentions of this repo)
[//]: # (Must not have its own title)
[//]: # (Must be less than 120 characters)
[//]: # (Must match GitHub's description)

Every question finds its best answer

[//]: # (Long Description)
[//]: # (OPTIONAL)
[//]: # (Must not have its own title)
[//]: # (A detailed description of the repo)



[//]: # (Keep this note to help people understand how to configure this repo.)
The configuration of this repo is managed by OpenTofu in [estate-repos](https://github.com/evoteum/estate-repos).

## Table of Contents

[//]: # (REQUIRED)
[//]: # (Delete as appropriate)

[//]: # (TOCGEN_TABLE_OF_CONTENTS_START)

1. [Security](#security)
1. [Background](#background)
1. [Install](#install)
1. [Usage](#usage)
1. [Infrastructure](#infrastructure)
1. [Contributing](#contributing)
1. [License](#license)



[//]: # (TOCGEN_TABLE_OF_CONTENTS_END)

## Security

[//]: # (OPTIONAL)
[//]: # (May go here if it is important to highlight security concerns.)

This service is public. There are no accounts or passwords, voting asks for a display name ("what does this group
know you as?") purely so one person's vote counts once and can be changed, not to authenticate anyone. That name is
stored (in a cookie, and against each vote) but not verified, so, like everything else in an event, it should be
assumed to be public and self-reported, not a real identity check. Don't be a dick to your mates.

## Background
[//]: # (OPTIONAL)
[//]: # (Explain the motivation and abstract dependencies for this repo)

Planzoco exists because group decisions made over chat get lost. A WhatsApp group planning a trip ends up with the
question, the suggestions, and the votes scattered across dozens of messages. By the time you need an answer, nobody
can reconstruct what the group actually agreed to.

Planzoco takes some inspiration from [poll.ly](https://poll.ly/), in that any participant can suggest an option and
everyone votes. Its point of difference is that one event is one link covering *every* open question, not one poll
per question. Organising an excursion with friends, or a family get together, usually means agreeing on several
things at once, such as,
- Where are we going?
- What are we doing?
- How are we getting there?

Planzoco lets the group answer all of them from a single shared link: no separate poll to create and share for each
question, no sign-up, and every open decision for the event lives in one place instead of buried in chat history.

## Install

[//]: # (Explain how to install the thing.)
[//]: # (OPTIONAL IF documentation repo)
[//]: # (ELSE REQUIRED)

In production, planzoco runs on our bare metal Kubernetes cluster and is deployed automatically by ArgoCD; see
[Infrastructure](#infrastructure). There's nothing to install manually.

To run it locally:

```sh
cd go/planzoco
export DATABASE_URL=postgres://planzoco:planzoco@localhost:5432/planzoco?sslmode=disable
go run .
```

Planzoco creates its own schema against that database on startup, so there's no separate migration step: just point
`DATABASE_URL` at any empty Postgres database. See [env.list](go/planzoco/env.list) for other supported environment
variables.

## Usage
[//]: # (REQUIRED)
[//]: # (Explain what the thing does. Use screenshots and/or videos.)

Planzoco has no accounts, anyone with the link can create, suggest, and vote.

1. **Create an event** from the homepage (e.g. "Saturday hike"); this gives you the one link you'll share with the
   group.
2. **Add questions** on the event page for each decision the group needs to make. The event page lists every
   question alongside its current leading answer, so it doubles as a live summary of what's settled and what's still
   open. A question with no options yet invites you to "Add options"; once it has at least one, that becomes "Vote!".
3. **Suggest options and vote** by opening a question; anyone can add an option, and votes update the leading
   answer(s) on the event page immediately. The first time you vote you're asked for a name (not a password, just
   what the group knows you as) so your vote counts once and you can change your mind later; it's remembered after
   that. Ties show as joint leaders until another vote breaks them.
4. **Share the event link** shown at the bottom of the event page: that single link covers every question in the
   event, so the group never needs a separate poll per decision.


[//]: # (Extra sections)
[//]: # (OPTIONAL)
[//]: # (This should not be called "Extra Sections".)
[//]: # (This is a space for ≥0 sections to be included,)
[//]: # (each of which must have their own titles.)

## Infrastructure

Planzoco originally ran on Evoteum's AWS ECS "golden path", provisioned by the shared `platform_service` OpenTofu
module, with DynamoDB for storage. It was migrated to our bare metal Kubernetes cluster: the application is packaged
as a Helm chart in [`chart/`](chart/), deployed via ArgoCD (registered in
[kubernetes-lab-services](https://github.com/evoteum/kubernetes-lab-services)), and backed by Postgres managed
in-cluster by [CloudNativePG](https://cloudnative-pg.io/) instead of DynamoDB. See
[ADR 0005](docs/adr/0005-migrate-to-bare-metal-kubernetes-and-cloudnativepg.md) for the full rationale.

The AWS infrastructure that previously backed it, in [`tofu/production`](tofu/production), has been wound down.

## EaC

We aim to follow the "Everything as Code" practice wherever possible. We prefer to define everything in code,
including things that might not "usually" be defined in code. For more info, see the estate-repos repo, where all our
repositories are defined and maintained.


[//]: # (## API)
[//]: # (OPTIONAL)
[//]: # (Describe exported functions and objects)



[//]: # (## Maintainers)
[//]: # (OPTIONAL)
[//]: # (List maintainers for this repository)
[//]: # (along with one way of contacting them - GitHub link or email.)



[//]: # (## Thanks)
[//]: # (OPTIONAL)
[//]: # (State anyone or anything that significantly)
[//]: # (helped with the development of this project)



## Contributing
[//]: # (REQUIRED)
If you need any help, please log an issue and one of our team will get back to you.

PRs are welcome.


## License
[//]: # (REQUIRED)

All our code is licenced under the AGPL-3.0. See [LICENSE](LICENSE) for more information.
