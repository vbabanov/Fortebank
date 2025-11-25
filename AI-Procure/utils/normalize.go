package utils

import (
	"regexp"
	"strings"
)

var legalPrefixes = regexp.MustCompile(`(?i)\b(ТОО|АО|ООО|ОАО|ЗАО|ПАО|ИП|ЧП|КХ|ПК|LLP|LTD|JSC|CJSC|LLC)\b\s*|` + `[\-\s_\.,\(\)\{\}\][\]\/\\]+`)

func NormalizeString(s string) string {
	s = strings.ToLower(s)
	s = legalPrefixes.ReplaceAllString(s, " ")
	s = strings.Join(strings.Fields(s), " ")
	s = strings.TrimSpace(s)
	return s
}
