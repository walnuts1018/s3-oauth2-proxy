resource "aws_s3_object" "index" {
  bucket       = aws_s3_bucket.test.id
  key          = "index.html"
  source       = "./src/index.html"
  etag         = filemd5("./src/index.html")
  content_type = "text/html"
}
