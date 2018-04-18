# Notes

## SSL

Here are some notes for how to use SSL.

### Certificate Signing Authorities

You can actually create your own Certificate Signing Authority using OpenSSL.
The following link describes the process of making your own Certificate Signing
Request (CSR) and Certificate Signing Authority (CA):

[The Data Center Overlords - Creating your own SSL Certificate Authority (and Dumping Self Signed Certs)](https://datacenteroverlords.com/2012/03/01/creating-your-own-ssl-certificate-authority/)


#### Creating a CA (General)

It's relatively straight-forward to create a CA. The following steps detail
the process.

The generalized steps are as follows:
1. Create a Private Key
2. Self-sign
3. Install Root CA on your various workstations

Subsequently, every device you manage via HTTPS needs to have its own certificate
created with the following steps:
1. Create CSR for device
2. Sign CSR with Root CA Key


#### Creating a CA (Breakdown)

The following steps are a more technical breakdown of the CA creation process.

1. Create Root Certificate (Done Once)
   This creates a root SSL certificate that you'll install on every target device,
   and a private key you'll use to sign the certificates on those devices.

 * Create the Root Key
  * `openssl genrsa -out rootCA.key 2048`
  * This command creates the 2048-bit Root CA Key: `rootCA.key`
  * Add `-des3` to password-protect the certificate
  * This key is the basis for ALL trust for your certificates. If you lose this,
    other people can forge certificates with this information, and cause your
    browser to accept insecure connections.

 * Self-Sign the Certificate
  * `openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 1024 -out rootCA.pem`
   * This will ask for a bunch of extra information for creating your cert.
   * NOTE: You can use a configuration file for skipping the interactive setup.
   * This will make the Root SSL Certificate: `rootCA.pem`


2. Install Root Certificate into Workstations
 * Go into your browser, find the "Trusted Root Certificate Authorities" section
   of the settings, and choose to install the new certificate here.


3. Create a Certificate (Done Once per Device)

 * Create a Private Key
  * `openssl genres -out device.key 2048`

 * Generate Certificate Signing Request
  * `openssl req -new -key device.key -out device.csr`
  * This will ask you for more information like the `rootCA.pem` file above.
  * NOTE: You MUST set `Common Name: [...]` to the REAL IP Address or
    name of the current system that you're creating the CSR for.
  * If CommonName is not set properly, you will receive a "cannot verify authenticity"
    error message.

 * Sign CSR
  * `openssl x509 -req -in device.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out device.crt -days 500 -sha256`
  * This signs the Certificate Signing Request, creating a new SSL Certificate
    for your device.
