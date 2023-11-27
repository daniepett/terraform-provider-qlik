resource "qlik_data_app" "example" {
  name        = "Example"
  description = "Description"
  type        = "LANDING_SAAS_MANAGED"
  project_id  = "data-project-id"
}