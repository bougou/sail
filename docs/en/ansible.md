# Ansible

## Ansible playbook

Sail uses ONLY one ansible playbook to construct the installation steps for all components.

But you do not need to write this playbook by youself, `sail` will automatically generate the playbook when it runs. The auto generated playbook file is named `.sail.yaml`.

For a specific environment (target/zone), the name of the product applied to the environment is recorded
in environment `vars.yaml` file, so `sail` will automatically generated a playbook file of the product and put `.sail.yaml` under `products/<productName>` and refers to it when running `ansible-playbook`.

The generated `.sail.yaml` playbook looks like the following. For each component of the product, there is a correspoding play block.


```yaml
- name: kafka
  hosts: "{{ _ansiblepattern_kafka | default('kafka') }}"
  any_errors_fatal: false
  gather_facts: true
  become: false
  roles:
    - role: kafka
      tags:
        - kafka
  tags:
    - play-kafka
```

The `.sail.yaml` playbook is composed of plays of the components of the product.
The order in the playbook also indicates the installation sequence for the components of the product.

But you can change the installation order of the components by setting `products/<productName>/order.yaml`.
The `order.yaml` is a list of components. You don't need to specify all components into `order.yaml`.
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
3. `sail` generates the `.sail.yaml` playbook.

> You can use `sail gen-sail -p <productName>` to manually generate it and get a look at `.sail.yaml`.

## Write Ansible Role

- [Write Ansible Role](ansible-roles.md)

## How sail uses Ansible

Based on [Ansible Official Best Practices](https://docs.ansible.com/ansible/latest/user_guide/playbooks_best_practices.html), Sail makes some changes.

1. Sail does not using `group_vars` and `host_vars`.

    In Sail, you MUST define the variables for the Product in `products/<productName>/{vars.yaml,components.yaml,components/*.yaml}`.

    And for specific environment (target/zone), there will be a correspoinding `targets/<target>/<zone>/vars.yaml` file which holding same variables structure defined for the product installed in the environment.

    To manage multiple environments, Sail uses a targets directory to store all configurations of the environments. Each environment (target/zone) holds its variables and inventory separately.

2. Sail uses **Product** concept to weaken the concept of **Ansible Playbook**. `sail` uses a default and auto generated playbook for the product.

    Sail promotes to write ansible tasks in the roles for product components. Ansible playbook is just a glue of the component hosts group and the corresponding roles. And `sail` will automatically generates the playbook. You can focus on the tasks on components.
