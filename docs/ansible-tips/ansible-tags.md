# ansible tags

You can add tags at different places in Ansible.

- tags to a task

```yaml
- service:
    name: ntpd
    state: started
    enabled: yes
  tags: ntp
```

- tags to a role

```yaml
roles:
  - role: webserver
    tags:
      - webserver
```

- tags to `include`

```yaml
- include: foo.yml
  tags:
    - foo
    - bar
```

- tags to a play

```yaml
- hosts: all
  roles:
    - role: webserver
  tags:
    - bar
```

## `always` tag

- [playbook tags](http://docs.ansible.com/ansible/playbooks_tags.html)
