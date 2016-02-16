default: test_cert.pem

CA/cakey.pem:
	openssl genrsa -out $@ 2048

cacert.pem: CA/cakey.pem
	openssl req -new -x509 -out $@ -key $< -subj "/C=JP/ST=Tokyo/O=KaoriYa/OU=Development/CN=ca.test.kaoriya.net" -days 9999

test_key.pem:
	openssl genrsa -out $@ 2048

test_csr.pem: test_key.pem
	openssl req -new -out $@ -key $< -subj "/C=JP/ST=Tokyo/O=KaoriYa/OU=Development/CN=target.test.kaoriya.net"

test_cert.pem: test_csr.pem cacert.pem CA/cakey.pem
	$(MAKE) ca-reset
	openssl ca -config CA/openssl.cnf -extfile CA/ip.cnf -batch -out $@ -infiles test_csr.pem

ca-reset: ca-clean ca-setup

ca-clean:
	rm -rf CA/certs CA/crl CA/newcerts CA/index.txt* CA/serial*

ca-setup:
	mkdir CA/certs CA/crl CA/newcerts
	touch CA/index.txt
	echo "01" > CA/serial

clean: ca-clean

distclean: clean
	rm -f CA/cakey.pem cacert.pem test_*.pem
