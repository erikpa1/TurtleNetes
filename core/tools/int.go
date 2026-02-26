// https://claude.ai/chat/b4f96a73-a7e0-4972-8ada-bb27fb9d9436
package tools

import "strconv"

func StringToInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return result
}
