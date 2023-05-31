# proxy-go
Proxy for reproduction of bug which include change of response to specific request

## Run proxy with Podman
	podman-compose -f compose.yaml up -d
Tested on versions:

	podman version: 4.4.1
	podman-composer version  1.0.3

## Connect insights-client through proxy
Add to insights-client.conf line:

	proxy=http://<IP address>:3129
Register insights-client:

	insights-client --register
