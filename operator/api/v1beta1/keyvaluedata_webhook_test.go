//nolint
package v1beta1_test

import (
	"github.com/miprokop/crd-kvd/api/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

var _ = Describe("KeyValueData webhook", func() {
	Context("Inside of a new namespace", func() {
		const (
			KeyValueDataName      = "test-keyvaluedata"
			KeyValueDataNamespace = "default"
			timeout               = time.Second * 10
			interval              = time.Millisecond * 250
		)
		BeforeEach(func() {
			k8sClient = fake.NewClientBuilder().WithScheme(scheme.Scheme).Build()
			reconciler.Client = k8sClient
			v1beta1.WebhookClient = k8sClient
		})
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
			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(keyValueData.ValidateCreate()).Should(Not(Succeed()))
			ContainSubstring("key: test-key-string, err : this key already exists")
		})

		It("Create resource with not the same keys value", func() {
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
			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName + "1", Namespace: KeyValueDataNamespace}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(keyValueData1.ValidateCreate()).Should(Succeed())
			Expect(k8sClient.Create(ctx, keyValueData1)).Should(Succeed())

			keyValueDataLookupKey = types.NamespacedName{Name: KeyValueDataName + "2", Namespace: KeyValueDataNamespace}
			_, err = reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(keyValueData2.ValidateCreate()).Should(Succeed())
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

			Expect(keyValueData.ValidateCreate()).Should(Not(Succeed()),
				ContainSubstring(`"empty key or value"`))
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

			Expect(keyValueData.ValidateCreate()).Should(Not(Succeed()),
				ContainSubstring("empty key or value"))
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

			Expect(keyValueData.ValidateCreate()).Should(Not(Succeed()),
				ContainSubstring("empty data field"))
		})

		It("Update resource with the same key value and the same resource", func() {
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
			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			createdKeyValueData := v1beta1.KeyValueData{}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, keyValueDataLookupKey, &createdKeyValueData)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(k8sClient.Update(ctx, &createdKeyValueData)).Should(Succeed())
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
			Expect(keyValueData.ValidateUpdate(nil)).Should(Not(Succeed()),
				ContainSubstring("empty key or value"))
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
			Expect(keyValueData.ValidateUpdate(nil)).Should(Not(Succeed()),
				ContainSubstring("empty key or value"))
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
			Expect(keyValueData.ValidateUpdate(nil)).Should(Not(Succeed()),
				ContainSubstring("empty data field"))
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
			Expect(keyValueData.ValidateUpdate(nil)).Should(Succeed())
		})

		It("Update resource with existing data field", func() {
			By("Empty value")
			keyValueData1 := &v1beta1.KeyValueData{
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
			keyValueData2 := &v1beta1.KeyValueData{
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
						"test-key-string": "test-value-string",
					},
				},
			}
			Expect(keyValueData1.ValidateCreate()).Should(Succeed())
			Expect(k8sClient.Create(ctx, keyValueData1)).Should(Succeed())
			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(keyValueData2.ValidateUpdate(nil)).Should(Not(Succeed()))
			ContainSubstring("key: test-key-string, err : this key already exists")
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

			Expect(keyValueData.ValidateCreate()).Should(Succeed())
			keyValueDataLookupKey := types.NamespacedName{Name: KeyValueDataName, Namespace: KeyValueDataNamespace}
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: keyValueDataLookupKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(keyValueData.ValidateDelete()).Should(Succeed())
		})
		It("Default resource", func() {
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
			keyValueData.Default()
		})
	})
})
