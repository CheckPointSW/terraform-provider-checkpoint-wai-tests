# AI Applications Catalog

AI application IDs identify specific AI services (e.g., ChatGPT, Google Gemini, Claude) in access and chats rule policies. These IDs are used in the `genai_application` field within `services_and_application`.

## Usage in Terraform

```hcl
policy = jsonencode({
  action  = "block"
  logging = "enabled"
  services_and_application = {
    mode = "selected"
    genai_application = [
      {
        id   = 1       # Application ID from the Apps Catalog
        mode = "all"
      }
    ]
  }
})
```

## Looking Up Application IDs

Use the Apps Catalog API to search for applications by name, description, or URL:

| Operation | Endpoint | Description |
|-----------|----------|-------------|
| Search apps | `POST /apps/search` | Search AI applications by name, description, or URL. |
| Get apps by IDs | `POST /apps/by-ids` | Retrieve application details for a list of known IDs. |

For the full API reference, see the [Apps Catalog API documentation](https://app.swaggerhub.com/apis/Check-Point/checkpoint-ai-security/1.0.0#/Apps%20Catalog).
