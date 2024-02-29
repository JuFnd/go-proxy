package scanner

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func DirBusterScan(baseURL string, dictFilePath string) map[string]bool {
    foundFiles := make(map[string]bool)

    data, err := ioutil.ReadFile(dictFilePath)
    if err != nil {
        fmt.Println("Failed to read dictionary file:", err)
        return foundFiles
    }

    lines := strings.Split(string(data), "\n")
    u, _ := url.Parse(baseURL)
    for _, line := range lines {
        lineTrimmed := strings.TrimSpace(line)
        if len(lineTrimmed) > 0 {
            u.Path = lineTrimmed
            resp, err := http.Head(u.String())
            if err == nil && resp.StatusCode != 404 {
                foundFiles[lineTrimmed] = true
            }
        }
    }

    return foundFiles
}

func Scan(request *http.Request, dictFilePath string) map[string]bool {
    parsedDictFilePath := strings.ReplaceAll(dictFilePath, "\\", "/")
    foundFiles := DirBusterScan(request.URL.String(), parsedDictFilePath)
    return foundFiles
}