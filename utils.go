package main

func strGetFirstOne(vars ...string) string {
	for _, v := range vars {
		if v != "" {
			return v
		}
	}
	return ""
}
