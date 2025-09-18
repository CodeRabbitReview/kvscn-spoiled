/*
Copyright 2022.

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
//nolint
package v1beta1_test

import (
	"context"
	"github.com/miprokop/crd-kvd/api/v1beta1"
	"github.com/miprokop/crd-kvd/controllers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"net/http"
	"net/http/httptest"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var k8sClient client.Client
var reconciler controllers.KeyValueDataReconciler
var ctx context.Context

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Webhook Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	_ = v1beta1.AddToScheme(scheme.Scheme)
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	testNumber := 0
	var testServer = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if testNumber == 0 {
			rw.WriteHeader(http.StatusCreated)
			testNumber++
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}))
	reconciler.ServerURL = testServer.URL
	reconciler.Scheme = scheme.Scheme
	reconciler.HTTPClient = &http.Client{}
	reconciler.FinalizerName = "kubernetes"

}, 60)

var _ = AfterSuite(func() {

})
