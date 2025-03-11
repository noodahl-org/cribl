terraform {
  required_providers {
    cribl = {
      source = "noodahl-org/cribl"
    }
  }
}

provider "cribl" {
  username = "admin"
  password = "changeme"
  base_url = "http://localhost:9000"
}

resource "cribl_pipeline" "example" {
  id          = "test01"
  description = "Foo pipeline"
  timeout_ms  = 3000
  tags        = ["foo"]
  output      = "default"
}

resource "cribl_output_s3" "example" {
  id                        = "output_example"
  type                      = "s3"
  default_id                = "output_example"
  description               = "S3 output that depends on test01 pipeline"
  bucket                    = "my-cribl-output-bucket"
  region                    = "us-west-2"
  stage_path                = "/tmp/cribl/s3-staging"
  dest_path                 = "cribl-outputs/${formatdate("YYYY/MM/DD", timestamp())}"
  aws_authentication_method = "auto" # Can be "auto", "secret", or "manual"

  # Authentication - choose one method
  # Option 1: Using AWS roles (recommended for production)
  # No additional fields needed if using IAM roles
  # OR
  # aws_secret = "my-aws-credentials"
  # OR
  # aws_api_key = "YOUR_ACCESS_KEY_ID"
  # aws_secret_key = "YOUR_SECRET_ACCESS_KEY"

  # S3 file management settings
  format                 = "json"
  compress               = "gzip"
  add_id_to_stage_path   = true
  max_file_size_mb       = 50
  max_file_open_time_sec = 60
  max_file_idle_time_sec = 30

  # Link to pipeline
  pipeline = cribl_pipeline.example.id

  # Advanced settings
  deadletter_enabled = true
  deadletter_path    = "/tmp/cribl/deadletter"

  depends_on = [
    cribl_pipeline.example
  ]
}

