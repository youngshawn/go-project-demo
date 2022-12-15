pid_file = "/tmp/pidfile"

vault {
  address = "http://vault-course:8200"
  retry {
    num_retries = 5
  }
}

listener "tcp" {
  address = "127.0.0.1:8201"
  tls_disable = true
}

cache {
  use_auto_auth_token = true
}

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
