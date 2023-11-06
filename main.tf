/*
Based on https://github.com/futurice/terraform-examples/blob/master/google_cloud/minecraft/main.tf
*/

# We require a project to be provided upfront
# Create a project at https://cloud.google.com/
# Make note of the project ID
# We need a storage bucket created upfront too to store the terraform state
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.3.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.3.0"
    }
  }

  backend "gcs" {
    prefix = "minecraft/state"
    bucket = "minecraft-626"
  }
}

provider "google" {
  project = local.project
  region  = local.region
}

provider "google-beta" {
  project = local.project
  region  = local.region
}

# Create service account to run service with no permissions
resource "google_service_account" "minecraft" {
  account_id   = "minecraft"
  display_name = "minecraft"
}

# Permenant IP address, stays around when VM is off
resource "google_compute_address" "minecraft" {
  name   = "minecraft-ip"
  region = local.region
}

# Create a private network so the minecraft instance cannot access
# any other resources.
resource "google_compute_network" "minecraft" {
  name = "minecraft"
}

# Open the firewall for Minecraft traffic
resource "google_compute_firewall" "minecraft" {
  name    = "minecraft"
  network = google_compute_network.minecraft.name
  # Minecraft client port
  allow {
    protocol = "tcp"
    ports    = ["25565"]
  }
  # ICMP (ping)
  allow {
    protocol = "icmp"
  }
  # SSH (for RCON-CLI access)
  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["minecraft"]
}

resource "google_compute_firewall" "serverless-to-connector" {
  name    = "serverless-to-connector"
  network = google_compute_network.minecraft.name
  # Minecraft client port
  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }
  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }
  allow {
    protocol = "icmp"
  }
  source_ranges = ["35.199.224.0/19"]
  target_tags   = ["vpc-connector"]
}

resource "google_compute_firewall" "minecraft-allow-internal" {
  name    = "minecraft-allow-internal"
  network = google_compute_network.minecraft.name
  # Minecraft client port
  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }
  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }
  allow {
    protocol = "icmp"
  }
  source_ranges = ["10.8.0.0/28"]
  target_tags   = ["minecraft"]
}


# VM to run Minecraft, we use preemptable which will shutdown within 24 hours
resource "google_compute_instance" "minecraft" {
  name         = "minecraft"
  machine_type = "n1-standard-1"
  zone         = local.zone
  tags         = ["minecraft"]

  # Run itzg/minecraft-server docker image on startup
  # The instructions of https://hub.docker.com/r/itzg/minecraft-server/ are applicable
  # For instance, Ssh into the instance and you can run
  #  docker logs mc
  #  docker exec -i mc rcon-cli
  # Once in rcon-cli you can "op <player_id>" to make someone an operator (admin)
  # Use 'sudo journalctl -u google-startup-scripts.service' to retrieve the startup script output
  metadata_startup_script = "docker run -d -p 25565:25565 -p 25575:25575 -e EULA=TRUE -e ENABLE_RCON=true -e VERSION=1.20.2 -v /var/minecraft:/data --name mc -e MEMORY=2G -e RCON_PASSWORD=minecraft --rm=true itzg/minecraft-server:latest;"

  metadata = {
    enable-oslogin = "TRUE"
  }
      
  boot_disk {
    auto_delete = false # Keep disk after shutdown (game data)
    source      = google_compute_disk.minecraft.self_link
  }

  network_interface {
    network = google_compute_network.minecraft.name
    access_config {
      nat_ip = google_compute_address.minecraft.address
    }
  }

  service_account {
    email  = google_service_account.minecraft.email
    scopes = ["userinfo-email"]
  }

  scheduling {
    preemptible       = true # Closes within 24 hours (sometimes sooner)
    automatic_restart = false
  }
}

# Permanent Minecraft disk, stays around when VM is off
resource "google_compute_disk" "minecraft" {
  name  = "minecraft"
  type  = "pd-standard"
  zone  = local.zone
  image = "cos-cloud/cos-stable"
}

resource "google_project_service" "vpcaccess-api" {
  project = local.project
  service = "vpcaccess.googleapis.com"
}

resource "google_cloud_run_v2_service" "server-monitor" {
  name     = "server-monitor"
  project = local.project
  location = local.region

  template {
    containers {
      image = "us-west2-docker.pkg.dev/minecraft-626/server-monitor/mc-server-monitor:latest" # change to other one
      env {
        name = "SERVER_ADDRESS"
        value = ":8080"
      }

      env {
        name = "RCON_ADDRESS"
        value = "${google_compute_instance.minecraft.network_interface.0.network_ip}:25575"
      }

      env {
        name = "RCON_PASSWORD"
        value = "minecraft"
      }
      
      resources {
        limits = {
          cpu    = "1000m"
          memory = "512Mi"
        }
      }
    }

    scaling {
      max_instance_count = 1
    }

    vpc_access {
      # Use the VPC Connector
      connector = google_vpc_access_connector.connector.id
      # all egress from the service should go through the VPC Connector
      egress = "ALL_TRAFFIC"
    }
  }
}

resource "google_cloud_run_service_iam_binding" "server-monitor" {
    location = google_cloud_run_v2_service.server-monitor.location
    service  = google_cloud_run_v2_service.server-monitor.name
    role     = "roles/run.invoker"
    members = [
      "allUsers"
    ]
}

resource "google_vpc_access_connector" "connector" {
  name          = "vpc-con"
  ip_cidr_range = "10.8.0.0/28"
  network       = google_compute_network.minecraft.name
  machine_type = "f1-micro"
  region = local.region
}