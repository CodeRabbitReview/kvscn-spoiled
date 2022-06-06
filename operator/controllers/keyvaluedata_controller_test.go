package controllers

import (
	"context"
	"github.com/miprokop/crd-kvd/api/v1beta1"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"

	. "github.com/onsi/ginkgo"
)

var one int32 = 1
var zero int32 = 0

var _ = Describe("KeyValueData controller", func() {
	const (
		KeyValueDataName1     = "test-keyvaluedata1"
		KeyValueDataName2     = "test-keyvaluedata2"
		KeyValueDataNamespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)
	Context("Success creation in the server", func() {
		It("Should change KeyValueData Status to one successes action", func() {
			By("Creating new KeyValueData successfully")
			ctx := context.Background()
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName1,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())

			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName1, Namespace: KeyValueDataNamespace}
			createdKeyValueData := v1beta1.KeyValueData{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, &createdKeyValueData)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})
	Context("Unsuccessful creation in the server", func() {
		It("Should change KeyValueData Status to one Unsuccessful action", func() {
			By("Creating new KeyValueData unsuccessfully")
			ctx := context.Background()
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName2,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string2": "test-value-string2",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())

			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName2, Namespace: KeyValueDataNamespace}
			createdKeyValueData := &v1beta1.KeyValueData{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, createdKeyValueData)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})
})
