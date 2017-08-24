/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"fmt"
	"time"

	"github.com/cloudflare/cfssl/auth"
	"k8s.io/client-go/tools/record"
	capi "k8s.io/kubernetes/pkg/apis/certificates/v1beta1"
	"k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
	"k8s.io/kubernetes/pkg/controller/certificates"

	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/remote"
)

type remoteSigner struct {
	remote              string
	authKey             string
	client              clientset.Interface
	certificateDuration time.Duration
	recorder            record.EventRecorder
}

func NewRemoteSigner(remote, authKey string, client clientset.Interface, certificateDuration time.Duration, recorder record.EventRecorder) (*remoteSigner, error) {
	if remote == "" {
		return nil, fmt.Errorf("remote should not be emtpy")
	}
	if authKey == "" {
		return nil, fmt.Errorf("authKey should not be emtpy")
	}
	return &remoteSigner{
		remote:              remote,
		authKey:             authKey,
		client:              client,
		certificateDuration: certificateDuration,
		recorder:            recorder,
	}, nil
}

func (s *remoteSigner) handle(csr *capi.CertificateSigningRequest) error {
	if !certificates.IsCertificateRequestApproved(csr) {
		return nil
	}
	csr, err := s.sign(csr)
	if err != nil {
		return fmt.Errorf("error auto signing csr: %v", err)
	}
	_, err = s.client.Certificates().CertificateSigningRequests().UpdateStatus(csr)
	if err != nil {
		return fmt.Errorf("error updating signature for csr: %v", err)
	}
	return nil
}

func (s *remoteSigner) sign(csr *capi.CertificateSigningRequest) (*capi.CertificateSigningRequest, error) {
	var usages []string
	for _, usage := range csr.Spec.Usages {
		usages = append(usages, string(usage))
	}
	provider, err := auth.New(s.authKey, nil)
	if err != nil {
		return nil, err
	}
	policy := &config.Signing{
		Default: &config.SigningProfile{
			Usage:          usages,
			Expiry:         s.certificateDuration,
			ExpiryString:   s.certificateDuration.String(),
			RemoteProvider: provider,
		},
	}
	policy.OverrideRemotes(s.remote)
	cfs, err := remote.NewSigner(policy)
	if err != nil {
		return nil, err
	}

	csr.Status.Certificate, err = cfs.Sign(signer.SignRequest{
		Request: string(csr.Spec.Request),
	})
	if err != nil {
		return nil, err
	}

	return csr, nil
}
