# YACS (go-mc)
YACS uses [go-mc](https://github.com/Tnze/go-mc) and [masscan](https://github.com/zan8in/masscan) libraries to scan the internet for minecraft servers, saving the results to the MongoDB database.

# Stuff YACS collects:
- Server IP
- MOTD
- Server Software
- Players
- Favicon
Yes, it's not that much, but it was supposed to be used only to collect data rescanner could work with.

# Setup
Ubuntu (sudo):
```
apt install masscan golang
git clone https://github.dev/Noroda/yacs-gomc.git
nano main.go # see lines 171 and 160
cd yacs-gomc
go build .
```

# Usage
Figure it out yourself.
