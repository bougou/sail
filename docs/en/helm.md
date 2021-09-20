# Helm

When deploying the product, all components or some components may need to deployed to K8S.

Sail deploys the pod components by using **Helm Chart**.
Sail supports two mode to use helm chart.

1. Treat the product as a whole helm chart. `_sail_helm_mode: "product"`
2. Treat each component as standalone helm chart. `_sail_helm_mode: "component"`

> These two modes CAN NOT be used simultaneously for a product.

You can specify the helm mode by set `_sail_helm_mode` to proper value in `<target_name>/<zone_name>/vars.yaml`.

## Treat the product as a whole helm chart

You can use the `products/<productName>` as a standarad helm chart directory.

```bash
components/*.yaml
components.yaml
vars.yaml

Chart.yaml      # helm Chart.yaml
values.yaml     # helm values.yaml
templates/      # helm templates
crds/           # helm crds
charts/         # helm charts
```

You can put K8S manifest files of each component under `templates`.

```bash
templates/foobar-api-service.yaml
templates/foobar-api-deployment.yaml
templates/foobar-api-configmap.yaml
templates/foobar-web-service.yaml
templates/foobar-web-deployment.yaml
templates/foobar-web-configmap.yaml
templates/...
```

## Treat each component as standalone helm chart.

The chart directory for each component: `products/<productName>/roles/<roleName>/helm/<componentName>/`.
It's a standarad helm chart directory.

```bash
Chart.yaml
values.yaml
templates/
crds/
charts/
```

Even in `_sail_helm_mode: "component"` mode，you can optionally still define a global values.yaml file `products/<productName>/values.yaml`。

## Helm Runing

Todo.
