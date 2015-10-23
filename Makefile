install:
	- mkdir -p crypto
	- openssl genrsa -out crypto/app.rsa 1024
	- openssl rsa -in crypto/app.rsa -pubout > crypto/app.rsa.pub
