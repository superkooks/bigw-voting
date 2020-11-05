terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "1.22.2"
    }
  }
}

variable "do_token" {}

variable "ssh_key_path" {}

provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_ssh_key" "UnseenUniversity" {
  name = "UnseenUniversity"
}

data "digitalocean_ssh_key" "Vetinari" {
  name = "Vetinari"
}


resource "digitalocean_droplet" "bigwpeer" {
    count = 2

    image = "ubuntu-20-04-x64"
    name = "bigwpeer-${count.index}"
    region = "sgp1"
    size = "s-1vcpu-1gb"
    private_networking = true
    ssh_keys = [
        data.digitalocean_ssh_key.UnseenUniversity.id,
        data.digitalocean_ssh_key.Vetinari.id
    ]

    connection {
        host = self.ipv4_address
        user = "root"
        type = "ssh"
        private_key = file(var.ssh_key_path)
        timeout = "2m"
    }

    provisioner "remote-exec" {
        inline = [
            "export PATH=$PATH:/usr/bin",
            
            # Install git
            "sudo apt-get update",
            "sudo apt-get -y install git",

            # Install golang
            "wget https://golang.org/dl/go1.15.3.linux-amd64.tar.gz",
            "tar -C /usr/local -xzf go1.15.3.linux-amd64.tar.gz",
            "export PATH=$PATH:/usr/local/go/bin",

            # Setup env, download code
            "export GOPATH=$HOME/go/",
            "echo \"GOPATH=$HOME/go/\" >> ~/.bashrc",
            "echo \"PATH=$PATH:/usr/bin\" >> ~/.bashrc",
            "mkdir go",
            "cd go",
            "mkdir src",
            "cd src",
            "git clone https://github.com/SuperKooks/bigw-voting.git",
            "go get github.com/huin/goupnp/dcps/internetgateway2",
            "go get github.com/spf13/pflag",
            "go get github.com/gdamore/tcell",
            "go get gitlab.com/tslocum/cview"
        ]
    }
}

output "vm_ip0" {
    value = digitalocean_droplet.bigwpeer[0].ipv4_address
}

output "vm_ip1" {
    value = digitalocean_droplet.bigwpeer[1].ipv4_address
}

