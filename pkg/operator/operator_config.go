package operator

import (
	"context"

	monitoringv1alpha1 "github.com/GoogleCloudPlatform/prometheus-engine/pkg/operator/apis/monitoring/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	RuleEvaluatorName = "rule-evaluator"
)

// setupOperatorConfigControllers ensures a rule-evaluator
// deployment as part of managed collection.
func setupOperatorConfigControllers(op *Operator) error {
	// Canonical filter to only capture events for the generated
	// rule evaluator deployment.
	objFilter := namespacedNamePredicate{
		namespace: op.opts.OperatorNamespace,
		name:      RuleEvaluatorName,
	}

	err := ctrl.NewControllerManagedBy(op.manager).
		Named("operator-config").
		For(
			&monitoringv1alpha1.OperatorConfig{},
		).
		Owns(
			&corev1.ConfigMap{},
			builder.WithPredicates(objFilter)).
		Owns(
			&appsv1.Deployment{},
			builder.WithPredicates(objFilter)).
		Complete(newOperatorConfigReconciler(op.manager.GetClient(), op.opts))

	if err != nil {
		return errors.Wrap(err, "operator-config controller")
	}
	return nil
}

// operatorConfigReconciler reconciles the OperatorConfig CRD.
type operatorConfigReconciler struct {
	client client.Client
	opts   Options
}

// newOperatorConfigReconciler creates a new operatorConfigReconciler.
func newOperatorConfigReconciler(c client.Client, opts Options) *operatorConfigReconciler {
	return &operatorConfigReconciler{
		client: c,
		opts:   opts,
	}
}

// Reconcile ensures the OperatorConfig resource is reconciled.
func (r *operatorConfigReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := logr.FromContext(ctx).WithValues("operatorconfig", req.NamespacedName)
	logger.Info("reconciling operatorconfig")

	var config = &monitoringv1alpha1.OperatorConfig{}
	if err := r.client.Get(ctx, req.NamespacedName, config); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "get operatorconfig")
	}
	if err := r.ensureRuleEvaluatorConfig(ctx, config); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "ensure rule-evaluator config")
	}
	if err := r.ensureRuleEvaluatorDeployment(ctx); err != nil {
		return reconcile.Result{}, errors.Wrap(err, "ensure rule-evaluator deploy")
	}

	return reconcile.Result{}, nil
}

// ensureRuleEvaluatorConfig reconciles the ConfigMap for rule-evaluator.
func (r *operatorConfigReconciler) ensureRuleEvaluatorConfig(ctx context.Context, config *monitoringv1alpha1.OperatorConfig) error {
	amConfigs, err := makeAlertManagerConfigs(&config.Rules.Alerting)
	if err != nil {
		return errors.Wrap(err, "make alertmanager config")
	}
	cm, err := makeRuleEvaluatorConfigMap(amConfigs, RuleEvaluatorName, r.opts.OperatorNamespace, "prometheus.yml")
	if err != nil {
		return errors.Wrap(err, "make rule-evaluator configmap")
	}

	// Upsert rule-evaluator ConfigMap.
	if err := r.client.Update(ctx, cm); err != nil {
		if err := r.client.Create(ctx, cm); err != nil {
			return errors.Wrap(err, "create rule-evaluator config")
		}
	} else if err != nil {
		return errors.Wrap(err, "update rule-evaluator config")
	}
	return nil
}

// makeRuleEvaluatorConfigMap creates the ConfigMap for rule-evaluator.
func makeRuleEvaluatorConfigMap(amConfigs []yaml.MapSlice, name, namespace, filename string) (*corev1.ConfigMap, error) {
	// Prepare and encode the Prometheus config used in rule-evaluator.
	pmConfig := yaml.MapSlice{}
	pmConfig = append(pmConfig,
		yaml.MapItem{
			Key: "alerting",
			Value: yaml.MapSlice{
				{
					Key:   "alertmanagers",
					Value: amConfigs,
				},
			},
		},
	)
	cfgEncoded, err := yaml.Marshal(pmConfig)
	if err != nil {
		return nil, errors.Wrap(err, "marshal Prometheus config")
	}

	// Create rule-evaluator ConfigMap.
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			filename: string(cfgEncoded),
		},
	}
	return cm, nil
}

// ensureRuleEvaluatorDeployment reconciles the Deployment for rule-evaluator.
func (r *operatorConfigReconciler) ensureRuleEvaluatorDeployment(ctx context.Context) error {
	deploy := r.makeRuleEvaluatorDeployment()

	// Upsert rule-evaluator ConfigMap.
	if err := r.client.Update(ctx, deploy); err != nil {
		if err := r.client.Create(ctx, deploy); err != nil {
			return errors.Wrap(err, "create rule-evaluator deployment")
		}
	} else if err != nil {
		return errors.Wrap(err, "update rule-evaluator deployment")
	}
	return nil
}

// makeRuleEvaluatorDeployment creates the Deployment for rule-evaluator.
func (r *operatorConfigReconciler) makeRuleEvaluatorDeployment() *appsv1.Deployment {
	podLabels := map[string]string{
		LabelAppName: RuleEvaluatorName,
	}
	podAnnotations := map[string]string{
		// TODO(pintohutch): does this need to be unique to the rule-evaluator?
		AnnotationMetricName: collectorComponentName,
	}
	spec := appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: podLabels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      podLabels,
				Annotations: podAnnotations,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "rule-evaluator",
						Image: "gcr.io/gpe-test-1/prometheus-engine/rule-evaluator:bench_20211309_1504",
					},
				},
			},
		},
	}
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: r.opts.OperatorNamespace,
			Name:      RuleEvaluatorName,
		},
		Spec: spec,
	}
	return deploy
}

// makeAlertManagerConfigs creates the alertmanager_config entries as described in
// https://prometheus.io/docs/prometheus/latest/configuration/configuration/#alertmanager_config.
func makeAlertManagerConfigs(spec *monitoringv1alpha1.AlertingSpec) ([]yaml.MapSlice, error) {
	var configs []yaml.MapSlice
	for _, am := range spec.Alertmanagers {
		var cfg yaml.MapSlice
		// Timeout, APIVersion, PathPrefix, and Scheme all resort to defaults if left unspecified.
		if am.Timeout != "" {
			cfg = append(cfg, yaml.MapItem{Key: "timeout", Value: am.Timeout})
		}
		// Default to V2 Alertmanager version.
		if am.APIVersion != "" {
			cfg = append(cfg, yaml.MapItem{Key: "api_version", Value: am.APIVersion})
		}
		// Default to / path prefix.
		if am.PathPrefix != "" {
			cfg = append(cfg, yaml.MapItem{Key: "path_prefix", Value: am.PathPrefix})
		}
		// Default to http scheme.
		if am.Scheme != "" {
			cfg = append(cfg, yaml.MapItem{Key: "scheme", Value: am.Scheme})
		}

		configs = append(configs, cfg)
	}
	return configs, nil
}
