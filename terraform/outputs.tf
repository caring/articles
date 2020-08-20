# Add any output values to this file
output "service_url" {
  value = module.service_dns_record.fqdn
}
