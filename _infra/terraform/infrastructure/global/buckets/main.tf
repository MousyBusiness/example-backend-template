resource "google_storage_bucket" "example_bucket" {
  name          = "${var.project}-example"
  location      = "EU"
  force_destroy = true
}

// !!! PUBLIC !!!
resource "google_storage_bucket" "downloads_bucket" {
  name          = "${var.project}-downloads"
  location      = "EU"
  force_destroy = false
}

data "google_iam_policy" "viewer" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "allUsers",
    ]
  }
}

resource "google_storage_bucket_iam_policy" "policy" {
  bucket = google_storage_bucket.downloads_bucket.name
  policy_data = data.google_iam_policy.viewer.policy_data
}
