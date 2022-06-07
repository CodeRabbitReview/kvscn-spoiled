//nolint
package v1beta1_test

import (
	"context"
	"github.com/miprokop/crd-kvd/api/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("KeyValueData controller", func() {
	Context("Inside of a new namespace", func() {
		const (
			KeyValueDataName      = "test-keyvaluedata"
			KeyValueDataNamespace = "default"
		)
		AfterEach(func() {
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}
			k8sClient.Delete(ctx, keyValueData)
		})
		ctx := context.Background()
		It("Create resource with the same key value", func() {
			By("Creating new KeyValueData successfully")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())
			Expect(k8sClient.Create(ctx, keyValueData)).ShouldNot(Succeed())
		})

		It("Create resource with the not the same keys value", func() {
			By("Creating new KeyValueData successfully")
			keyValueData1 := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName + "1",
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string1": "test-value-string1",
					},
				},
			}

			keyValueData2 := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName + "2",
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string2": "test-value-string2",
					},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData1)).Should(Succeed())

			Expect(k8sClient.Create(ctx, keyValueData2)).Should(Succeed())
		})

		It("Create resource with the empty value", func() {
			By("Empty value")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string1": "",
					},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData)).ShouldNot(Succeed())
		})

		It("Create resource with the empty key", func() {
			By("Empty key")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"": "test-value-string1",
					},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData)).ShouldNot(Succeed())
		})

		It("Create resource with the empty data", func() {
			By("Empty data")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData)).ShouldNot(Succeed())
		})

		It("Update resource with the same key value", func() {
			By("Creating new KeyValueData successfully")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())
			Expect(k8sClient.Update(ctx, keyValueData)).ShouldNot(Succeed())
		})

		It("Update resource with the empty value", func() {
			By("Empty value")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}
			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())
			keyValueData.Spec = v1beta1.KeyValueDataSpec{
				Data: map[string]string{
					"test-key-string": "",
				},
			}
			Expect(k8sClient.Update(ctx, keyValueData)).ShouldNot(Succeed())
		})

		It("Update resource with the empty key", func() {
			By("Empty key")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())
			keyValueData.Spec = v1beta1.KeyValueDataSpec{
				Data: map[string]string{
					"": "test-value-string",
				},
			}
			Expect(k8sClient.Update(ctx, keyValueData)).ShouldNot(Succeed())
		})

		It("Update resource with the empty data", func() {
			By("Empty value")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())
			keyValueData.Spec = v1beta1.KeyValueDataSpec{
				Data: map[string]string{},
			}
			Expect(k8sClient.Update(ctx, keyValueData)).ShouldNot(Succeed())
		})

		It("Update resource successfully", func() {
			By("Empty value")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())
			keyValueData.Spec = v1beta1.KeyValueDataSpec{
				Data: map[string]string{
					"bla": "bla",
				},
			}
			Expect(k8sClient.Update(ctx, keyValueData)).Should(Succeed())
		})

		It("Delete resource", func() {
			By("Empty value")
			keyValueData := &v1beta1.KeyValueData{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KeyValueData",
					APIVersion: "key-value.teamdev.com/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      KeyValueDataName,
					Namespace: KeyValueDataNamespace,
				},
				Spec: v1beta1.KeyValueDataSpec{
					Data: map[string]string{
						"test-key-string": "test-value-string",
					},
				},
			}

			Expect(k8sClient.Create(ctx, keyValueData)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, keyValueData)).Should(Succeed())
		})
	})
})
