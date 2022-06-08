package linkservice

// function that truncates protocol, www and trailing slashes
func transformLink(link string) string {
	if len(link) >= 7 && link[0:7] == "http://" {
		link = link[7:]
	} else if len(link) >= 8 && link[0:8] == "https://" {
		link = link[8:]
	}
	if len(link) >= 4 && link[0:4] == "www." {
		link = link[4:]
	}
	lastSlash := len(link)
	for lastSlash > 0 && link[lastSlash-1] == '/' {
		lastSlash--
	}
	return link[0:lastSlash]
}
