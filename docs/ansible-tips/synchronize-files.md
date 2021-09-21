# Synchronize Files in Ansible

## synchronize

Suppose the following ***three*** kind of nodes.

- `ansbile local host` (deploy machine)
- `ansible target`
- `a third node`

### synchronize from `ansible local host` to `ansible target`

> using default `mode=push`

```yml playbook.yml
- hosts: target
  tasks:
    - synchronize:
        src: "{{ dir_on_ansible_local }}"
        dest: "{{ dir_on_target }}"
```

### synchronize from `ansible target` to `ansible local host`

> using `mode=pull`

```yml playbook.yml
- hosts: target
  tasks:
    - synchronize:
        mode: pull
        src: "{{ dir_on_target }}"
        dest: "{{ dir_on_ansible_local }}"
```

### synchronize from `a third node` to `ansible target`

> using default `mode=push` and `delegate_to` a third node

```yml playbook.yml
- hosts: target
  tasks:
    - synchronize:
        src: "{{ dir_on_third_node }}"
        dest: "{{ dir_on_target }}"
      delegate_to: "{{ a_third_node }}"
```

### synchronize from `ansible target` to `a third node`

> using `mode=pull` and `delegate_to` a third node

```yml playbook.yml
- hosts: target
  tasks:
    - synchronize:
        mode: pull
        src: "{{ dir_on_target }}"
        dest: "{{ dir_on_third_node }}"
      delegate_to: "{{ a_third_node }}"
```

### synchronize from `ansible local host` to `a third node`

> **`synchronize` module can not work it out**.
> Please directly using `rsync` command in ansible `shell` module and `delegate_to` to `ansible_local_host(127.0.0.1)`

```yml playbook.yml
- hosts: target
  tasks:
    - shell: rsync -azv {{ dir_on_ansible_local }} {{ a_third_node }}:{{ dir_on_third_node }}
      delegate_to: 127.0.0.1
```

### synchronize from `a third node` to `ansible local host`

> **`synchronize` module can not work it out**.
> Please directly using `rsync` command in ansible `shell` module and `delegate_to` to `127.0.0.1`

```yml playbook.yml
- hosts: target
  tasks:
  - shell: rsync -azv {{ a_third_node }}:{{ dir_on_thid_node }} {{ dir_on_ansible_local }}
    delegate_to: 127.0.0.1
```
