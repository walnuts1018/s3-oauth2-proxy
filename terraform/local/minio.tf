locals {
  minio_access_key = "access_key"
  minio_secret_key = "secret_key"
}

module "minio" {
  source             = "../modules/minio"
  bucket_name_suffix = ""
  minio_access_key   = local.minio_access_key
  minio_secret_key   = local.minio_secret_key
}
