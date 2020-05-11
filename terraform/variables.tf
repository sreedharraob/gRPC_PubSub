variable "project" {
  type        = string
  description = "name of the project"
}

variable "region" {
  type        = string
  description = "name of the region"
}

variable "pubsub_topic_name" {
  type        = string
  description = "name of the pubsub topic"
}

variable "pubsub_subscription_name" {
  type        = string
  description = "name of the pubsub subscription"
}

variable "pubsub_sa_name" {
  type        = string
  description = "service account name to publish, read and view pubsub topics and subscriptions"
}

variable "pubsub_sa_topic_roles" {
  type        = list(string)
  description = "list of iam role bindings required for pubsub account"
}

variable "pubsub_sa_subscriber_roles" {
  type        = list(string)
  description = "list of iam role bindings required for pubsub account"
}