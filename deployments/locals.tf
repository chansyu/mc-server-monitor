# You need to fill these locals out with the project, region, zone and mc-version
# Then to boot it up, run:-
#   gcloud auth application-default login
#   terraform init
#   terraform apply
locals {
  # The Google Cloud Project ID that will host and pay for your Minecraft server
  project = "INSERT PROJECT NAME"
  region  = "us-west2"
  zone    = "us-west2-a"
  bucket  = "INSERT BUCKET NAME"
  mc_version = "1.20.2"
}