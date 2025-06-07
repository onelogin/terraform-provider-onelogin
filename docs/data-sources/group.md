# onelogin_group Data Source

Use this data source to get information about a specific OneLogin group by ID.

## Example Usage

```hcl
data "onelogin_group" "engineering" {
  id = 123456
}

output "group_name" {
  value = data.onelogin_group.engineering.name
}
```

## Argument Reference

* `id` - (Required) The ID of the group.

## Attribute Reference

* `name` - The name of the group.
* `reference` - A reference identifier for the group.
