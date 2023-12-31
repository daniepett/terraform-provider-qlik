---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qlik_source_entities Data Source - terraform-provider-qlik"
subcategory: ""
description: |-
  
---

# qlik_source_entities (Data Source)



## Example Usage

```terraform
data "qlik_source_entities" "example" {
  source_connection_id = "source-connection-id"
  database             = ""
  table_pattern        = "%"
  schema_pattern       = "%"
  entity_type          = "TABLE"
  project_id           = "project-id"
  app_id               = "data-app-id"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_id` (String)
- `database` (String)
- `entity_type` (String)
- `project_id` (String)
- `schema_pattern` (String)
- `source_connection_id` (String)
- `table_pattern` (String)

### Read-Only

- `entities` (Attributes List) (see [below for nested schema](#nestedatt--entities))

<a id="nestedatt--entities"></a>
### Nested Schema for `entities`

Read-Only:

- `data_app_id` (String)
- `database` (String)
- `id` (String)
- `name` (String)
- `project_id` (String)
- `schema` (String)
- `type` (String)
