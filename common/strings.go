package common

import (
	"fmt"

	"github.com/rs/xid"
)

// GenerateLongUniqueID generate a long string unique id
func GenerateLongUniqueID() string {
	return fmt.Sprintf("%s-%s-%s", xid.New().String(), xid.New().String(), xid.New().String())
}
