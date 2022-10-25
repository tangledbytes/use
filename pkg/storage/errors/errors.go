// errors package is a collection of errors that are used
// in the storage package.
package errors

import "fmt"

var (
	ErrStorageNotInitialized = fmt.Errorf("storage is not initialized")
	ErrKeyNotFound           = fmt.Errorf("key not found")
	ErrCorruptStorage        = fmt.Errorf("storage is corrupted")
	ErrReadOnlyStorage       = fmt.Errorf("storage is read only")
)
