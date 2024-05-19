# Socktopus

!["Logo"](https://raw.githubusercontent.com/Zuma206/socktopus/main/logo.png)

## Mission Statement

A simple websocket microservice made to be paired with serverless applications. It supports multiple secrets, allowing one instance to be used in multiple applications. It's written in Go to make the best possible use of the resources you give it.

## Authenticating

To authenticate as a server, it's very simple. First, add an environment variable called `YOURAPP_SOCKTOPUS_SECRET` with a value generated from the server's homepage. Then, when calling `/send` or `/kick` you just need to provide `yourapp` and the secret value to authenticate.

To authenticate as a client, things are less simple. You'll need to contact the application's server, and recieve a token. The token will look like this:

`hex(timestampNow) + 'h' + hex(appName) + 'h' + hex(connectionId) + 'h' + hex(sha256Sign(timestamp+connectionId, secret))`

To easily generate this token, use the SDK. If your app isn't written in TS/JS, use a trusted crypto library instead.

## API Reference

```http
GET /recieve
```

Parameter Type: Query Parameters
| Parameter | Type | Description |
| :-------- | :------- | :--------------------------------------------------------- |
| `token` | `string` | A token issued by the main application server (see above) |

```http
POST /send
```

Parameter Type: JSON
| Parameter | Type | Description |
| :-------- | :------- | :--------------------------------------------------------- |
| `secretName` | `string` | The name of the secret that will be provided |
| `secret` | `string` | The secret value stored under the secret name that is provided |
| `messages` | `array` | A list of objects, each containing a `recipient` string and a `content` string. Each message `content` will be send to the `recipient` connectionId |

```http
POST /kick
```

Parameter Type: JSON
| Parameter | Type | Description |
| :-------- | :------- | :--------------------------------------------------------- |
| `secretName` | `string` | The name of the secret that will be provided |
| `secret` | `string` | The secret value stored under the secret name that is provided |
| `connections` | `array` | A list of connections to terminate. This endpoint may take up to 5 seconds, as the server will first wait for the connection's tokens to expire so that they cannot reconnect|

## License

This project is licensed to you under the MIT License
