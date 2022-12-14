pid_file = "/tmp/pidfile"

auto_auth {
  method {
    type   = "approle"
    config = {
      role_id_file_path = "/vault/config/approle/course-roleid"
      secret_id_file_path = "/vault/config/approle/course-secretid"
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
  source = "/vault/config/appconf/course.tmpl"
  destination = "/vault/config/appconf/course.yaml.tmp"
  exec {
    command = ["dd", "if=/vault/config/appconf/course.yaml.tmp", "of=/vault/config/appconf/course.yaml" ]
    timeout = "5s"
  }
}
