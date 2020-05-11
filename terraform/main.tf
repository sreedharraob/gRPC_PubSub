resource "google_project_service" "pubsub_service" {
  disable_on_destroy = false
  project            = var.project
  provider           = google-beta
  service            = "pubsub.googleapis.com"
}

resource "google_pubsub_topic" "default" {
  provider = google-beta
  name     = var.pubsub_topic_name
  project  = var.project

  message_storage_policy {
    allowed_persistence_regions = [var.region]
  }

  depends_on = [google_project_service.pubsub_service]
}

resource "google_pubsub_subscription" "default" {
  provider                   = google-beta
  project                    = var.project
  name                       = var.pubsub_subscription_name
  topic                      = google_pubsub_topic.default.name
  ack_deadline_seconds       = 20
  message_retention_duration = "86400s"

  expiration_policy {
    ttl = "172800s"
  }

  depends_on = [google_pubsub_topic.default]
}

resource "google_service_account" "pubsub_sa" {
  project      = var.project
  provider     = google-beta
  account_id   = var.pubsub_sa_name
  display_name = var.pubsub_sa_name

  depends_on = [google_pubsub_subscription.default]
}

resource "google_pubsub_topic_iam_binding" "default" {
  count    = length(var.pubsub_sa_topic_roles)
  provider = google-beta
  topic    = google_pubsub_topic.default.name
  role     = var.pubsub_sa_topic_roles[count.index]
  members  = ["serviceAccount:${google_service_account.pubsub_sa.email}"]

  depends_on = [google_service_account.pubsub_sa]
}

resource "google_pubsub_subscription_iam_binding" "default" {
  count        = length(var.pubsub_sa_subscriber_roles)
  provider     = google-beta
  subscription = google_pubsub_subscription.default.name
  role         = var.pubsub_sa_subscriber_roles[count.index]
  members      = ["serviceAccount:${google_service_account.pubsub_sa.email}"]

  depends_on = [google_service_account.pubsub_sa]
}

resource "google_service_account_key" "pubsub_sa_key" {
  provider           = google-beta
  service_account_id = google_service_account.pubsub_sa.name

  depends_on = [google_service_account.pubsub_sa]
}

resource "local_file" "pubsub_sa_key_json" {
  content  = base64decode(google_service_account_key.pubsub_sa_key.private_key)
  #filename = "${path.module}/../../../../../../../.ssh/gcp_pubsub_sa.json"
  filename = "${path.module}/gcp_pubsub_sa.json"
}