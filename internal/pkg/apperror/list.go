package apperror

import(
	"github.com/pkg/errors"
)

// ErrNotFound is error for case when smth not found
var ErrNotFound error = errors.New("Not found")

