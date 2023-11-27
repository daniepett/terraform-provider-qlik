resource "qlik_space" "example" {
  name        = "Some Name"
  description = "I describe the space"
  type        = "data"
}


resource "qlik_data_project" "example" {
  name               = "example"
  space_id           = qlik_spaces.example.id
  lakehouse_type     = "SNOWFLAKE"
  type               = "DATA_PIPELINE"
  storage_connection = "storage-connection-id"
  description        = "This is a desc"
}