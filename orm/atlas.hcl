env "postgres" {
  src = "ent://ent/schema"
  url = getenv("POSTGRES_DSN")
  dev = getenv("POSTGRES_DSN")
  schemas = ["public"]
  migration {
    dir = "file://ent/migrate/migrations"
    revisions_schema = "public"
  }
}
