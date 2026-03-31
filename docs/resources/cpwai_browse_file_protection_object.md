# Resource: cpwai_browse_file_protection_object

Manages a File Protection policy object in Browse Protect. A reusable file handling configuration that can be attached to rules to control how file uploads and downloads are processed.

## Example Usage

```hcl
resource "cpwai_browse_file_protection_object" "internal_file_policy" {
  name        = "Internal File Protection"
  description = "Emulate downloads and block unsupported uploads on internal domains"

  object_data = jsonencode({
    file_download_action_unsupported = "Block"
    file_upload_action_unsupported   = "Block"
    file_download_size_limit_mb      = 100
    file_upload_size_limit_mb        = 50
    file_download_error_fail_close   = "FailClose"
    file_upload_error_fail_close     = "FailClose"
  })
}
```

### Using with a Rule

After creating a file protection object, attach it to a secure browsing rule to enable file emulation on internal domains:

```hcl
resource "cpwai_browse_file_protection_object" "internal_file_policy" {
  name        = "Internal File Protection"
  description = "Emulate downloads and block unsupported uploads on internal domains"

  object_data = jsonencode({
    file_download_action_unsupported = "Block"
    file_upload_action_unsupported   = "Block"
    file_download_size_limit_mb      = 100
    file_upload_size_limit_mb        = 50
    file_download_error_fail_close   = "FailClose"
    file_upload_error_fail_close     = "FailClose"
  })
}

resource "cpwai_browse_secure_browsing_rule" "protect_internal" {
  name        = "Protect Internal Domains"
  description = "Enable secure browsing and file emulation for internal corporate domains"
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
    FILE_EMULATION = [cpwai_browse_file_protection_object.internal_file_policy.id]
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
| `file_download_action_unsupported` | String | Action for unsupported file types on download: `"Block"` or `"Allow"`. |
| `file_upload_action_unsupported` | String | Action for unsupported file types on upload: `"Block"` or `"Allow"`. |
| `file_download_size_limit_mb` | Number | Maximum file size in MB for downloads (1-100). |
| `file_upload_size_limit_mb` | Number | Maximum file size in MB for uploads (1-100). |
| `file_download_error_fail_close` | String | Behavior on download emulation error: `"FailClose"` or `"FailOpen"`. |
| `file_download_size_limit_fail_close` | String | Behavior when download exceeds size limit: `"FailClose"` or `"FailOpen"`. |
| `file_upload_error_fail_close` | String | Behavior on upload error: `"FailClose"` or `"FailOpen"`. |
| `file_upload_size_limit_fail_close` | String | Behavior when upload exceeds size limit: `"FailClose"` or `"FailOpen"`. |
| `file_download_action_per_supported_file` | Map | Per-file-type download actions for supported file types. Map of file type to action (`"Block"`, `"Allow"`, `"WaitForEmulation"`, `"EmulateInBackground"`, `"Detect"`, `"Extract"`). |
| `file_download_action_per_unsupported_file` | Map | Per-file-type download actions for unsupported file types. Map of extension to action (`"Block"`, `"Allow"`). |
| `file_upload_action_per_supported_file` | Map | Per-file-type upload actions for supported file types. Map of file type to action. |
| `file_upload_action_per_unsupported_file` | Map | Per-file-type upload actions for unsupported file types. Map of extension to action (`"Block"`, `"Allow"`). |

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

| Attribute | Type | Description |
|-----------|------|-------------|
| `id` | String | Unique identifier of the object (UUID). Assigned by the server on creation. |

## Import

An existing File Protection object can be imported using its object ID:

```shell
terraform import cpwai_browse_file_protection_object.example <object_id>
```
