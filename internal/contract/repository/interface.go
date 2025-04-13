package repository

type IRepository interface {
	// SaveEntity saves the entity with the given key and data in the namespace
	SaveEntity(namespace string, key string, data []byte) error
	// LoadEntity loads the entity with the given key in the namespace
	LoadEntity(namespace string, key string) ([]byte, error)
}
