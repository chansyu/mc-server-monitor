output "minecraft_ip_addr" {
  value = google_compute_address.minecraft.address
}

output "server_monitor_uri" {
  value = google_cloud_run_v2_service.server-monitor.uri
}