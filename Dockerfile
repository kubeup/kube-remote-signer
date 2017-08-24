FROM debian:jessie
RUN apt update && apt install -y ca-certificates
ADD kube-remote-signer /kube-remote-signer
