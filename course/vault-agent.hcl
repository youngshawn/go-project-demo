pid_file = "/tmp/pidfile"

auto_auth {
  method {
    type   = "approle"
    config = {
      role_id_file_path = "/tmp/.course-roleid"
      secret_id_file_path = "/tmp/.course-secretid"
      remove_secret_id_file_after_reading = false
    }
  }

  sink "file" {
      config = {
          path = "/tmp/.vault-token-via-agent"
      }
  }
}

vault {
  address = "http://127.0.0.1:8200"
}

template_config {
  exit_on_retry_failure = true
  
}

template {
  error_on_missing_key = true
  source = "./course.tmpl"
  destination = "./course.yaml.tmp"
  exec {
    command = ["cp", "./course.yaml.tmp", "./course.yaml" ]
    timeout = "5s"
  }
}
