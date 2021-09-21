# ansible jinja2

## loop iterate

> http://jinja.pocoo.org/docs/dev/templates/#for

```py
{% for i in a %}
     {{ loop }}
     {{ loop.index }}  # 1 indexed
     {{ loop.index0 }} # 0 indexed
{% endfor %}
```

## iterate dict

```py
{% for key, value in ETC_HOSTS.iteritems() %}
{{ key }} {{ value }}
{% endfor %}
```

## join list into strings

Vars:

```yaml
# a yaml variable - list
docker_insecure_registries:
  - '127.0.0.1'
  - 'gcr.io'
  - '10.5.252.61'
```

Jinja2 template file:
```j2
{
    "insecure-registries": ["{{'\",\"'.join(docker_insecure_registries) }}" ],
}
```

Rendered file:

```j2
{
    "insecure-registries": ["127.0.0.1","gcr.io","10.5.252.61"],
}
```


## transform list to a Comma Separated string

> combine `loop.last` and `whitespace control`

```py
{% for fruit in fruits %}
    {{ fruit }}
    {%- if not loop.last %},{% endif %}
{% endif %}
```
