# avchecker-cli

Command line tool to check a website availability.

```man
usage: avchecker-cli [<flags>] <URL>

Flags:
  --help               Show help.
  -r, --reporter=stdout
                       Reporter to use to publish stats (http|redis)
  -u, --report-url=REPORT-URL
                       URL to report. HTTP(s) when using HTTP or Redis dial info
                       for redis.
  -q, --queue-name=QUEUE-NAME
                       Name of queue to use when using Redis
  --format=json        Format to POST stats (only json available for now)
  --log-level=info     Set log level.
  -o, --log-file=LOG-FILE
                       File to log the output
  -p, --print          Print the output to stdout
  --slack-url=SLACK-URL
                       URL to send message to Slack when server is not
                       availabile
  --check-interval=10  Interval in seconds between availability check
  --publish-interval=60
                       Interval in seconds between stats publication
  -f, --fatal-threshold=0.8
                       Sucess ratio under which a fatal error should be logged
  -x, --extra-fields=EXTRA-FIELDS
                       Extra fields to send with the report stats data
  -H, --request-headers=REQUEST-HEADERS
                       Extra headers to send with the request
  --request-type="GET"
                       The method to use when sending message to the server to
                       check
  --request-body=REQUEST-BODY
                       The body to send to the server to check if using POST
  --version            Show application version.

Args:
  <URL>  URL address to check.
```


Here is a realword example, with headers added to the request, data saved in log file, and alert to Slack on low availability.

```sh
$ ./avchecker-cli -o avchecker.log --slack-url=https://hooks.slack.com/services/WEBHOOK/SECURE/TOKEN -H My-Auth-Token=VERY_SECURE_AUTH_TOKEN http://example.com >> avchecker-stats.log 2>&1
```

If you want to run this in an SSH session, and want to keep it running after exiting, you can simply use `nohup` and '&':

```sh
$ nohup ./avchecker-cli -o avchecker.log --slack-url=https://hooks.slack.com/services/WEBHOOK/SECURE/TOKEN -H My-Auth-Token=VERY_SECURE_AUTH_TOKEN http://example.com >> avchecker-stats.log 2>&1 &
```
