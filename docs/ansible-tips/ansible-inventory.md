# ansible inventory

## one host belongs to multiple groups

```yaml
group-a:
  vars:
    testkey: a
  children:
    group-b:
      vars:
        testkey: b
      children:
        group-b-1:
          vars:
            testkey: "b1"
          hosts:
            10.30.14.5: {}
            10.30.12.5: {}
        group-b-2:
          vars:
            testkey: "b2"
          hosts:
            10.30.14.5: {}
        group-b-3:
          hosts:
            10.30.1.160: {}
```

`group_names` represents all groups to which the host belongs. It is a list.

For `10.30.14.5`, it's `group_names` is `['group-a', 'group-b', 'group-b-1', 'group-b-2']`

If a variable is defined in multiple groups, then the child group will override the parent group.

By default Ansible merges groups at the same parent/child level in ASCII order, and the last group loaded overwrites the previous groups.

see: [ansible: How variables are merged](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#how-variables-are-merged)
