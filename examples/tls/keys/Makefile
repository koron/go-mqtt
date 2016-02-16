default: test_cert.pem

test_key.pem:
	openssl genrsa -out $@

test_csr.pem: test_key.pem
	openssl req -new -out $@ -key $< -subj "/C=JP/ST=Tokyo/O=KaoriYa/OU=Development/CN=target.test.kaoriya.net"

test_cert.pem: test_csr.pem
	$(MAKE) ca-reset
	openssl ca -config CA/openssl.cnf -extfile CA/ip.cnf -batch -out $@ -infiles $<

ca-reset: ca-clean ca-setup

ca-clean:
	rm -rf CA/certs CA/crl CA/newcerts CA/index.txt* CA/serial*

ca-setup:
	mkdir CA/certs CA/crl CA/newcerts
	touch CA/index.txt
	echo "01" > CA/serial

clean: ca-clean
