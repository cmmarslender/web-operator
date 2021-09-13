# Web Operator

A kubernetes operator for deploying a simple/typical webapp consisting of a deployment, service, and ingress.

## Dev

Generated using kubebuilder.

Project initialized with `kubebuilder init --domain k8s.cmm.io --repo github.com/cmmarslender/web-operator --component-config`

Additional APIs created with commands like `kubebuilder create api --group webapp --version v1 --kind SimpleApp`

Resource reconciliation utilizes `github.com/banzaicloud/operator-tools/pkg/reconciler` as a reconciler.

### Config

Configuration is managed with a type defined with kubebuilder. To add a new config type:

`kubebuilder create api --group config --version v1 --kind Config --resource --controller=false --make=false`
