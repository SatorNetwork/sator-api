package storage

// Options struct
type Options struct {
	Key            string
	Secret         string
	Endpoint       string
	Region         string
	Bucket         string
	URL            string
	ForcePathStyle bool
	DisableSSL     bool
}
