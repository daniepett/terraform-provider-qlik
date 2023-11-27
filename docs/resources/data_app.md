---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qlik_data_app Resource - terraform-provider-qlik"
subcategory: ""
description: |-
  
---

# qlik_data_app (Resource)



## Example Usage

```terraform
resource "qlik_data_app" "example" {
  name        = "Example"
  description = "Description"
  type        = "LANDING_SAAS_MANAGED"
  project_id  = "data-project-id"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `project_id` (String)
- `type` (String)

### Optional

- `description` (String)

### Read-Only

- `id` (String) The ID of this resource.