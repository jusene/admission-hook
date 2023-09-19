#!/bin/bash

usage() {
    cat <<EOF
Generate certificate suitable for use with an sidecar-injector webhook service.

This script uses k8s' CertificateSigningRequest API to a generate a
certificate signed by k8s CA suitable for use with sidecar-injector webhook
services. This requires permissions to create and approve CSR. See
https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster for
detailed explantion and additional instructions.

The server key/cert k8s CA cert are stored in a k8s secret.

usage: ${0} [OPTIONS]

The following flags are required.

       --service          Service name of webhook.
       --namespace        Namespace where webhook service and secret reside.
       --secret           Secret name for CA certificate and server certificate/key pair.
EOF
    exit 1
}

while [[ $# -gt 0 ]]; do
    case ${1} in
        --service)
            service="$2"
            shift
            ;;
        --secret)
            secret="$2"
            shift
            ;;
        --namespace)
            namespace="$2"
            shift
            ;;
        *)
            usage
            ;;
    esac
    shift
done

[ -z ${service} ] && service=admission-webhook-svc
[ -z ${secret} ] && secret=admission-webhook-certs
[ -z ${namespace} ] && namespace=default

if [ ! -x "$(command -v openssl)" ]; then
    echo "openssl not found"
    exit 1
fi

mkdir certs
echo "creating certs in certs directory "

cat > certs/csr.conf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
prompt = no
[req_distinguished_name]
CN = ${service}.${namespace}.svc
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${service}.${namespace}
DNS.2 = ${service}.${namespace}.svc
DNS.3 = ${service}
EOF

# 生成ca证书
openssl genrsa -out certs/ca.key 2048
openssl req -x509 -new -nodes -key certs/ca.key -days 365 -out certs/ca.crt -subj "/CN=admission_ca"

# 生成tls密钥和证书
openssl genrsa -out certs/server.key 2048
openssl req -new -key certs/server.key -out certs/server.csr -config certs/csr.conf
openssl x509 -req -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/server.crt -days 365 -extensions v3_req -extfile certs/csr.conf

# 为webhook创建kubernetes tls secret
kubectl create secret tls ${secret} --cert=certs/server.crt --key=certs/server.key -n ${namespace}

# 设置CA_BUNDLE
export CA_BUNDLE=`cat certs/ca.crt | base64 | tr -d '\n'`
echo ${CA_BUNDLE}