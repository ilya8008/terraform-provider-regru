# terraform-provider-regru
Terraform provider for reg.ru cloud servers.

# Binaries

Download binaries from https://github.com/ilya8008/terraform-provider-regru/releases/

# Building from source

```
go get github.com/ilya8008/terraform-provider-regru
cd $GOPATH/src/github.com/dmacvicar/terraform-provider-regru
go build -o terraform-provider-regru
```
# Installing

Place `terraform-provider-regru` file in `~/terraform.d/plugins` folder.


# main.tf example

```
resource "regru_server" "kvm_server" {
    name = "test"
    size = "cloud-1"
    image = "ubuntu-18-04-amd64"
    token = "your api token"
}
```
