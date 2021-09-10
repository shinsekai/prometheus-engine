package operator

import (
	"context"

	monitoringv1alpha1 "github.com/GoogleCloudPlatform/prometheus-engine/pkg/operator/apis/monitoring/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// setupOperatorConfigControllers ensures a rule-evaluator
// deployment as part of managed collection.
func setupOperatorConfigControllers(op *Operator) error {
	// Canonical filter to only capture events for the generated
	// rule evaluator deployment.
	objFilter := namespacedNamePredicate{
		namespace: op.opts.OperatorNamespace,
		name:      nameRuleEvaluator,
	}

	err := ctrl.NewControllerManagedBy(op.manager).
		Named("operator-config").
		For(
			&monitoringv1alpha1.OperatorConfig{},
		).
		Owns(
			&corev1.ConfigMap{},
			builder.WithPredicates(objFilter)).
		Complete(newOperatorConfigReconciler(op.manager.GetClient(), op.opts))

	if err != nil {
		return errors.Wrap(err, "operator-config controller")
	}
	return nil
}

type operatorConfigReconciler struct {
	client client.Client
	opts   Options
}

func newOperatorConfigReconciler(c client.Client, opts Options) *operatorConfigReconciler {
	return &operatorConfigReconciler{
		client: c,
		opts:   opts,
	}
}

// TODO(pintohutch): create the actual deployment.
func (r *operatorConfigReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := logr.FromContext(ctx).WithValues("operatorconfig", req.NamespacedName)

	logger.Info("reconciling OperatorConfig")

	var operatorConfig monitoringv1alpha1.OperatorConfig
	err := r.client.Get(ctx, req.NamespacedName, &operatorConfig)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
