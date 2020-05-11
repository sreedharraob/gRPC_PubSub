project                  = "sree-dev-01"
region                   = "europe-west2"
pubsub_topic_name        = "grpc-topic"
pubsub_subscription_name = "grpc-sub-01"
pubsub_sa_name           = "grpc-pubsub-sa"
pubsub_sa_topic_roles = [
  "roles/pubsub.viewer",
  "roles/pubsub.publisher"
]
pubsub_sa_subscriber_roles = [
  "roles/pubsub.viewer",
  "roles/pubsub.subscriber",
]