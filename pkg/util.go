package util

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileReturnHelper Ensures that we can return safely with pointers that might be nil
// Essentially just a wrapper for what was a bunch of repetative code
func ReconcileReturnHelper(result *reconcile.Result, err error) (reconcile.Result, error) {
	if result == nil {
		result = &reconcile.Result{}
	}

	return *result, err
}

// ReconcilerStateHelper Returns the desired state based on a bool flag
func ReconcilerStateHelper(enabled bool) reconciler.DesiredState {
	if enabled {
		return reconciler.StatePresent
	}

	return reconciler.StateAbsent
}
