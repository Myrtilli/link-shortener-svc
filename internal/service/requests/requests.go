package requests

import (
	"net/url"
	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreateShortLinkRequest struct {
	URL string `json:"url"`
}

func (r *CreateShortLinkRequest) Validate() error {
	if strings.TrimSpace(r.URL) == "" {
		return errors.New("URL field is required")
	}

	if _, err := url.ParseRequestURI(r.URL); err != nil {
		return errors.Wrap(err, "URL field is not a valid URI")
	}

	return nil
}
