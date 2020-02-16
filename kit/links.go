package kit

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinks makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links (to same domain) in the
// HTML document.
func ExtractLinks(url string) (links []string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	domain := hostDomain(resp.Request.URL)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getting %s: %s", url, http.StatusText(resp.StatusCode))
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}

				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					log.Printf("parsing URL %s: %v", link, err)
					continue // ignore bad URLs
				}

				if strings.Contains(link.Hostname(), domain) {
					links = append(links, link.String())
				}
			}
		}
	}

	forEachNode(doc, visitNode)
	return links, nil
}

// Apply f() recursively to a parent HTML node and all its children.
func forEachNode(n *html.Node, f func(n *html.Node)) {
	if f != nil {
		f(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, f)
	}
}

// Return primary domain name of URL.
func hostDomain(url *url.URL) string {
	domains := strings.Split(url.Hostname(), ".")
	return strings.Join(domains[len(domains)-2:], ".")
}
