# GitHub Variables

GitHub Variables is a convenient and cheap ($0) way to store CI/CD variables for GitHub Actions. Unfortunately, unlike
most other value stores, one must be an organisation administrator just to see the names of these variables, making it
difficult for those without that permission to use them.

As we are an open source organisation, these names are already publicly available in our code, so publishing them in one
convenient location poses no additional security risk. As such, here are the names of the GitHub Variables that are
available for use within GitHub Actions here in the `planzoco` repo.

## Organisation
### Secrets
- `AWS_ROLE_TO_ASSUME`
- `CLOUDFLARE_ACCOUNT_ID`
- `CLOUDFLARE_API_TOKEN`
- `CLOUDFLARE_TOKEN_MANAGER`
- `EVOTEUMBOT_APP_ID`
- `EVOTEUMBOT_APP_INSTALLATION_ID`
- `EVOTEUMBOT_APP_PEM_FILE`
- `EVOTEUMBOT_CLIENT_ID`
- `QUAY_ROBOT_PASSWORD`
- `QUAY_TOKEN`
- `TOFU_STATE_BUCKET_NAME`

### Variables
- `APP_NAME`
- `AWS_REGION`
- `BUILD_FLAGS`
- `FAIL_FAST`
- `LANGUAGE_VERSION`
- `QUAY_DOMAIN`
- `QUAY_OAUTH2_USERNAME`
- `QUAY_ROBOT_USERNAME`
- `QUAY_URL`
- `TARGET_ARCH`
- `TARGET_OS`
- `TERRAFORM_DOCS_IS_TOFU_COMPATIBLE`
- `TOFU_IS_JUNIT_COMPATIBLE`

## planzoco repo
If a repository variable and an organisation variable share the same name, the repository variable takes precedence.

### Secrets


### Variables
- `ARTEFACT_TYPE`
- `LANGUAGE`
- `LANGUAGE_VERSION`
- `QUAY_REPOSITORY_PATH`
- `QUAY_REPOSITORY_URL`
- `SOURCE_PATH`


---

> [!WARNING]  
> This file is controlled by OpenTofu in the [estate-repos](https://github.com/evoteum/estate-repos) repository.  
>  
> Manual changes will be **overwritten** within 24 hours.
