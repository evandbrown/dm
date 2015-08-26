# DM
An opinionated tool for using the Deployment Manager service from Google Cloud Platform. Makes common tasks easy with basic commands. `dm` assumes that the directory you execute it from should be associated with a particular deployment.

Multiple deployments are supported. Each deployment has an entry in the `.dm.toml` config file created by `dm`. If multiple deployments are present, the `--name` param must be used to indicate which deployment to perform the command on.

* `dm deploy` creates a new deployment using `config.yaml` in your pwd.
* `dm ls` lists available deployments in a directory
* `dm rm` deletes the deployment.
* `dm stat` describes the deployment defined in `.dm`.

## Local parameters
You can provide parameters when a deployment is created or update. The convention is to define them in `vars.yaml` in a dict named `variables`, like:


```yaml
variables:
  image: "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20150710"
  num_servers: "2"
```

This creates two vars - `image` and `num_servers`, each with default values. They can be used in the Deployment Manager config.yaml using the {{ var \`var_name\` }} syntax. This is modled after [Packer's](https://packer.io) support for variables:

```yaml
---
imports:
  - path: config.jinja
resources:
  -
    name: loader
    type: config.jinja
    properties:
      image: "{{var `image`}}"
      num_servers: "{{ var `num_servers`}}"
      datestamp: "1728992022"
      test_duration: 60m
```

Variables (including those with default values defined in `vars.yaml`) can be specified or overridden at deploy time:

```bash
dm deploy              \
  --var image=debian-8 \
  --var num_servers=2
```

Variables can also receive a default value from the environment. In `vars.yaml` use the {{env `ENV_VAR_NAME`}} function:

```yaml
variables:
  image: "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20150710"
  num_servers: "{{env `NUM_SERVERS`}}"
```

Variables defined in `vars.yaml` must have a value at deploy time. The following is invalid:

`vars.yaml`:

```yaml
variables:
  image: "https://www.googleapis.com/compute/v1/projects/debian-cloud/global/images/debian-8-jessie-v20150710"
  num_servers: ""
```

```shell
dm deploy -p evandbrown17
FATA[0000] An error occurred                             error=variable num_servers is empty
```

## Git integration
Future support for a git-based workflow that associates tags and branches with deployments.

## Common workflows
TODO

## Contributing
Contributions are so totally welcome. I suggest opening an Issue first to describe what you'd like to contribute and make sure it's not already being worked on.

To get the project working and submit a PR, you should:

1. Fork it
1. Clone it to your workstation and `cd` to it
1. `go get ./...` to download dependencies
1. Create a feature branch: `git checkout -b your-new-feature`
1. Some day there will be tests that you would run at this step :)
1. `go install` to build and install to `$GOPATH/bin`
1. Do some user acceptance tests
1. Push your branch and create a PR! 
