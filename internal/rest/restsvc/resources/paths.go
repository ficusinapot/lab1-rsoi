package resources

import "fmt"

const APIPrefix = "/api/v1"

func ResourceLocation(collectionPath string, id int) string {
	return fmt.Sprintf("%s%s/%d", APIPrefix, collectionPath, id)
}
