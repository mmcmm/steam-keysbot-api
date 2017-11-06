package internal

// Err ...
type Err struct {
	Func      uint8  `json:"func" validate:"nonzero"`
	Message   uint8  `json:"message" validate:"nonzero"`
	CreatedAt string `json:"created_at" validate:"nonzero"`
}
