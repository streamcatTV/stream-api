package models

// Stream contains stream models.
type Stream struct {
	ID          int64  `db:"id" json:"id,omitempty"`
	StreamName  string `db:"stream_name" json:"stream_name,omitempty" valid:"alphanum,required"`
	Type        string `db:"type" json:"type,omitempty" valid:"alphanum,required"`
	Description string `db:"description" json:"description,omitempty" valid:"alphanum"`
	URL         string `db:"url" json:"url,omitempty" valid:"url,required"`
	SecretKey   string `db:"secret_key" json:"secret_key,omitempty" valid:"alphanum,required"`
	PublicKey   string `db:"public_key" json:"public_key,omitempty" valid:"alphanum,required"`
	Private     bool   `db:"private" json:"private,omitempty" valid:"required"`
}
