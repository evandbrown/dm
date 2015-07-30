# DM
An opinionated tool for using the Deployment Manager service from Google Cloud Platform. Makes common tasks easy with basic commands. `dm` assumes that the directory you execute it from should be associated with a particular deployment.

`dm create` creates a new deployment using `config.yaml` in your pwd. It writes a `.dm` file with the resulting deployment and project names.
`dm update` updates the deployment defined in `.dm`.
`dm delete` deletes the deployment defind in `.dm`. Optionally destroys Google Cloud Storage buckets and their objects.
`dm status` describes the deployment defined in `.dm`.
`dm resources` lists the resources created by the deployment defined in `.dm`.

## Local parameters

## Git integration

## Common workflows

TODO:
1. Write toml on create
2. Implement delete
3. Support parameters
