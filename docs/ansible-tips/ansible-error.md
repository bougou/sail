# ansible errors

## ansible yum ok but failed

```bash
# on the target, run:
$ yum-complete-transaction --cleanup-only
```

## ansible can not create tmp file

```ini
TASK [Gathering Facts] ****************************************************************************************************************
Monday 19 April 2021  16:41:14 +0800 (0:00:00.410)       0:00:00.410 **********
fatal: [10.10.0.7]: UNREACHABLE! => changed=false
  msg: 'Authentication or permission failure. In some cases, you may have been able to authenticate and did not have permissions on the target directory. Consider changing the remote tmp path in ansible.cfg to a path rooted in "/tmp". Failed command was: ( umask 77 && mkdir -p "` echo /tmp/.ansible-$USER/ansible-tmp-1618821674.92-213924690448765 `" && echo ansible-tmp-1618821674.92-213924690448765="` echo /tmp/.ansible-$USER/ansible-tmp-1618821674.92-213924690448765 `" ), exited with result 1'
  unreachable: true
```

Ansible can't create temporary files. The reason may be: 1 the system disk is full or 2 the system disk is in a readonly state.

