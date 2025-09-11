resource "aws_s3_bucket" "test" {
  bucket = format("test%s", var.bucket_name_suffix)
}
