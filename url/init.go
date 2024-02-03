package url

import "regexp"

type URL struct {
	raw     string
	matches []string
}

func NewURL(raw string) *URL {
	pattern := regexp.MustCompile(`^(https?://)(?:www\.)?([^/]*)([^?#]*)(\?[^#]*)?(\#.*)?$`)
	matches := pattern.FindStringSubmatch(raw)
	if matches == nil {
		return nil
	}

	return &URL{
		raw:     raw,
		matches: matches,
	}
}

func (u *URL) Host() string {
	return u.matches[2]
}

func (u *URL) Path() string {
	return u.matches[3]
}

func (u *URL) Query() string {
	return u.matches[4]
}

func (u *URL) Fragment() string {
	return u.matches[5]
}

func (u *URL) Protocol() string {
	return u.matches[1]
}

func (u *URL) SameHostAs(other *URL) bool {
	return u.Host() == other.Host()
}

func (u *URL) String() string {
	return u.raw
}
