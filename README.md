## Uplog

[![CI](https://github.com/SOORAJTS2001/uplog/actions/workflows/ci.yml/badge.svg)](
https://github.com/SOORAJTS2001/uplog/actions/workflows/ci.yml
)
[![Deployment](https://github.com/SOORAJTS2001/uplog/actions/workflows/backend-deploy.yml/badge.svg)](
https://github.com/SOORAJTS2001/uplog/actions/workflows/backend-deploy.yml
)
![Python](https://img.shields.io/badge/python-3.12%20|%203.13%20|%203.14-blue)
![License](https://img.shields.io/github/license/SOORAJTS2001/uplog)


A free and open-source log monitoring platform that works in milli-seconds. No signup, no dependencies, and no code rewrites - just plug in the CLI and watch your logs stream live.



### Breakdown
- CLI
    -   It is a go binary, which could be used to monitor any cli logs from program/process
    -   These logs are batched and send to the backend server, the shareable url would be shown right in the terminal
    -   It doesn't buffer the output, so you could see it on your cli as soon as it comes.

### Usage
```python
# main.py
import time
for i in range(10):
    print(i)
    time.sleep(1)
```

```bash
uplog python main.py
```

These logs are intermediately written to temporary log files, which would be deleted after successful log update

> This repository is in active development, hence change in documentation is expected
