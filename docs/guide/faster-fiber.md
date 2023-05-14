---
id: faster-woody
title: ⚡ Make Woody Faster
sidebar_position: 7
---

## Custom JSON Encoder/Decoder
Since Woody v2.32.0, we use **encoding/json** as default json library due to stability and producibility. However, the standard library is a bit slow compared to 3rd party libraries. If you're not happy with the performance of **encoding/json**, we recommend you to use these libraries:
- [goccy/go-json](https://github.com/goccy/go-json)
- [bytedance/sonic](https://github.com/bytedance/sonic)
- [segmentio/encoding](https://github.com/segmentio/encoding)
- [mailru/easyjson](https://github.com/mailru/easyjson)
- [minio/simdjson-go](https://github.com/minio/simdjson-go)
- [wI2L/jettison](https://github.com/wI2L/jettison)

```go title="Example"
package main

import "github.com/ximispot/woody"
import "github.com/goccy/go-json"

func main() {
	app := woody.New(woody.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	# ...
}
```

### References
- [Set custom JSON encoder for client](../api/client.md#jsonencoder)
- [Set custom JSON decoder for client](../api/client.md#jsondecoder)
- [Set custom JSON encoder for application](../api/woody.md#config)
- [Set custom JSON decoder for application](../api/woody.md#config)