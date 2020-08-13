## pulse
HTTP client for Feedzai's API

#### usage

```
import "github.com/jcaberio/go-pulse"

client, err := pulse.New(&pulse.Options{
    Username:   "firstname.lastname@paymaya.com",
    Password:   "yourpassword",
    BaseURL:    "https://feedzai-pulse-stg.voyagerinnovation.com",
    Timeout:    1 * time.Minute,
})

client.UploadList("list.csv", "listid_123")
```
