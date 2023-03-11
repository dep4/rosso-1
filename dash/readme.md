# DASH

ISO/IEC 23009-1:2019

- <https://wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP>
- https://standards.iso.org/ittf/PubliclyAvailableStandards

## Dotted keys

TOML supports dotted keys:

~~~toml
physical.color = "orange"
physical.shape = "round"
~~~

JSON result:

~~~json
{
  "physical": {
    "color": "orange",
    "shape": "round"
  }
}
~~~

https://toml.io/en/v1.0.0#keys

I tried this YAML:

~~~yaml
physical.color: orange
physical.shape: round
~~~

but the resultant JSON is not desirable:

~~~json
{
  "physical.color": "orange",
  "physical.shape": "round"
}
~~~

instead requiring explicit indentation:

~~~yaml
physical:
   color: orange
   shape: round
~~~

https://github.com/yaml/yaml-spec/issues/302
