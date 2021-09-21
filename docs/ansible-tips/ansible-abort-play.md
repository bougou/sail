# Abort ansible playboo

1. [`maximum-failure-percentage`](http://docs.ansible.com/ansible/latest/playbooks_delegation.html#maximum-failure-percentage)

    set `maximum-failure-percentage = 0` means exit the play if any one node fails.

2. [`any_errors_fatal`](http://docs.ansible.com/ansible/latest/playbooks_delegation.html#interrupt-execution-on-any-error)

    With the `any_errors_fatal` option set to `true`, any failure on any host in ***a multi-host play*** will be treated as fatal and Ansible will exit immediately without waiting for the other hosts.
