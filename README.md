# terraform-provider-regru
Terraform provider for reg.ru cloud servers.

# Building from source


resource "regru_server" "kvm_server" {
    name = "test"
    size = "cloud-1"
    image = "ubuntu-18-04-amd64"
    token = "your api token"
}
