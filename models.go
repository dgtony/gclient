package gclient

// Here are request/response models defined to be
// used in client's top level methods: Get/Set/etc.

/* interfaces */

type ValueReader interface {
	GetValue() interface{}
}

type KeyReader interface {
	GetKeys() []string
}

/* requests */

type GetKeyRequest struct {
	Key string `json:"key"`
}

type GetKeySubKeyRequest struct {
	Key    string `json:"key"`
	SubKey string `json:"subkey"`
}

type GetKeySubIndexRequest struct {
	Key      string `json:"key"`
	SubIndex int    `json:"subindex"`
}

type SetKeyRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl"`
}

type RemoveKeyRequest struct {
	Key string `json:"key"`
}

type GetStoredKeysRequest struct {
	Mask string `json:"mask"`
}

/* succeeded responses */

type GetKeyResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func (r *GetKeyResponse) GetValue() interface{} {
	return r.Value
}

type SetKeyResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	TTL   int         `json:"ttl"`
}

func (r *SetKeyResponse) GetValue() interface{} {
	return r.Value
}

type GetKeySubResponse struct {
	Key      string      `json:"key"`
	SubKey   string      `json:"subkey,omitempty"`
	SubIndex int         `json:"subindex,omitempty"`
	Value    interface{} `json:"value"`
}

func (r *GetKeySubResponse) GetValue() interface{} {
	return r.Value
}

type GetStoredKeysResponse struct {
	Mask string   `json:"mask"`
	Keys []string `json:"keys"`
}

func (r *GetStoredKeysResponse) GetKeys() []string {
	return r.Keys
}

/* error response */

type ErrorResponse struct {
	Code   int
	Reason string
}
