# Resource: cpwai_browse_domains_object

Manages a Domains policy object in Browse Protect. A reusable domain list that can be attached to rules to define which domains a policy applies to.

## Example Usage

```hcl
resource "cpwai_browse_domains_object" "org_internal_domains" {
  name        = "Organization Internal Domains"
  description = "Internal corporate domains that require browsing protection"

  object_data = jsonencode({
    domains = [
      {
        domain     = "portal.myorg.com"
        match_mode = "exact_host"
      },
      {
        domain     = "apps.myorg.com"
        match_mode = "any_subdomain"
      },
      {
        domain     = "intranet.myorg.com"
        match_mode = "any_subdomain"
      }
    ]
  })
}
```

### Using with a Rule

After creating a domains object, attach it to a rule as protected domains via the `objects` attribute:

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

resource "cpwai_browse_secure_browsing_rule" "protect_internal" {
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

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | String | Yes | Human-readable object name. |
| `description` | String | No | Explanation of the object's purpose. |
| `object_data` | String (JSON) | Yes | The object configuration as a JSON string (use `jsonencode({...})`). See [Object Data Fields](#object-data-fields) below. |

### Object Data Fields

The `object_data` JSON object supports the following fields:

| Field | Type | Description |
|-------|------|-------------|
| `domains` | Array | List of domain entries. Each entry has `domain` (string) and optional `match_mode`. |

Each domain entry supports:

| Field | Type | Description |
|-------|------|-------------|
| `domain` | String | The domain name (e.g. `example.com`). |
| `match_mode` | String | `exact_host` (matches only the exact hostname, default) or `any_subdomain` (matches any subdomain). |

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Unique identifier of the object (UUID). Assigned by the server on creation. |

## Import

An existing Domains object can be imported using its object ID:

```shell
terraform import cpwai_browse_domains_object.example <object_id>
```
