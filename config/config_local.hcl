server {
  bind_addr = ":8080"
  // one of "DEBUG", "INFO", "WARN", "ERROR"
  log_level = "INFO"
  // in seconds
  ctx_timeout = 5
}

api {
  postgres {
    dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
  }
}
