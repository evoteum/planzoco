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

1. [Security](#security)
1. [Background](#background)
1. [Install](#install)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

## Security

[//]: # (OPTIONAL)
[//]: # (May go here if it is important to highlight security concerns.)

This service is public. No login details are currently stored so customers should assume everything that they
write is public.

## Background
[//]: # (OPTIONAL)
[//]: # (Explain the motivation and abstract dependencies for this repo)

Planzoco takes some inspiration from [poll.ly](https://poll.ly/), in that all participants can suggest and vote on an item. Where
planzoco takes it further is the inclusion of multiple questions. When organising an excursion with friends, or a
family get together, there are many questions that must be agreed apon, such as,
- Where are we going?
- What are we doing?
- How are we getting there?

Planzoco allows all participants to suggest an option, then everyone gets to vote on it.

## Install

[//]: # (Explain how to install the thing.)
[//]: # (OPTIONAL IF documentation repo)
[//]: # (ELSE REQUIRED)

Run the container.

## Usage
[//]: # (REQUIRED)
[//]: # (Explain what the thing does. Use screenshots and/or videos.)

Run the container.


[//]: # (Extra sections)
[//]: # (OPTIONAL)
[//]: # (This should not be called "Extra Sections".)
[//]: # (This is a space for ≥0 sections to be included,)
[//]: # (each of which must have their own titles.)

## EaC

You may notice that there is less OpenTofu (the FLOSS fork of Terraform) code here than you might have expected.
Evoteum has an ECS base deployment platform, so the only module that each repos needs is the platform service module.
This provides a truly repeatable "golden path" deployment to ensure that everyone can focus on just getting code to
customers faster.

Other modules are defined in the tofu-modules repo, while the central platform infrastructure is defined in the
platform-infrastructure module.

We also aim to follow the "Everything as Code" practice wherever possible. We prefer to define everything in code,
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
