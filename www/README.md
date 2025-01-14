# LOGJAM www/

The `www/` source-code directory-paths contains HTTP handlers.

The directory paths roughly correspond to the HTTP request-paths.

So, for example, the source-code path `www/hapi/v1/ws/servehttp.go` corresponds to the HTTP request-path `/hapi/v1/ws`

And, for example, the fictitious source-code path `www/apple/banana/cherry/servehttp.go` would correspond to the HTTP request-path `/apple/banana/cherry`

## WILDCARD

The exception to this is `www/FALLBACK/servehttp.go`, which is what is called when with the following routes:

- Any path that begins with `/log/`
- Any matches this pattern `/@{handle}/host`

## Usage

To have all the code under this directory activate and work, something needs to `import` it — even if it is just for the side-effects.
I.e.:

```golang
import _ "github.com/reiver/logjam/www"
```

Note the underscore (_).
