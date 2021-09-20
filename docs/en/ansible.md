# Ansible

Sail uses ONLY one ansible playbook to construct the installation steps for all components.
This is playbook SHOULD be named `sail.yaml` and be placed at `products/<productName>/sail.yaml`.

The `sail.yaml` playbook looks like the following:

```yaml
- name: kafka
  hosts: "{{ _ansiblepattern_kafka | default('kafka') }}"
  roles:
    - role: kafka
      tags:
        - kafka
  tags:
    - play-kafka

- name: etcd
  hosts: "{{ _ansiblepattern_etcd | default('etcd') }}"
  roles:
    - role: etcd
      tags:
        - etcd
  tags:
    - play-etcd
```

For each component, there should be an seperated ansible play (playbook is composed of plays).

The order in the playbook also indicates the installation sequence for the components of the product.

## Generate `sail.yaml`

If the product contains a lot of components, `sail` can help to generate the `sail.yaml` playbook for you.
Please run `sail gen-sail -p <productName>`.

You can indicate the installation order for the components by `order.yaml`.
The `order.yaml` is a list of components. You don't need to specified ALL components into `order.yaml`.
Those unspecified components are automatically appended to the last by alphabetical order.

```yaml
# order.yaml content
- mysql
- kafka
- myapp-web
- myapp-api
...
```

1. `sail` gets the components list of the product by parsing `components.yaml` and `components/*.yaml`.
2. `sail` gets the order by parsing `order.yaml`.
3. `sail` generates the `sail.yaml` playbook.
