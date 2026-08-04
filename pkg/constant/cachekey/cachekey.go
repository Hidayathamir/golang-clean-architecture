package cachekey

import "fmt"

func UserID(id int64) string {
	return fmt.Sprintf("user:%d", id)
}
