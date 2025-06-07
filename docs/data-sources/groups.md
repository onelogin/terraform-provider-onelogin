# onelogin_groups Data Source

Use this data source to get a list of all OneLogin groups.

## Example Usage

```hcl
data "onelogin_groups" "all" {}

output "all_groups" {
  value = data.onelogin_groups.all.groups
}

output "first_group_name" {
  value = data.onelogin_groups.all.groups[0].name
}
```

## Attribute Reference

* `groups` - A list of groups. Each group has the following attributes:
  * `id` - The ID of the group.
  * `name` - The name of the group.
  * `reference` - A reference identifier for the group.
