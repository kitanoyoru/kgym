package key

import "crypto"

type Key struct {
	ID        string            `json:"kid"`
	Private   crypto.PrivateKey `json:"private"`
	Public    crypto.PublicKey  `json:"public"`
	Algorithm string            `json:"alg"`
	Active    bool              `json:"active"`
}
