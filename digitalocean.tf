terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "1.22.2"
    }
  }
}

variable "do_token" {}

variable "ssh_key_path" {}

variable "vm_count" {}

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
  count = var.vm_count

  image              = "ubuntu-20-04-x64"
  name               = "bigwpeer-${count.index}"
  region             = "sgp1"
  size               = "s-1vcpu-1gb"
  private_networking = true
  ssh_keys = [
    data.digitalocean_ssh_key.UnseenUniversity.id,
    data.digitalocean_ssh_key.Vetinari.id
  ]

  connection {
    host    = self.ipv4_address
    user    = "root"
    type    = "ssh"
    agent   = true
    timeout = "2m"
  }

  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",

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
      "cd bigw-voting",
      "go mod download all"
    ]
  }
}

output "ip" {
  value = join("\n", digitalocean_droplet.bigwpeer[*].ipv4_address)
}
