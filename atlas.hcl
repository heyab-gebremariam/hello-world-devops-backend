# Declare the input variable expected from the command line
variable "db_url" {
  type        = string
  description = "The URL of the development database"
  // default = "postgres://user:pass@host:port/db?sslmode=disable"
}

data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./atlas/main.go",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = var.db_url

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}