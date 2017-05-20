global
    maxconn 4096

defaults
    log global
    mode    http
    option  httplog
    option  dontlognull
    retries 3
    option redispatch
    maxconn 2000
    timeout connect 5000
    timeout client 50000
    timeout server 50000

frontend http-in
    bind :8000
    default_backend http

backend http
{{range .}}     server {{.Name}} {{.Host}}:{{.Port}} maxconn 32
{{end}}
