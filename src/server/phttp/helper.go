package phttp

func IsHtml5SupportSubtitle(ext string) bool {
	switch ext {
	case "vtt":
		return true
	}
	return false
}
