package url

import (
	"fmt"
	"github.com/kataras/iris/core/errors"
	"regexp"
	"strings"
)

type QueryMap map[string]interface{}

type UserInfo struct {
	Username string
	Password string
}

type URL struct {
	Scheme      string
	UserInfo    *UserInfo
	Host        string
	Domain      string
	Subdomain   string
	Port        string
	Path        string
	QueryString string
	Query       QueryMap
	Fragment    string
}

func (u *URL) Uri(path string, params QueryMap) *URL {
	u.ReplacePath(path)

	return u
}

func (u *URL) UriAppend(path string, params QueryMap) *URL {
	u.AppendPath(path)
	u.parseParams(params)

	return u
}

func (u *URL) AppendPath(path string) *URL {
	u.Path = fmt.Sprintf(`%s/%s`, strings.TrimRight(u.Path, "/"), strings.TrimLeft(path, "/"))

	return u
}

func (u *URL) ReplacePath(path string) *URL {
	u.Path = strings.TrimLeft(path, "/")

	return u
}

func (u *URL) GetHost() string {
	return u.Host
}

func (u *URL) GetDomain() string {
	return u.Domain
}

func (u *URL) GetPath() string {
	return u.Path
}

func (u *URL) GetBaseDomain() string {
	return fmt.Sprintf("www.%s", u.Domain)
}

func (u *URL) GetUrlString() string {
	urlString := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	if u.Port != "" {
		urlString = fmt.Sprintf("%s:%s", urlString, u.Port)
	}

	if u.Path != "" {
		urlString = fmt.Sprintf("%s/%s", urlString, u.Path)
	}

	if u.QueryString != "" {
		urlString = fmt.Sprintf("%s?%s", urlString, u.QueryString)
	}

	if u.Fragment != "" {
		urlString = fmt.Sprintf("%s#%s", urlString, u.Fragment)
	}

	return urlString
}

func (u *URL) String() string {
	return u.GetUrlString()
}

func (u *URL) parseParams(params QueryMap) {
	u.Query = params

	var queryString string
	for key, value := range params {
		if v, ok := value.([]interface{}); ok {
			for _, val := range v {
				queryString = fmt.Sprintf("&%s[]=%v", key, val)
			}
		} else {
			queryString = fmt.Sprintf("&%s=%v", key, value)
		}
	}

	u.QueryString = queryString
}

// Parse the given url into URL object
func Parse(url string) (*URL, error) {
	// URL regex `^[scheme]?[user_info]?[host][port]?[path]?[query]?#?[fragment]?`
	regexRule := fmt.Sprintf(`(?i)^%s?%s?%s`, getSchemeRegexRule(), getMainRegexRule(), getUrlTrailRegexRule())
	regex, _ := regexp.Compile(regexRule)

	matches := regex.FindStringSubmatch(url)

	// If given url matches the regex rule, result will have all 11 group
	if len(matches) != 11 {
		return nil, errors.New(fmt.Sprintf("%s is an invalid url!", url))
	}

	urlObject := &URL{
		Scheme:      matches[1],
		UserInfo:    &UserInfo{Username: matches[2], Password: matches[3]},
		Host:        matches[4],
		Domain:      matches[6],
		Subdomain:   matches[5],
		Port:        matches[7],
		Path:        matches[8],
		QueryString: matches[9],
		Fragment:    matches[10],
	}

	return urlObject, nil
}

func getSchemeRegexRule() string {
	return `(?:(https?):(?:\/\/)?)`
}

func getMainRegexRule() string {
	return fmt.Sprintf(`%s?%s%s`, getUserInfoRegexRule(), getHostRegexRule(), getPortRegexRule())
}

func getUserInfoRegexRule() string {
	return `(?:([a-zA-z0-9\-]+):([a-zA-z0-9\-]*)?@)`
}

func getHostRegexRule() string {
	return `((?:([a-zA-z0-9\-\.]+)\.)?([a-zA-z0-9\-]+\.[a-zA-z]{2,}))`
}

func getPortRegexRule() string {
	return `(?::(\d+))`
}

func getUrlTrailRegexRule() string {
	return fmt.Sprintf(`%s?%s?#?%s?`, getPathRegexRule(), getQueryRegexRule(), getFragmentRegexRule())
}

func getPathRegexRule() string {
	return `(?:\/([\+~%\/\.\w\-_]*)`
}

func getQueryRegexRule() string {
	return `(?:\?([\-\+=&;%@\.\w_]*)?)`
}

func getFragmentRegexRule() string {
	return `([\.\!\/\\\w\-_%]*)?)`
}
