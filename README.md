## pulse
HTTP client for Feedzai's API

#### usage

```
import "code.corp.voyager.ph/jorick.caberio/go-pulse"

client, err := pulse.NewClient(&pulse.Options{
    Username:   "firstname.lastname@paymaya.com",
    Password:   "yourpassword",
    BaseURL:    "https://feedzai-pulse-stg.voyagerinnovation.com",
    Timeout:    1 * time.Minute,
})

client.UploadList("list.csv", "listid_123")
```
