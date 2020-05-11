provider "google" {
  credentials = file("~/.ssh/gcp-owner-sree-dev-01.json")
  project     = "sree-dev-01"
}

provider "google-beta" {
  credentials = file("~/.ssh/gcp-owner-sree-dev-01.json")
  project     = "sree-dev-01"
}

provider "local" {
}