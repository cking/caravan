# caravan
caravan is a GitOps operator for Hashicorp Nomad.

## How to use

### Setting up caravan

#### Environment variables
```
NOMAD_ADDR - Required to overide the default of http://127.0.0.1:4646.
NOMAD_TOKEN - Required with ACLs enabled.
NOMAD_CACERT - Required with TLS enabled.
NOMAD_CLIENT_CERT - Required with TLS enabled.
NOMAD_CLIENT_KEY - Required with TLS enabled.
GIT_REPO - get repository to checkout
GIT_BRANCH - Branch to use, default is 'main'
GIT_PATH - only index files in this directory and its subdirectories
CARAVAN_INTERVAL - Run caravan in these intervals (defaults to 1 Minute)
```

### Usage
Set the environment variables by a tool of your choice, or place a `.env` file in the working directory.
Then just run caravan:
```bash
./caravan
```

## Run as Nomad job
```yaml
job "caravan" {
  datacenters = ["dc1"]

  group "caravan" {
    count = 1
    task "caravan" {
      driver = "exec"
      config {
        command = "caravan"
      }
  
      env {
        GIT_REPO = "https://github.com/user/some_nomad_job_file_repo"
      }
  
      artifact {
        source      = "https://github.com/cking/caravan/releases/download/v0.0.1/caravan_0.0.1_linux_amd64.tar.gz"
        destination = "local"
        mode        = "any"
      }
    }
  }
}
```

## Acknowledgement
[nomad-gitops-operator](https://github.com/jonasvinther/nomad-gitops-operator) by @jonasvinthjer