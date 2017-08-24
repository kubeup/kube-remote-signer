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
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/pflag"
)

// KubeCertificatesController is the main context object for the package.
type KubeCertificatesController struct {
	Kubeconfig          string
	Remote              string
	AuthKey             string
	CertificateDuration metav1.Duration
}

// Create a new instance of a KubeCertificatesController with default parameters.
func NewKubeCertificatesController() *KubeCertificatesController {
	s := &KubeCertificatesController{
		CertificateDuration: metav1.Duration{Duration: 8760 * time.Hour},
		Remote:              os.Getenv("REMOTE_SIGNER_REMOTE"),
		AuthKey:             os.Getenv("REMOTE_SIGNER_AUTH_KEY"),
	}
	return s
}

// AddFlags adds flags for a specific KubeCertificatesController to the
// specified FlagSet.
func (s *KubeCertificatesController) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Kubeconfig, "kubeconfig", s.Kubeconfig, "Path to kubeconfig file with authorization and master location information.")

	fs.StringVar(&s.Remote, "remote", s.Remote, "Remote CFSSL signing server")
	fs.StringVar(&s.AuthKey, "authkey", s.AuthKey, "Secret key used to authenticate with remote CFSSL signing server")
	fs.DurationVar(&s.CertificateDuration.Duration, "certificate-duration", s.CertificateDuration.Duration, "The length of duration signed certificates will be given.")
}
