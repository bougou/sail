# ansible filters


## ansible `ipaddr` jinja filter

Jinaj2 tempalte file:

```j2
addr: {{ IPMI_LAN_ADDR }}
address: {{ IPMI_LAN_ADDR | ipaddr('address') }}
network: {{ IPMI_LAN_ADDR | ipaddr('network') }}
netmask: {{ IPMI_LAN_ADDR | ipaddr('netmask') }}
broadcast: {{ IPMI_LAN_ADDR | ipaddr('broadcast') }}
prefix: {{ IPMI_LAN_ADDR | ipaddr('prefix') }}
subnet: {{ IPMI_LAN_ADDR | ipaddr('subnet') }}
```

eg 1: `IPMI_LAN_ADDR="10.0.34.45/255.255.240.0"`

```json
addr: 10.0.34.45/255.255.240.0
address: 10.0.34.45
network: 10.0.32.0
netmask: 255.255.240.0
broadcast: 10.0.47.255
prefix: 20
subnet: 10.0.32.0/20
```

eg 2: `IPMI_LAN_ADDR="10.0.34.45/22"`

```json
addr: 10.0.34.45/22
address: 10.0.34.45
network: 10.0.32.0
netmask: 255.255.252.0
broadcast: 10.0.35.255
prefix: 22
subnet: 10.0.32.0/22
```

## [Jinja2 Builtin Filters](http://jinja.pocoo.org/docs/2.10/templates/#builtin-filters)

- abs
- attr
- batch
- capitalize
- center
- default
- dictsort
- escape
- filesizeformat
- first
- float
- format
- groupby
- indent
- int
- join
- last
- length
- list
- lower
- map
- max
- min
- pprint
- random
- reject
- rejectattr
- replace
- reverse
- round
- safe
- select
- selectattr
- slice
- sort
- string
- striptags
- sum
- title
- tojson
- trim
- truncate
- unique
- upper
- urlencode
- urlize
- wordcount
- wordwrap
- xmlattr

## [Ansible added filters](https://docs.ansible.com/ansible/latest/user_guide/playbooks_filters.html)

- to_json
- to_yaml
- to_human_json
- to_human_yaml
- from_json
- from_yaml
- mandatory
- default
  - default('somevalue')
  - default('somevalue', true), default('somevalue', false)
  - default(omit)
- min
- max
- flatten
- unqiue
- union
- intersect
- difference
- symmetric_difference
- dict2items
- subelements
- random
- shuffle
- log
- pow
- root
- json_query
- ipaddr
- ipv4
- ipv6
- parse_cli
- parse_xml
- hash
- checksum
- password_hash
- combine
- comment
- urlsplit
- regex_search
- quote
- ternary
- join
- basename
- win_basename
- win_splitdrive
- dirname
- win_dirname
- expanduser
- expandvars
- realpath
- splitext
- b64encode
- to_uuid
- strftime
- type_debug

## Jinja Bultin Tests

- callable
- defined
- divisibleby
- eq
- escaped
- even
- ge
- gt
- in
- iterable
- le
- lower
- lt
- mapping
- ne
- none
- number
- odd
- sameas
- sequence
- string
- undefined
- upper
