Goldilocks
==========

Start and stop services based on contents of Bitcoin wallet.

Configuration
-------------

Goldilocks is run as its own user, from a JSON or YAML config file. If you run it as a different user than specified in the config, it will attempt to switch to that user, and die with a message if it fails. This is for security.

In order to allow Goldilocks the ability to control system-level things without full root access, you will want to customize your [/etc/sudoers file](https://help.ubuntu.com/community/Sudoers). You may want to wrap Goldilocks' behavior in wrapper scripts that ignore all arguments, so as to not even provide unfettered access to the specific commands you want to run.

```bash
#!/bin/bash
# /bin/goldilocks_start

service nginx start
```

```bash
#!/bin/bash
# /bin/goldilocks_stop

service nginx stop
```

```
# /etc/sudoers

User_Alias GOLD_NG = goldilocks_nginx
Cmnd_Alias GOLD_NG_CMNDS = /bin/goldilocks_start, /bin/goldilocks_stop

GOLD_NG ALL = (root) NOPASSWD: GOLD_NG_CMNDS 
```

Goldilocks itself is capabable of detecting the current balance in a Bitcoin wallet, and regularly transferring money to a different wallet, via a connection to a local bitcoind process over standard RPC.

```json
{
  "rpc_alias" : {
    "default" : "https://alladin:opensesame@localhost/"
  },
  "services" : [
    {
      "name" : "nginx",
      "description" : "Turn nginx on and off",
      "address" : "<some long bitcoin addr>",
      "threshold" : "0 BTC",
      "commands" : {
        "start": "sudo /bin/goldilocks_start",
        "stop": "sudo /bin/goldilocks_stop",
        "status": "pgrep nginx"
      }
    }
  ],
  "schedule" : [
    {
      "from" : "<same bt addr as earlier>",
      "to" : "<personal bt addr>",
      "amount" : "0.002 BTC",
      "frequency" : "0 5 * * *"
    }
  ],
  "template" : [
    {
      "name" : "overview",
      "source" : "/srv/goldilocks/templates/overview",
      "output" : "/srv/www/gl/index.html"
    },
    {
      "name" : "global json dump",
      "source" : "core.json",
      "output" : "/srv/www/gl/core.json"
    }
  ]
}
```

### Services

Service entries use bitcoin RPC to verify that the address has more money than the threshold, starting or stopping the service as necessary based on the commands specified.

### Schedules

Schedule entries regularly move money from one address to another using RPC, using the cron schedule syntax.

### Templates

Template entries use the [Go template module](http://golang.org/pkg/html/template/) to produce HTML static files, that are updated whenever Goldilocks internal state updates. You can either serve these directly by saving them to your web server's static file tree, which allows for very efficient/cheap availability.

For more "involved" uses, you can also save a dump of all (non-security-risk) internal state as a JSON, YAML, or bencode file, using special sources "core.json", "core.yaml", and "core.bencode", respectively.

Your OS's filesystem caching should be perfectly sufficient to ensure that this information is cached in-memory for your application, without requiring a disk hit.
