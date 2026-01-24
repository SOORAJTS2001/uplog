> This repository is in active development, hence change in documentation is expected

<h1 style="display: flex; align-items: center; gap: 8px;">
  <img
    src="docs/uplog_icon.png"
    alt="Uplog icon"
    height="36"
  />
  Uplog
</h1>



[![CI](https://github.com/SOORAJTS2001/uplog/actions/workflows/ci.yml/badge.svg)](
https://github.com/SOORAJTS2001/uplog/actions/workflows/ci.yml
)
[![Deployment](https://github.com/SOORAJTS2001/uplog/actions/workflows/backend-deploy.yml/badge.svg)](
https://github.com/SOORAJTS2001/uplog/actions/workflows/backend-deploy.yml
)
[![CLI Build](https://github.com/SOORAJTS2001/uplog/actions/workflows/cli-release.yml/badge.svg)](
https://github.com/SOORAJTS2001/uplog/actions/workflows/cli-release.yml
)
![Python](https://img.shields.io/badge/python-3.12%20|%203.13%20|%203.14-blue)
![License](https://img.shields.io/github/license/SOORAJTS2001/uplog)

A free and open-source log monitoring platform that works in milli-seconds. No signup, no dependencies, and no code
rewrites - just plug in the CLI and watch your logs stream live.

### Installation

- Linux/MacOS

```bash
curl -fsSL https://uplog.live/install.sh | sh
```

- Windows
    - Download executables from [releases](https://github.com/SOORAJTS2001/uplog/releases)

### Breakdown

- CLI
    - It is a go binary, which could be used to monitor any cli logs from program/process
    - These logs are batched and send to the backend server, the shareable url would be shown right in the terminal
    - It doesn't buffer the output, so you could see it on your cli as soon as it comes.

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

Set batch upload size

```bash
uplog --poll <batch_size> python main.py
```

Tag a session

```bash
uplog --tag <tag> python main.py
```

Note: You could tag and batch at the same time

List all recorded sessions
```bash
uplog list
```

Delete all  recorded sessions
```bash
uplog purge
```

Delete a single session
```bash
uplog delete <session_id>
```

These logs are intermediately written to temporary log files, which would be deleted after successful log update
