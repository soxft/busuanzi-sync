package transfer

import "fmt"

func GetSyncKey() string {
	return fmt.Sprintf("%s:s:synced", Prefix)
}
