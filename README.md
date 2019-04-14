# go-mod-dep-source-finder [![](https://hush-house.pivotal.io/api/v1/teams/main/pipelines/go-mod-license-finder/badge)](https://hush-house.pivotal.io/teams/main/pipelines/go-mod-license-finder)

Given dependencies from `go list`, retrieves the URL of each dependency matching the versions.


## Example

```sh
# Get the URL of a single dependency
#
go-mod-dep-source-finder 'golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2' | jq
{
  "original": "golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2",
  "discovered": "https://go.googlesource.com/crypto/+/c2843e01d9a2/"
}

# Get the URL of all of your dependencies
#
go list -m all | tail -n +2 | go-mod-license-finder -
...
```

## Install

### Latest using `go`

Having `$GOPATH/bin` in your `$PATH`:

```sh
go get -u github.com/cirocosta/go-mod-dep-source-finder
```


### Stable

Binaries are distributed through GitHub releases.

TODO


### Docker

Container images are continuously shipped to DockerHub:

```sh
docker run \
  cirocosta/go-mod-dep-source-finder \
  'golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2'
```


## LICENSE

Apache 2 - see [./LICENSE](LICENSE).

