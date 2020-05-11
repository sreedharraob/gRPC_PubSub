output "pubsub_project_service_id" {
  value = google_pubsub_subscription.default.path
}

output "pubsub_subscription_path" {
  value = google_pubsub_subscription.default.path
}

output "pubsub_subscription_id" {
  value = google_pubsub_subscription.default.id
}

output "pubsub_topic_id" {
  value = google_pubsub_topic.default.id
}

output "pubsub_sa_email" {
  value = google_service_account.pubsub_sa.email
}

output "pubsub_sa_name" {
  value = google_service_account.pubsub_sa.name
}

output "pubsubs_sa_topic_iam_bindings" {
  value = google_pubsub_topic_iam_binding.default.*.role
}

output "pubsub_sa_subscription_iam_bindings" {
  value = google_pubsub_subscription_iam_binding.default.*.role
}

output "pubsub_sa_key_filename" {
  value = local_file.pubsub_sa_key_json.filename
}