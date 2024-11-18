# Micro Hiera

Merge multiple yaml files together

## Features

* Deep merge of lists and maps
* Duplicate value detection


## examples

Example file has `.duplicate_value_in_over_as_base` defined with the same value in 2 files

```sh
$ ./micro-hiera merge ./examples/input/*
DUPLICATE_OVERRIDE_VALUE: ./examples/input/a_over_1.yml duplicate value at path:.duplicate_value_in_over_as_base

violation                      count
--------------------------------------------------
DUPLICATE_OVERRIDE_VALUE       1
```

Correcting the issue by deleting `.duplicate_value_in_over_as_base` from `./examples/input/a_over_1.yml`

```sh
$ ./micro-hiera merge ./examples/input/*
duplicate_value_in_over_as_base: true
file: a_over_1.yml
list:
- a_over_1.yml
- a_base.yml
map:
  file: a_over_1.yml
  file_list:
  - a_over_1.yml
  - a_base.yml
new_key_for_a_over_1: true
```