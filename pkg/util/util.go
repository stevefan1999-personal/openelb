package util

import (
	"context"
	"net"
	"os"

	"github.com/openelb/openelb/pkg/constant"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// IsDeletionCandidate checks if object is candidate to be deleted
func IsDeletionCandidate(obj metav1.Object, finalizer string) bool {
	return obj.GetDeletionTimestamp() != nil && ContainsString(obj.GetFinalizers(), finalizer)
}

// NeedToAddFinalizer checks if need to add finalizer to object
func NeedToAddFinalizer(obj metav1.Object, finalizer string) bool {
	return obj.GetDeletionTimestamp() == nil && !ContainsString(obj.GetFinalizers(), finalizer)
}

// Find node first NodeInternalIP, should check result
func GetNodeIP(node corev1.Node) net.IP {
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeInternalIP {
			return net.ParseIP(address.Address)
		}
	}

	return nil
}

func GetNodeName() string {
	return os.Getenv(constant.EnvNodeName)
}

func DutyOfCNI(metaOld metav1.Object, ObjectNew metav1.Object) bool {
	_, okNew := ObjectNew.GetLabels()[constant.OpenELBCNI]

	if metaOld == nil {
		return okNew
	}

	_, okOld := metaOld.GetLabels()[constant.OpenELBCNI]

	if okOld == okNew && okOld {
		return true
	}

	return false
}

type CheckFn func() bool

func Check(ctx context.Context, c client.Client, obj client.Object, f CheckFn) bool {
	key := client.ObjectKeyFromObject(obj)

	if err := c.Get(ctx, key, obj); err != nil {
		return false
	}

	return f()
}

type CreateFn func() error

func Create(ctx context.Context, c client.Client, obj client.Object, f CreateFn) error {
	err := f()
	if err != nil {
		return err
	}

	if err := c.Create(ctx, obj); err != nil {
		return err
	}

	return nil
}
func EnvNamespace() string {
	ns := os.Getenv(constant.EnvOpenELBNamespace)
	if ns == "" {
		return constant.OpenELBNamespace
	}
	return ns
}
