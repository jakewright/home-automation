# Provision a Raspberry Pi to be a home automation host

## Setting up Raspbian

These instructions will assume you're using Mac OS and that the Raspberry Pi has a static IP address of 192.168.1.210.

### SSH access

Add your SSH key to the Raspberry Pi

```
brew install ssh-copy-id
ssh-copy-id pi@192.168.1.210
```

Set up agent forwarding by adding the following to `~/.ssh/config`

```
Host 192.168.1.*
  ForwardAgent yes
```


### Syslog
Add the following config to `/etc/rsyslog.conf` to override the default log format and use the RFC-5424 format instead.

```
# Use RFC-5424 format
$ActionFileDefaultTemplate RSYSLOG_SyslogProtocol23Format
```

Place `home-automation.conf` at `/etc/rsyslog.d/home-automation.conf`.
This will forward all logs to the logstash instance.



### Installing software

```
sudo apt-get update
```

Install git on the Raspberry Pi
```
sudo apt-get install git
```

Install Python 3
```
sudo apt-get install python3
```

Install pip
```
sudo apt-get install python3-pip
```

Or pip for python2
```
sudo apt-get install python-pip
```

Install virtualenv
```
sudo pip3 install virtualenv
```

Install node
```
# Download
wget https://nodejs.org/download/release/v8.1.0/node-v8.1.0-linux-armv6l.tar.gz

# Unzip
tar -xvf node-v8.1.0-linux-armv6l.tar.gz

# Test that it works
node-v8.1.0-linux-armv6l/bin/node -v

# Copy to /usr/local to be able to globally run node
cd node-v8.1.0-linux-armv6l/
sudo cp -R bin/ /usr/local/
sudo cp -R include/ /usr/local/
sudo cp -R lib/ /usr/local/
sudo cp -R share/ /usr/local/
```

Note that the `cp` commands will not overwrite the existing directories. The contents of the `bin`, `include`, `lib`, and
`share` folders will be copied *into* the existing directories in `/usr/local/`.

Now `node -v` should output `v8.1.0`.

### Clone the repository

```
cd /home/pi/
mkdir home-automation
```


### Docker

Docker was tested on the Raspberry Pi but it lacked enough memory to run a Node app.
A Python service would run successfully but response times were twice as long as running it locally.

For reference, this is how to install Docker on a Raspberry Pi.
(Information from https://blog.alexellis.io/getting-started-with-docker-on-raspberry-pi/)

```
curl -sSL https://get.docker.com | sh
```

Set Docker to auto-start
```
sudo systemctl enable docker
```

Reboot the pi or start docker with
```
sudo systemctl start docker
```

Add the pi user to the docker group
```
sudo usermod -aG docker pi
```
