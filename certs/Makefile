gen_all: gen_altsubj gen_ca gen_server_csr gen_server_crt gen_client_csr gen_client_crt
gen_altsubj:
	echo "subjectAltName = DNS:localhost" > altsubj.ext

gen_ca:
	openssl req -new -x509 \
		-newkey rsa:4096 \
		-keyout ca.key \
		-out ca.crt \
		-sha256 -days 365 -nodes \
		-subj "/CN=localhost" \
		-addext "subjectAltName = DNS:localhost"


gen_server_csr:
	openssl req -new -sha256 \
		-newkey rsa:2048 \
		-keyout server.key \
		-out server.csr \
		-nodes \
		-subj "/CN=localhost" \
		-addext "subjectAltName = DNS:localhost"


gen_server_crt:
	openssl x509 -req -days 365 \
		-sha256 -in server.csr \
		-CA ca.crt -CAkey ca.key \
		-set_serial 1 -out server.crt \
		-extfile altsubj.ext

gen_client_csr:
	openssl req -new -sha256 \
		-newkey rsa:2048 \
		-keyout client.key \
		-out client.csr \
		-nodes \
		-subj "/CN=localhost" \
		-addext "subjectAltName = DNS:localhost"


gen_client_crt:
	openssl x509 -req -days 365 \
		-sha256 -in client.csr \
		-CA ca.crt -CAkey ca.key \
		-set_serial 2 -out client.crt \
		-extfile altsubj.ext
