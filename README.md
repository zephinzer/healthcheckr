# Healthcheckr

- [Healthcheckr](#healthcheckr)
- [Usage](#usage)
  - [Job mode](#job-mode)
  - [Worker mode](#worker-mode)
    - [Worker configuration documentation](#worker-configuration-documentation)
      - [Configuration properties](#configuration-properties)
      - [Channel type configuration](#channel-type-configuration)
  - [Debug mode](#debug-mode)
    - [Telegram chat ID retrieval](#telegram-chat-id-retrieval)
- [Development](#development)
  - [Getting started](#getting-started)
  - [Executing via `go`](#executing-via-go)

# Usage

## Job mode

To run `healthcheckr` in job mode:

```sh
# with defaults
healthcheckr verify http;

# with all available flags
healthcheckr verify http \
  --expect-body-regex 'google' \
  --expect-response-time-ms 1000 \
  --expect-status-code 200 \
  --log-level 4 \
  --use-hostname yahoo.com \
  --use-method get \
  --use-path / \
  --use-query abc=def \
  --use-query ghi=jkl \
  --use-scheme https \
  --use-user-agent healthcheckr/example;

# get help on available flags
healthcheckr verify http --help;
```

## Worker mode

`healthcheckr` can also run as a long-running worker process with multiple healthchecks configured.

```sh
# with defaults
healthcheckr start worker;

# with all available flags
healthcheckr start worker \
  --config-path /path/to/config/file.yaml \
  --server-addr 0.0.0.0 \
  --server-port 8080;

# get help on flags
healthcheckr start worker --help;
```

Configuration is done via YAML:

```yaml
http:
  - scheme: https
    hostname: google.com
    path: /
    queries:
      - aaa=bbb
      - bbb=ccc
    method: get
    userAgent: healthcheckr/1.0/example
    timeoutMs: 3000
    expectStatusCode: 200
    expectBodyRegexes:
      - Google
    intervalMs: 5000
    failureThreshold: 3
    channels:
      - telegram-default
channels:
  - name: telegram-default
    type: telegram
    apiKey:
      fromEnv: TELEGRAM_BOT_TOKEN
    chatId:
      value: "123456789"
  - name: slack-default
    type: slack
    url:
      fromEnv: SLACK_WEBHOOK_URL
```

An example is available at [`examples/config.yaml`](examples/config.yaml).

### Worker configuration documentation

#### Configuration properties

| Property | Type | Description |
| --- | --- | --- |
| `http[]` | `list(object)` | `http` is the root level property that defines a list of HTTP-based checks. |
| `http[].scheme` | `string` | This defines the scheme to use for the TCP connection. Only `http` and `https` is supported. |
| `http[].hostname` | `string` | This defines the hostname component of the URL. |
| `http[].path` | `string` | This defines the path component of the URL. Defaults to `"/"` |
| `http[].queries` | `list(string)` | This defines a list of `key=value` strings which are used as the query value of the HTTP-based request |
| `http[].method` | `string` | This defines the method to use for the request. Defaults to a `"GET"` request. |
| `http[].userAgent` | `string` | This defines a custom User Agent string to use for the request. Defaults to `"healthcheckr/1.0"` |
| `http[].timeoutMs` | `string` | This defines a timeout for the request in terms of milliseconds. Defaults to `5000` milliseconds. |
| `http[].expectStatusCode` | `integer` | This defines the expected integer status code of the response. Defaults to `200`. |
| `http[].expectBodyRegexes` | `list(string)` | This defines a list of regular expressions to match against the response body. |
| `http[].intervalMs` | `integer` | This defines the interval between checks in terms of milliseconds. Defaults to `5000` milliseconds |
| `http[].failureThreshold` | `integer` | This defines the maximum number of check failures before a notification is triggered to one of the defined channels. |
| `http[].alertMinimumIntervalS` | `integer` | This defines the minimum duration between alerts in terms of seconds. Assuming a value of `60`, this means failure notifications will happen only once every 60 seconds even if failures beyond the `.failureThreshold` |
| `http[].channels` | `list(string)` | This defines a list of channel names for which this HTTP check should notify upon failure/resolution. This string should be a value defined in one `channels[].name` otherwise an error will be thrown. |
| `channels[]` | `list(object)` | `channels` is a root level property defining a list of channels to which checks can send notifications via. |
| `channels[].name` | `string` | This defines the name of the channel and must be unique across all channels. |
| `channels[].type` | `string` | This defines the type of channel which affects how the `apiKey` and `chatId`  properties are consumed. Currently only `"telegram"` and `"slack"` are supported. See notes at the bottom of this section for instructions on what fields to define. |
| `channels[].apiKey` | `ChannelValue` | This defines the API key to use for this channel where applicable. |
| `channels[].apiKey.fromEnv` | `string` | This defines the environment variable from which to retrieve the value of the API key. |
| `channels[].apiKey.value` | `string` | This defines the literal value of the API key. Takes precedence over the `.fromEnv` property. |
| `channels[].chatId` | `ChannelValue` | This defines the chat ID to use for this channel where applicable. |
| `channels[].chatId.fromEnv` | `string` | This defines the environment variable from which to retrieve the value of the chat ID. |
| `channels[].chatId.value` | `string` | This defines the literal value of the chat ID. Takes precedence over the `.fromEnv` property. |
| `channels[].url` | `ChannelValue` | This defines the URL to use to send notifications to where applicable. |
| `channels[].url.fromEnv` | `string` | This defines the environment variable from which to retrieve the value of the URL. |
| `channels[].url.value` | `string` | This defines the literal value of the URL. Takes precedence over the `.fromEnv` property. |

#### Channel type configuration

When `"telegram"` is used for `channels[].type`:
- Set the `apiKey` to the Bot Token by @BotFather
- Set the `chatId` to the ID of the chat which notifications should be sent to

When `"slack"` is used:
- Set the `url` property to the webhook URL

## Debug mode

### Telegram chat ID retrieval

To use a Telegram channel for alerting:

1. Create a new bot with [@BotFather](https://t.me/BotFather) and receive a Telegram bot token (`${BOT_TOKEN}` from here). Set this in your terminal by running `export BOT_TOKEN=${BOT_TOKEN}` or define a `.envrc` file locally. For container deployments, specify these in the deployment manifest in the recommended secure way of your platform.
2. Start `healthcheckr` in Telegram debug mode while specifying the `${BOT_TOKEN}`:
    ```sh
    healthcheckr debug telegram --bot-token ${BOT_TOKEN};
    ```
3. Add the Telegram bot to a chat
4. Use `/info` to trigger a response that will indicate the chat ID (`${CHAT_ID}` from here)
5. Start `healthcheckr` specifying the Telegram chat ID and Telegram bot token as one of the channels in the configuration file
    ```yaml
    # ... other properties ...
    channels:
      # ... other channels ...
      - name: telegram-01
        type: telegram
        apiKey:
          fromEnv: BOT_TOKEN
        chatId: ${CHAT_ID}
    # ... other properties ...
    ```
6. Use the channel by specifying its name in the check configuration:
    ```yaml
    # ... other properties ...
    http:
      - scheme: https
        hostname: google.com
        # ... other properties ...
        channels:
          - telegram-01
    # ... other properties ...
    ```

# Development

## Getting started

Run the following to do a smoke test on Google:

```sh
go run . verify http;
```

To test the worker mode, you'll need to create a Telegram bot and a Slack webhook.

Then create a `.envrc` file and define the following:

```sh
export TELEGRAM_BOT_TOKEN=xxx
export TELEGRAM_CHAT_ID=xxx
export SLACK_WEBHOOK_URL=xxx
```

And then run:

```sh
go run . start worker -c ./examples/config.yaml;
```


## Executing via `go`

To execute locally without compiling, replace all `healthcheckr` invocations with `go run .`
