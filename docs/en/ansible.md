# Ansible

Sail uses ONLY one ansible playbook to construct the installation steps for all components.
This playbook SHOULD be named `sail.yaml` and be placed at `products/<productName>/sail.yaml`.

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

## Write Ansible Role

- [Write Ansible Role](ansible-roles.md)

## How sail uses Ansible

Based on [Ansible Official Best Practices](https://docs.ansible.com/ansible/latest/user_guide/playbooks_best_practices.html), Sail makes some changes.

1. Not using `group_vars` and `host_vars`.

    In Sail, you MUST define the variables for the Product in `products/<productName>/{vars.yaml,components.yaml,components/*.yaml}`.
    And for specific environment (target/zone), there will be a correspoinding `targets/<target>/<zone>/vars.yaml` file which holding same variables defined for the product installed in the environment.

    To manage multiple environments, Sail uses a targets directory to store all configurations of the environments. Each target/zone holds its variables and inventory separately.

2. Using **Product** concept to weaken the concept of **Ansible Playbook**. Each product may have only one playbook named `sail.yaml`.
