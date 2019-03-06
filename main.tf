#resource "regru_ssh" "ssh_key" {
#    name = "${var.ssh_name}"
#    public_key = "${var.public_key}"
#    token = "${var.token}"
#}


resource "regru_server" "kvm_server" {
    name = "${var.name}-${count.index}"
    size = "${var.size}"
    image = "${var.image}"
    token = "${var.token}"
    count = "${var.count}"
}

#resource "regru_server" "kvm_server" {
#    name = "${var.name}"
#    size = "${var.size}"
#    image = "${var.image}"
#    token = "${var.token}"
#    count = "${var.count}"
#}

