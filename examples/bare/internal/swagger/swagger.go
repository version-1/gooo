package swagger

import (
	"embed"
	"fmt"
	"io/fs"
)

const hostURL = "http://localhost:8080"

func Index() []byte {
	return []byte(fmt.Sprintf(`<!DOCTYPE html>
            <html>
            <head>
              <meta charset="UTF-8">
              <title>Swagger</title>
              <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
               <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI" />
            </head> 
            <body>
              <div id="swagger-ui"></div>
              <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
              <script>
                window.onload = () => {
                  window.ui = SwaggerUIBundle({
                    url: '%s/api/v1/swagger/swagger.yml',
                    dom_id: '#swagger-ui',
                  });
                };
              </script>
            </body>
          </html>
        `, hostURL))

}

//go:embed *.yml
var swaggerConf embed.FS

func SwaggerYAML() ([]byte, error) {
	f, err := fs.ReadFile(swaggerConf, "swagger.yml")
	if err != nil {
		return f, err
	}

	return f, nil
}
