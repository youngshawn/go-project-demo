pid_file = "/tmp/pidfile2"

auto_auth {
  method {
    type   = "approle"
    config = {
      role_id_file_path = "/vault/config/course-roleid"
      secret_id_file_path = "/vault/config/course-secretid"
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
  address = "http://vault-course:8200"
}

template_config {
  exit_on_retry_failure = true
  
}

template {
  error_on_missing_key = true
  source = "/vault/config/course.tmpl"
  destination = "/vault/file/course.yaml.tmp"
  exec {
    command = ["cp", "/vault/file/course.yaml.tmp", "/vault/file/course.yaml" ]
    timeout = "5s"
  }
}
