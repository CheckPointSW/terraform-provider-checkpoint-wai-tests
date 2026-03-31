# Resource: cpwai_browse_secure_browsing_rule

Manages a Secure Browsing rule in the Browse Protect policy. Controls web browsing security features such as phishing protection, password reuse detection, safe search, and search reputation.

## Example Usage

```hcl
resource "cpwai_browse_secure_browsing_rule" "protect_internal" {
  name        = "Protect Internal Domains"
  description = "Enable secure browsing protection for internal corporate domains"
  order       = 0
  active      = true

  policy = jsonencode({
    password_reuse    = "prevent"
    zero_phishing     = "detect"
    safe_search       = "on"
    search_reputation = "on"
    protected_domains = [
      {
        domain     = "portal.myorg.com"
        match_mode = "exact_host"
      },
      {
        domain     = "apps.myorg.com"
        match_mode = "any_subdomain"
      }
    ]
  })

  source = [
    {
      assignment_type = "ASSIGNMENT_TYPE_ENTIRE_ORG"
    }
  ]
}
```

### Example with Protected Domains Object

Create a reusable domains object and attach it to a secure browsing rule as protected domains:

```hcl
resource "cpwai_browse_domains_object" "org_internal_domains" {
  name        = "Organization Internal Domains"
  description = "Internal corporate domains that require browsing protection"

  object_data = jsonencode({
    domains = [
      { domain = "portal.myorg.com",   match_mode = "exact_host" },
      { domain = "apps.myorg.com",     match_mode = "any_subdomain" },
      { domain = "intranet.myorg.com", match_mode = "any_subdomain" }
    ]
  })
}

resource "cpwai_browse_secure_browsing_rule" "protect_internal_with_objects" {
  name        = "Protect Internal Domains"
  description = "Enable secure browsing for internal corporate domains"
  order       = 0
  active      = true

  policy = jsonencode({
    password_reuse    = "prevent"
    zero_phishing     = "prevent"
    safe_search       = "on"
    search_reputation = "on"
  })

  source = [
    { assignment_type = "ASSIGNMENT_TYPE_ENTIRE_ORG" }
  ]

  objects = {
    PROTECTED_DOMAINS = [cpwai_browse_domains_object.org_internal_domains.id]
  }
}
```

## Argument Reference

The following arguments are supported:

| Argument | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `name` | String | Yes | — | Human-readable rule name. Used for display in the Infinity Portal UI and Terraform state. |
| `description` | String | No | `""` | Free-text explanation of the rule's purpose. |
| `order` | Number | No | Computed | Position in the rulebase (0-based). Lower values are evaluated first. |
| `active` | Boolean | No | `true` | Whether the rule is currently enforced. Set to `false` to disable without deleting. |
| `policy` | String (JSON) | Yes | — | The enforcement policy as a JSON string (use `jsonencode({...})`). See the [Policy Reference](../guides/policy-reference.md) for the Secure Browsing policy schema. |
| `source` | List | No | — | List of users, groups, or the entire organization this rule applies to (source entries). An empty list means the rule matches all users. See [Source](#source) below. |
| `objects` | Map of lists | No | — | Policy objects attached to this rule. A map of object type to a list of object IDs. Valid types: `PROTECTED_DOMAINS`, `EXCLUDE_DOMAINS`. |

### Source

Each element in the `source` list supports:

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `assignment_type` | String | Yes | Assignment scope: `ASSIGNMENT_TYPE_USER`, `ASSIGNMENT_TYPE_GROUP`, or `ASSIGNMENT_TYPE_ENTIRE_ORG`. |
| `assignment_id` | String | No | UUID of the user or group. Not required for `ASSIGNMENT_TYPE_ENTIRE_ORG`. |
| `display_name` | String | No | Display name. Computed by the server if not provided. |

### Policy Fields (Secure Browsing)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `password_reuse` | String | No | Password reuse detection action: `"detect"`, `"prevent"`, `"ask"`, or `"off"`. |
| `zero_phishing` | String | No | Phishing protection action: `"detect"`, `"prevent"`, `"ask"`, or `"off"`. |
| `safe_search` | String | No | Safe search enforcement: `"on"` or `"off"`. |
| `search_reputation` | String | No | Search reputation filtering: `"on"` or `"off"`. |
| `domain_exclusions` | Array | No | List of domains to exclude from secure browsing. See [Domain Item](#domain-item) below. |
| `protected_domains` | Array | No | List of domains to protect with secure browsing. See [Domain Item](#domain-item) below. |

#### Domain Item

Each element in `domain_exclusions` or `protected_domains`:

| Field | Type | Description |
|-------|------|-------------|
| `domain` | String | **Required.** Valid domain name or IPv4 address. |
| `match_mode` | String | `"exact_host"` (default) or `"any_subdomain"`. |

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Unique identifier of the rule (UUID). Assigned by the server on creation. |
| `version` | Number | Rule version number, incremented on every modification. |

## Import

An existing Browse Secure Browsing rule can be imported using its rule ID:

```shell
terraform import cpwai_browse_secure_browsing_rule.example <rule_id>
```
