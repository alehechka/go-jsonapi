package jsonapi

// // Node is used to represent a generic JSON API Resource
// type Node struct {
// 	Type          string                 `json:"type"`
// 	ID            string                 `json:"id,omitempty"`
// 	ClientID      string                 `json:"client-id,omitempty"`
// 	Attributes    map[string]interface{} `json:"attributes,omitempty"`
// 	Relationships map[string]interface{} `json:"relationships,omitempty"`
// 	Links         *LinkMap               `json:"links,omitempty"`
// 	Meta          *Meta                  `json:"meta,omitempty"`
// }

// // Payload is used to encapsulate the One and Many payload types
// type Payload interface {
// 	clearIncluded()
// }

// // OnePayload is used to represent a generic JSON API payload where a single
// // resource (Node) was included as an {} in the "data" key
// type OnePayload struct {
// 	Data     *Node    `json:"data"`
// 	Included []*Node  `json:"included,omitempty"`
// 	Links    *LinkMap `json:"links,omitempty"`
// 	Meta     *Meta    `json:"meta,omitempty"`
// }

// func (p *OnePayload) clearIncluded() {
// 	p.Included = []*Node{}
// }

// // ManyPayload is used to represent a generic JSON API payload where many
// // resources (Nodes) were included in an [] in the "data" key
// type ManyPayload struct {
// 	Data     []*Node  `json:"data"`
// 	Included []*Node  `json:"included,omitempty"`
// 	Links    *LinkMap `json:"links,omitempty"`
// 	Meta     *Meta    `json:"meta,omitempty"`
// }

// func (p *ManyPayload) clearIncluded() {
// 	p.Included = []*Node{}
// }

// func Marshal(models interface{}) (Payload, error) {
// 	return nil, nil
// }
