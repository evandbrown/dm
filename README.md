# DM
An opinionated tool for using the Deployment Manager service from Google Cloud Platform. Makes common tasks easy with basic commands. `dm` assumes that the directory you execute it from should be associated with a particular deployment.

Multiple deployments are supported. Each deployment has an entry in the `.dm.toml` config file created by `dm`. If multiple deployments are present, the `--name` param must be used to indicate which deployment to perform the command on.

* `dm deploy` creates a new deployment using `config.yaml` in your pwd.
* `dm ls` lists available deployments in a directory
* `dm rm` deletes the deployment.
* `dm stat` describes the deployment defined in `.dm`.

## Local parameters
Future support for using local paramters in a configuration.

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
