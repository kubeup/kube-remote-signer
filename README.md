# kube-remote-signer

kube-remote-signer is a Kubernetes controller doing CSR signing using remote CFSSL CA server.
It is built together with [kube-ca][kube-ca] as a proof of concept about using external CA for CSR signing
to enhance security for Kubernetes clusters.

Currently we have to put the CA private key and certificate on master nodes and pass them to the builtin
certificate controller running in [kube-controller-manager][kube-controller-manager] to support the token
based node bootstrapping process. It is a burden to manage the CA private key properly and there are risks
about key leaking which would leads to critical security incidents.

By moving the signer out of the Kubernetes cluster, we could reduce security risk and simplify the
configuration process for master servers.

## Features

- CSR controller using remote CFSSL CA server
- HMAC authentication to avoid unauthorized access

## Installation

Fist we have to disable the internal certificate controller.

```
sed -i '/controllers/ s/$/,-csrsigning/' /etc/kubernetes/manifests/kube-controller-manager.yaml
```

Then we create a `Secret` containing the remote address and HMAC key.

```
kubectl create secret generic remote-signer-config -n kube-system --from-literal=REMOTE_SIGNER_REMOTE=PATH_TO_CFSSL_SERVER --from-literal=REMOTE_SIGNER_AUTH_KEY=HMAC_KEY
```

Last we run `kube-remote-signer` in the cluster.

```
kubectl create -f https://raw.githubusercontent.com/kubeup/kube-remote-signer/master/kube-remote-signer.yaml
```

## License

Apache Version 2.0

[kube-ca]: https://github.com/kubeup/kube-ca
[kube-controller-manager]: https://kubernetes.io/docs/admin/kube-controller-manager/
