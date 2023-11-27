resource "qlik_data_connection" "example" {
  name       = "example"
  space_id   = "space-id"
  gateway_id = "data-gateway-id"
  type       = "reptgt_qdisnowflake"
  connection_parameters = {
  }

}