# vars

```yaml
cert:
  valid_days: 800
  valid_hours: "{{ cert['valid_days'] * 24 }}"    # Error, leads to recursive loop
```

```yaml
data_dir: /data
log_dir: "{{ data_dir }}/log" # Right
```
