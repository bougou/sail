# determine if file exists

```yaml
- name: determine existence of /path/to/app.conf
  stat:
    path: "/path/to/app.conf"
  register: app_conf_stat_result

# Only copy file when not exist
- name: copy kong.conf from kong.conf.default if not exist
  shell: cp /path/to/app.conf.default /path/to/app.conf
  when: not app_conf_stat_result.stat.exists
```
