# Provision a Raspberry Pi to be a home automation host


Add the following config to `/etc/rsyslog.conf` to override the default log format and use the RFC-5424 format instead.

```
# Use RFC-5424 format
$ActionFileDefaultTemplate RSYSLOG_SyslogProtocol23Format
```

Place `home-automation.conf` at `/etc/rsyslog.d/home-automation.conf`.
This will forward all logs to the logstash instance.