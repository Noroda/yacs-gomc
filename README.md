# Crappiest Minecraft Server Scanner!
You **don't** want to use this, here's why:
- sometimes it doesn't wants to display motd
- sometimes scanner just skips IP address
- color codes and formatting codes are visible
- omg I actually fixed this

# ***Use something better instead of this!***
But if you decided to use this, then install masscan and run ``sudo go run main.go``

flags:
```
--range <ip-range> | IP range to scan | default is "127.0.0.1"
--port-range <port-range> | Port range to scan | default is "25565"
--output <file> | Change name of output file | default is "output.txt"
--rate <rate> | set masscan rate so it won't take so much time | default is "1000"
```
