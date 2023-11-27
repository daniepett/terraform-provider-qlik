---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qlik_space Resource - terraform-provider-qlik"
subcategory: ""
description: |-
  
---

# qlik_space (Resource)



## Example Usage

```terraform
resource "qlik_space" "example" {
  name        = "Some Name"
  description = "I describe the space"
  type        = "data"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `type` (String)

### Optional

- `description` (String)

### Read-Only

- `id` (String) The ID of this resource.
- `owner_id` (String)