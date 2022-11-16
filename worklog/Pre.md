
__In HTTP routes, additional placeholders are available (replace any `*`)__ :

| Placeholder                                        | Description                                                                                 |
| -------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| `{http.request.body}`                              | The request body (⚠️ inefficient; use only for debugging)                                    |
| `{http.request.cookie.*}`                          | HTTP request cookie                                                                         |
| `{http.request.duration}`                          | Time up to now spent handling the request (after decoding headers from client)              |
| `{http.request.duration_ms}`                       | Same as 'duration', but in milliseconds.                                                    |
| `{http.request.uuid}`                              | The request unique identifier                                                               |
| `{http.request.header.*}`                          | Specific request header field                                                               |
| `{http.request.host.labels.*}`                     | Request host labels (0-based from right); e.g. for foo.example.com: 0=com, 1=example, 2=foo |
| `{http.request.host}`                              | The host part of the request's Host header                                                  |
| `{http.request.hostport}`                          | The host and port from the request's Host header                                            |
| `{http.request.method}`                            | The request method                                                                          |
| `{http.request.orig_method}`                       | The request's original method                                                               |
| `{http.request.orig_uri.path.dir}`                 | The request's original directory                                                            |
| `{http.request.orig_uri.path.file}`                | The request's original filename                                                             |
| `{http.request.orig_uri.path}`                     | The request's original path                                                                 |
| `{http.request.orig_uri.query}`                    | The request's original query string (without `?`)                                           |
| `{http.request.orig_uri}`                          | The request's original URI                                                                  |
| `{http.request.port}`                              | The port part of the request's Host header                                                  |
| `{http.request.proto}`                             | The protocol of the request                                                                 |
| `{http.request.remote.host}`                       | The host (IP) part of the remote client's address                                           |
| `{http.request.remote.port}`                       | The port part of the remote client's address                                                |
| `{http.request.remote}`                            | The address of the remote client                                                            |
| `{http.request.scheme}`                            | The request scheme                                                                          |
| `{http.request.tls.version}`                       | The TLS version name                                                                        |
| `{http.request.tls.cipher_suite}`                  | The TLS cipher suite                                                                        |
| `{http.request.tls.resumed}`                       | The TLS connection resumed a previous connection                                            |
| `{http.request.tls.proto}`                         | The negotiated next protocol                                                                |
| `{http.request.tls.proto_mutual}`                  | The negotiated next protocol was advertised by the server                                   |
| `{http.request.tls.server_name}`                   | The server name requested by the client, if any                                             |
| `{http.request.tls.client.fingerprint}`            | The SHA256 checksum of the client certificate                                               |
| `{http.request.tls.client.public_key}`             | The public key of the client certificate.                                                   |
| `{http.request.tls.client.public_key_sha256}`      | The SHA256 checksum of the client's public key.                                             |
| `{http.request.tls.client.certificate_pem}`        | The PEM-encoded value of the certificate.                                                   |
| `{http.request.tls.client.certificate_der_base64}` | The base64-encoded value of the certificate.                                                |
| `{http.request.tls.client.issuer}`                 | The issuer DN of the client certificate                                                     |
| `{http.request.tls.client.serial}`                 | The serial number of the client certificate                                                 |
| `{http.request.tls.client.subject}`                | The subject DN of the client certificate                                                    |
| `{http.request.tls.client.san.dns_names.*}`        | SAN DNS names(index optional)                                                               |
| `{http.request.tls.client.san.emails.*}`           | SAN email addresses (index optional)                                                        |
| `{http.request.tls.client.san.ips.*}`              | SAN IP addresses (index optional)                                                           |
| `{http.request.tls.client.san.uris.*}`             | SAN URIs (index optional)                                                                   |
| `{http.request.uri.path.*}`                        | Parts of the path, split by `/` (0-based from left)                                         |
| `{http.request.uri.path.dir}`                      | The directory, excluding leaf filename                                                      |
| `{http.request.uri.path.file}`                     | The filename of the path, excluding directory                                               |
| `{http.request.uri.path}`                          | The path component of the request URI                                                       |
| `{http.request.uri.query.*}`                       | Individual query string value                                                               |
| `{http.request.uri.query}`                         | The query string (without `?`)                                                              |
| `{http.request.uri}`                               | The full request URI                                                                        |
| `{http.response.header.*}`                         | Specific response header field                                                              |
| `{http.vars.*}`                                    | Custom variables in the HTTP handler chain                                                  |
| `{http.shutting_down}`                             | True if the HTTP app is shutting down                                                       |
| `{http.time_until_shutdown}`                       | Time until HTTP server shutdown, if scheduled                                               |