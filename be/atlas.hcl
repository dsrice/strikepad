env "dev" {
  src = "file://schema.sql"
  dev = "postgres://postgres:password@localhost:5432/strikepad?sslmode=disable"
  url = "postgres://postgres:password@localhost:5432/strikepad?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "test" {
  src = "file://schema.sql"
  url = "postgres://postgres:password@localhost:5432/strikepad_test?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
}

env "production" {
  src = "file://schema.sql"
  url = env("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}