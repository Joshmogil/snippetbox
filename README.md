# letsgo
Let's learn go

sudo apt install mysql-server
sudo service mysql start

go run cmd/web/!(*_test).go

#run
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Switch to using the `snippetbox` database.
USE snippetbox;

go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost