// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

package structs

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (app *App) ParseRequest(r *http.Request) (string, string, error) {
	name, found := strings.CutPrefix(r.URL.Path, "/")
	if !found {
		return "", "", fmt.Errorf("prefix not found")
	}

	u, err := url.QueryUnescape(name)
	if err != nil {
		return "", "", fmt.Errorf("can not unescape url")
	}

	path := fmt.Sprintf("%s/%s", *app.Mount, u)

	return name, path, nil
}
