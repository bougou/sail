# cmdb

## `hosts.yaml`

The servers for components deployed by ansible should be kept under `targets/<target>/<zone>/hosts.yaml`.
This `hosts.yaml` is essentially an yaml formated hosts inventory recognized by Ansible.

### How to change SSH for a host

```yaml
all:
  vars:
    ansible_user: root
    ansible_port: 22
somegroupname:
  vars:
    ansible_user: root
    ansible_port: 22
  hosts:
    10.0.0.1:
      ansible_user: root
      ansible_port: 2222
```

If the SSH addr for the host is not equal to the keys under the `hosts` field,
you can set it to other address by setting `ansible_host` for the host.

```yaml
somegroupname:
  hosts:
    10.0.0.1:
      ansible_user: root
      ansible_port: 2222
      ansible_host: 127.0.0.1
```

## `platforms.yaml`

The k8s platform information for components deployed by helm should be kept under `targets/<target>/<zone>/platforms.yaml`.

Generally, all components of a product is deployed to one namespace under a K8S cluster.
In this case, you just needs to set one `all` record in `platforms.yaml`.

```yaml
all:
  k8s:
    kubeConfig: ~/.kube/config
    kubeContext: ""
    namespace: default
```

But, if the components are needed to deployed to different namespaces or different K8S cluster,
you can set seperated k8s information for each component.

```yaml
all:
  k8s:
    kubeConfig: ~/.kube/config
    kubeContext: ""
    namespace: default

component-name1:
  k8s:
    kubeConfig: ~/.kube/config
    kubeContext: ""
    namespace: another

component-name2:
  k8s:
    kubeConfig: ~/.kube/config
    kubeContext: ""
    namespace: another
```

The `all` record is used for all other components which are not explicitly defined.
