// creds from .secret file
provider "aws" {
  region = ""
  access_key = ""
  secret_key = "" 
}

resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name           = "UsersTable"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "UserID"

  attribute {
    name = "UserID"
    type = "S"
  }
}