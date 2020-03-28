package main

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

// indexVars holds the values passed to the index template.
type indexVars struct {
	Query    string
	ImageURL string
}

// tpl is the html to show at the root of the domain.
var tpl = `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta property="og:url" content="https://fournaan.com{{.ImageURL}}" />
    <meta property="og:image" content="https://fournaan.com{{.ImageURL}}" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:site" content="@mhemmings" />
    <meta name="twitter:title" content="Four naan, Jeremy?" />
    <meta name="twitter:image" content="https://fournaan.com{{.ImageURL}}" />

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css" integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">

    <title>Four naan, Jeremy?</title>
    <style>
      html {
        position: relative;
        min-height: 100%;
      }
      body {
        margin-bottom: 95px; /* Margin bottom by footer height */
        background: #f8f9fa;
      }
      .container{
        // padding-top: 25px;
        text-align: center;
        width: auto;
        max-width: 680px;
        padding: 25px 15px;
      }
      img {
        max-width: 100%;
        border-radius: 15px;
      }
      .input-group {
        margin-top: 25px;
      }
      .footer {
        position: absolute;
        bottom: 0;
        width: 100%;
        height: 40px; /* Set the fixed height of the footer here */
        line-height: 40px; /* Vertically center the text there */
        background-color: #ffffff;
        text-align: center;
        font-size: 10pt;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-10">
         <img src="{{.ImageURL}}">
        </div>
      </div>
      <div class="row justify-content-center">
        <div class="col-6">
          <form method="GET">
            <div class="input-group input-group-lg">
              <input type="text" class="form-control"
                name="four" value="{{.Query}}" maxlength="15">
              <div class="input-group-append">
                <button class="btn btn-outline-primary" type="submit">Create</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
    <footer class="footer">
       <a href="https://github.com/mhemmings/fournaan"><img height="20" width="20" src="https://cdn.jsdelivr.net/npm/simple-icons@latest/icons/github.svg"> Source</a> |
       <span class="text-muted">Made by <a href="https://twitter.com/mhemmings">Mark Hemmings</a> |
       Peep Show © Channel Four Television Corporation
       </span>
    </footer>
  </body>
</html>
`

var indexHtml = template.Must(template.New("index").Parse(tpl))

// indexHandler writes the index template to the http.ResponseWriter.
func indexHandler(w http.ResponseWriter, req *http.Request) {
	q := strings.TrimSpace(req.URL.Query().Get("four"))
	if q == "" {
		q = "naan"
	}

	u, err := url.Parse("/img/" + q)
	if err != nil {
		http.Error(w, "error creating image url", http.StatusInternalServerError)
		return
	}

	num := req.URL.Query().Get("num")
	name := req.URL.Query().Get("name")
	imgVals, _ := url.ParseQuery(u.RawQuery)
	if num != "" {
		imgVals.Set("num", num)
	}
	if name != "" {
		imgVals.Set("name", name)
	}
	u.RawQuery = imgVals.Encode()

	vars := indexVars{
		Query:    q,
		ImageURL: u.String(),
	}
	err = indexHtml.Execute(w, vars)
	if err != nil {
		http.Error(w, "error exectuting index template", http.StatusInternalServerError)
	}
}
