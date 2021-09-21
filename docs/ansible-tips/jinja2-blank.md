# Ansible Jinja Blank

```bash
{%-  # 移除块前的空白
-%}  # 移除块后的空白
```

Suppose:

```bash
# variable
seq="one two three four"
```

## Normal Mode

Jinja2:

```jinja2
{% for item in seq %}
    {{ item }}
{% endfor %}
```

Result:

```txt
    one
    two
    three
    four
```

## Trim blank after the block

Jinja2:

```jinja2
{% for item in seq -%}
    {{ item }}
{% endfor %}
```

Output:

```
one
two
three
four
```

> after the `-%}`, which means space between `-%}` and next `{{`

## Trim Blank before the block

Jinja2:

```jinja2
{% for item in seq %}
    {{ item }}
{%- endfor %}
```

Result:

```
    onetwothreefour
```

> before the `{%-`, which means blank between `}}` and `{%-` (line break).

## Trim Blank before and after the block

Jinja2:

```jinja2
{% for item in seq -%}
    {{ item }}
{%- endfor %}
```

Result:

```
onetwothreefour
```

## See also

- [`trim_blocks` argument of `templates` module](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/template_module.html#parameter-trim_blocks)
