# Why Terraform?

CloudLab’s Portal and XML-RPC are powerful, but not designed for **reproducible IaC**. Terraform provides:

!!! tip "Declarative desired state"
    You describe **what** you want; Terraform computes **how** to converge.

!!! info "Versioning and review"
    Store infrastructure alongside code; PRs become infra change reviews.

!!! example "Automation & CI/CD"
    Provision teaching labs or research clusters from a single pipeline run.

!!! note "Hybrid workflows"
    Combine CloudLab with AWS/Azure/GCP providers in one stack.
