#!/usr/bin/env bash

SSL_KEY_FILE=server.key
SSL_CERT_FILE=server.crt

# Use OpenSSL to generate ssl keys using RSA algorithm
openssl genrsa -out $SSL_KEY_FILE 2048;

# Use OpenSSL to generate keys with ECDSA algorithm
# openssl ecparam -genkey -name secp384r1 -out $SSL_KEY_FILE;

# Create ssl certificate based on the new private key
openssl req -new -x509 -sha256 -key $SSL_KEY_FILE -out $SSL_CERT_FILE -days 3650;
