terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

provider "google" {
  region = var.region
}

# GCP resources

resource "google_compute_instance" "test" {
  project = var.project_id
  name         = var.instance_name
  machine_type = var.machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1604-lts"
    }
  }

  network_interface {
    network = "default"
    subnetwork = "default"
    access_config {}
  }
}

resource "google_storage_bucket" "test" {
  project = var.project_id
  name = var.bucket_name
  force_destroy = true
}