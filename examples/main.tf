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

data "cribl_system" "build" {}

output "cribl_build" {
  value = data.cribl_system.build
}