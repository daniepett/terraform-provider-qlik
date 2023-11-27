data "qlik_source_entities" "example" {
  source_connection_id = "source-connection-id"
  database             = ""
  table_pattern        = "%"
  schema_pattern       = "%"
  entity_type          = "TABLE"
  project_id           = "project-id"
  app_id               = "data-app-id"
}

resource "qlik_data_app_source_selection" "example" {
  project_id           = "project-id"
  app_id               = "data-app-id"
  source_connection_id = "source-connection-id"
  source_selection     = data.qlik_source_entities.example.entities
}