data "qlik_source_entities" "example" {
  source_connection_id = "source-connection-id"
  database             = ""
  table_pattern        = "%"
  schema_pattern       = "%"
  entity_type          = "TABLE"
  project_id           = "project-id"
  app_id               = "data-app-id"
}