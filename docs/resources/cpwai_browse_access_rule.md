# Resource: cpwai_browse_access_rule

Manages an Access rule in the Browse Protect policy. Controls which web applications and domains users can access through the browser.

## Example Usage

```hcl
resource "cpwai_browse_access_rule" "allow_internal_access" {
  name        = "Allow Access to Internal Resources"
  description = "Allow employees to access internal corporate domains"
  order       = 0
  active      = true

  policy = jsonencode({
    action  = "allow"
    logging = "enabled"
    services_and_application = {
      mode = "selected"
      domain = [
        {
          domain     = "portal.myorg.com"
          match_mode = "exact_host"
        },
        {
          domain     = "apps.myorg.com"
          match_mode = "any_subdomain"
        }
      ]
    }
  })

  source = [
    {
      assignment_type = "ASSIGNMENT_TYPE_ENTIRE_ORG"
    }
  ]
}
```

-> **Note:** Browse access rules use `domain`, `url`, and `category` in `services_and_application`. To target specific AI applications (via `genai_application`, see [AI Applications Catalog](../guides/ai-applications-catalog.md)), use the Workforce AI access rule (`cpwai_workforce_ai_access_rule`) instead - Workforce AI rules are also available under the Browse product.

### Example with Domains Object

Create a reusable domains object and attach it to a rule using the `objects` attribute:

```hcl
resource "cpwai_browse_domains_object" "org_internal_domains" {
  name        = "Organization Internal Domains"
  description = "Internal corporate domains that require access control"

  object_data = jsonencode({
    domains = [
      { domain = "portal.myorg.com",   match_mode = "exact_host" },
      { domain = "apps.myorg.com",     match_mode = "any_subdomain" },
      { domain = "intranet.myorg.com", match_mode = "any_subdomain" }
    ]
  })
}

resource "cpwai_browse_access_rule" "allow_internal_with_object" {
  name        = "Allow Access to Internal Resources"
  description = "Allow employees to access internal corporate domains"
  order       = 0
  active      = true

  policy = jsonencode({
    action  = "allow"
    logging = "enabled"
    services_and_application = { mode = "all" }
  })

  source = [
    { assignment_type = "ASSIGNMENT_TYPE_ENTIRE_ORG" }
  ]

  objects = {
    DOMAIN_SERVICES = [cpwai_browse_domains_object.org_internal_domains.id]
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
| `policy` | String (JSON) | Yes | — | The enforcement policy as a JSON string (use `jsonencode({...})`). See the [Policy Reference](../guides/policy-reference.md) for the Access policy schema. |
| `source` | List | No | — | List of users, groups, or the entire organization this rule applies to (source entries). An empty list means the rule matches all users. See [Source](#source) below. |
| `objects` | Map of lists | No | — | Policy objects attached to this rule. A map of object type to a list of object IDs. Valid types: `FILE_EMULATION`, `DOMAIN_SERVICES`. |

### Source

Each element in the `source` list supports:

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `assignment_type` | String | Yes | Assignment scope: `ASSIGNMENT_TYPE_USER`, `ASSIGNMENT_TYPE_GROUP`, or `ASSIGNMENT_TYPE_ENTIRE_ORG`. |
| `assignment_id` | String | No | UUID of the user or group. Not required for `ASSIGNMENT_TYPE_ENTIRE_ORG`. |
| `display_name` | String | No | Display name. Computed by the server if not provided. |

### Policy Fields (Access)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | String | Yes | Access action: `"block"`, `"allow"`, or `"ask"`. |
| `logging` | String | No | `"enabled"` or `"disabled"`. |
| `services_and_application` | Object | No | Service scope. See below. |
| `download_file_protection` | String | No | File download protection action. Must not be `"na"` for Browse rules. Valid values: `"Block"`, `"Allow"`, `"EmulateInBackground"`, `"WaitForEmulation"`, `"Detect"`, `"Extract"`. |
| `upload_file_protection` | String | No | File upload protection action. Must not be `"na"` for Browse rules. Valid values: `"Block"`, `"Allow"`, `"EmulateInBackground"`, `"WaitForEmulation"`, `"Detect"`. |

#### services_and_application

| Field | Type | Description |
|-------|------|-------------|
| `mode` | String | `"all"` or `"selected"`. |
| `domain` | Array | List of domains. Each: `{ domain = "...", match_mode = "exact_host" \| "any_subdomain" }`. |
| `url` | Array | List of URLs. Each: `{ url = "https://..." }`. |
| `category` | Array | List of URLF categories. Each: `{ category_id = <int> }`. Any category except Generative AI Tools (60531762). See [Supported Categories](#supported-categories). |

~> **Browse constraint:** Browse access rules cannot use `genai_application` in `services_and_application`. To target specific AI applications, use the Workforce AI access rule (`cpwai_workforce_ai_access_rule`) instead.

#### Supported Categories

Common URLF categories that can be used in Browse rules:

| Category | `category_id` |
|----------|---------------|
| Uncategorized | `0` |
| Computers & Internet | `1` |
| Education | `3` |
| Entertainment | `5` |
| Financial Services | `7` |
| Gambling | `9` |
| Government / Military | `11` |
| Hacking | `13` |
| Health | `15` |
| Illegal / Questionable | `17` |
| Job Search / Careers | `21` |
| News / Media | `24` |
| Phishing | `31` |
| Pornography | `33` |
| Real Estate | `34` |
| Shopping | `42` |
| Business / Economy | `51` |
| Sports | `54` |
| Social Networking | `52000069` |
| File Storage / Sharing | `52000136` |
| Search Engines / Portals | `52000132` |
| Artificial Intelligence (AI) | `60517663` |

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Unique identifier of the rule (UUID). Assigned by the server on creation. |
| `version` | Number | Rule version number, incremented on every modification. |

## Import

An existing Browse Access rule can be imported using its rule ID:

```shell
terraform import cpwai_browse_access_rule.example <rule_id>
```
