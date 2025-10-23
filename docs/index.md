# CloudLab Terraform Provider
[![GitHub Tag](https://img.shields.io/github/v/tag/CSC478-WCU/terraform-provider-cloudlab?logo=terraform&label=latest&color=%237B42BC)](https://github.com/CSC478-WCU/terraform-provider-cloudlab/releases)
[![Terraform Provider Downloads](https://img.shields.io/terraform/provider/dt/846963?logo=terraform&label=Registry%20downloads&color=%237B42BC)](https://registry.terraform.io/providers/CSC478-WCU/cloudlab)
[![GitHub](https://img.shields.io/badge/GitHub-Repository-181717?logo=github)](https://github.com/CSC478-WCU/terraform-provider-cloudlab)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Tyler%20Geiger-0A66C2?logo=linkedin)](https://www.linkedin.com/in/tyler-geiger)

---

!!! example "Provider Information"

    **Author**: [Tyler Geiger](https://www.linkedin.com/in/tyler-geiger) — West Chester University  
    
    **Thanks to my professor:** [Dr. Linh Ngo](https://www.cs.wcupa.edu/LNGO).

    **Inspiration**: UCY-COAST — [terraform-provider-cloudlab](https://github.com/ucy-coast/terraform-provider-cloudlab)  

    Built as an independent provider to support our main full-stack Cloud Computing project *CampusConnect* at West Chester University.
    The goal was to enable reproducible Infrastructure-as-Code (IaC) deployments on CloudLab for academic research and teaching.
    Acknowledgments to the Terraform and CloudLab communities for foundational tooling and research infrastructure.


---

The **CloudLab Terraform Provider** bridges academic research environments with modern Infrastructure-as-Code (IaC) workflows.  
It allows you to define, version, and automate **CloudLab experiments** using Terraform—just like provisioning AWS or Azure infrastructure.

---

### Quick Links

- [Get Started →](usage/getting-started.md)
- [Authentication with PEM →](usage/auth.md)
- [CI/CD With Github Actions Example →](usage/github-actions.md)
- [Provider Architecture →](internals/architecture.md)

---

### Overview

CloudLab is a powerful academic testbed for cloud computing research, offering bare-metal and virtualized environments across multiple university aggregates.  
This provider brings CloudLab into the Terraform ecosystem, enabling reproducible, version-controlled infrastructure deployments.
