package validator

import (
	"errors"
	"net/url"
	"strings"
)

func ValidateURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return errors.New("url is required")
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return errors.New("invalid url")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("url must use http or https scheme")
	}

	if u.Host == "" {
		return errors.New("url must have a host")
	}

	return nil
}
