
# grpc

gRPC service and utility to check load balancing and connection stability

## Prepare service environment and configuration

- Create a new project and GKE cluster in Google Cloud
- Setup Contour balancer according to [example](https://github.com/projectcontour/contour/tree/master/examples/contour)
- Note somewhere the IP of Load balancer
- Change default GKE values in `Makefile` on your own

```Makefile
GKE_PROJECT_ID ?= your-project-id
GKE_PROJECT_ZONE ?= europe-west1-b
GKE_CLUSTER_NAME ?= your-cluster-name
```

- Choose domain host name for the service and point your DNS record on GKE Contour balancer using IP above
- Make the certificates for the domain name (Let's encrypt as option)
- Test to create in dry run  `grpc-service-tls` in corresponded Kubernetes environment (ex: `test`)

```sh
kubectl create secret tls grpc-service-tls --key key.pem --cert cert.pem --dry-run -o yaml
```

- Create `grpc-service-tls` in corresponded Kubernetes environment (ex: `test`)

```sh
kubectl create secret tls grpc-service-tls --key key.pem --cert cert.pem
```

- Change domain name in `.helm/values-test.yaml` or any other `values-name.yaml` to your own

```yaml
  ## Ingress route hosts
  ##
  hosts:
    ## gRPC service host
    - name: grpc
      host: grpc.your-test-domain.net
...
  ## Client connection to the service
  ##
  client:
    host: grpc.your-test-domain.net
```

- Change default namespace in `Makefile`

```Makefile
# Namespace: dev, prod, username ...
NAMESPACE ?= test
```

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

- `all` - run default complete set of commands (build the service)
- `vendor` - import all vendors (using dep)
- `compile` - build the service binary
- `certs` - download latests certs from an alpine image and prepare it for service container
- `build` - build container image
- `push` - push an image in docker registry
- `run` - build and run the service
- `logs` - show service logs from container
- `deploy` - deployment of the service into Kubernetes environment
- `charts` - validate helm templates (charts)
- `test` - run unit tests
- `cover` - show testing coverage for packages
- `fmt` - format Go packages with go fmt
- `lint` - use set of linters ( ~ 20) to check the service code
- `stop` - stop running container
- `start` - start existing container (if it was stopped before)
- `rm` - remove stopped container
- `version` - add next major/minor/patch version
- `clean` - remove binary and running container
- `bootstrap` - check and setup if something from utilities is not exist

## Versioned changes

All changes in the project described in [changelog](docs/CHANGELOG.md)

_Generated using ([Caldera boilerplate](https://github.com/takama/caldera))_
