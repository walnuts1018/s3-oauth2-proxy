.PHONY: cert
cert:
	mkcert -cert-file ./certs/s3-oauth2-proxy.local.walnuts.dev.pem -key-file ./certs/s3-oauth2-proxy.local.walnuts.dev-key.pem s3-oauth2-proxy.local.walnuts.dev
	mkcert -cert-file ./certs/authelia.local.walnuts.dev.crt -key-file ./certs/authelia.local.walnuts.dev.key authelia.local.walnuts.dev
	mkcert -install
