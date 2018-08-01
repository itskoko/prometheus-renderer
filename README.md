# prometheus-renderer
*Renders Prometheus queries as PNG images*

## Usage
Commands are located in `cmd/`

### render
render is a CLI to render a query as png image.

```
Usage of render
  -f string
    	Path to output file (default "out.png")
  -h int
    	Height (default 600)
  -s duration
    	Graph range (default 1h0m0s)
  -u string
    	URL of prometheus server (default "http://localhost:9090")
  -w int
    	Width (default 800)
```

### renderd
renderd is a HTTP server returning an png image for a query.

```
Usage of renderd
  -l string
    	Address to listen on (default ":8080")
  -u string
    	URL of prometheus server (default "http://localhost:9090")
```

#### HTTP API
The `/graph` GET endpoint takes the following arguments and returns an PNG image:

- `q`: Prometheus query (mandatory)
- `h`: Height (360)
- `w`: Width (360)
- `s`: Graph range in seconds (3600)
