output "instance_name" {
  value = google_compute_instance.test.name
}

output "public_ip" {
  value = google_compute_instance.test.network_interface[0].access_config[0].nat_ip
}

output "bucket_url" {
  value = google_storage_bucket.test.url
}