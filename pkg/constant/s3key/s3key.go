package s3key

import (
	"fmt"
	"strings"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/bucketname"
)

func Image(username string, filename string) string {
	timenow := time.Now().Unix()
	safeFilename := strings.ReplaceAll(filename, " ", "_")
	return fmt.Sprintf("%s/%d_%s", username, timenow, safeFilename)
}

func ImageURL(baseEndpoint string, key string) string {
	return fmt.Sprintf("%s/%s/%s", baseEndpoint, bucketname.Image, key)
}

// KeyFromImageURL is the inverse of ImageURL: it strips the base endpoint and
// bucket prefix so a stored URL can be turned back into an S3 object key.
func KeyFromImageURL(baseEndpoint string, url string) string {
	prefix := fmt.Sprintf("%s/%s/", baseEndpoint, bucketname.Image)
	return strings.TrimPrefix(url, prefix)
}
