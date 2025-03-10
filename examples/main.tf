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
  pipelines = [
    {
      id          = "test01"
      description = "Foo pipeline"
      timeout_ms  = 3000
      tags        = ["foo"]
      output      = "default"
    },
    {
      id          = "test02"
      description = "Foo Bar pipeline"
      timeout_ms  = 60000
      tags        = ["foo", "bar"]
      output      = "default"
    }
  ]
}

#resource "cribl_route" "route1" {
#  name        = "First Route"
#  pipeline_id = cribl_api.example.pipelines[0].id

# This implicit dependency ensures the pipeline is created before the route
#  depends_on = [cribl_api.example]
#}
