
# grpc

gRPC service and utility to check load balancing and connection stability

## Build the service

```sh
make build
```

## Deploy the service

```sh
make deploy
```

## Run the utility with simple ping

```sh
./bin/darwin-amd64/grpc client ping --count 600 --config config/test.conf
```

## Run the utility with reverse ping

```sh
./bin/darwin-amd64/grpc client reverse --count 600 --config config/test.conf
```

## Scale services

During the ping process scale services from 1 to N and see how it is going

```sh
kubectl scale --replicas 3 -n test deploy/grpc
kubectl scale --replicas 6 -n test deploy/grpc
kubectl scale --replicas 1 -n test deploy/grpc
...
```

## Other make commands

* `all` - run default complete set of commands (build the service)
* `vendor` - import all vendors (using dep)
* `compile` - build the service binary
* `certs` - download latests certs from an alpine image and prepare it for service container
* `build` - build container image
* `push` - push an image in docker registry
* `run` - build and run the service
* `logs` - show service logs from container
* `deploy` - deployment of the service into Kubernetes environment
* `charts` - validate helm templates (charts)
* `test` - run unit tests
* `cover` - show testing coverage for packages
* `fmt` - format Go packages with go fmt
* `lint` - use set of linters ( ~ 20) to check the service code
* `stop` - stop running container
* `start` - start existing container (if it was stopped before)
* `rm` - remove stopped container
* `version` - add next major/minor/patch version
* `clean` - remove binary and running container
* `bootstrap` - check and setup if something from utilities is not exist

## Versioned changes

All changes in the project described in [changelog](docs/CHANGELOG.md)

_Generated using ([Caldera boilerplate](https://github.com/takama/caldera))_
