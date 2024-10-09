```lsl --help
NAME:
   lsl - windows replica of linux's 'ls -l'

USAGE:
    [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --path value   root path to start from (default: ".")
   --depth value  search depth (default: 1)
   --verbose      use verbose mode
   --help, -h     show help
```

Example output:
```
$ lsl --path mydir/myprojects --depth 3
bar: 23.4 MB
  .git: 101.1 kB
    hooks: 22.9 kB
    info: 240.0 B
    logs: 3.4 kB
    objects: 72.7 kB
    refs: 153.0 B
  __pycache__: 4.0 kB
  flask_session: 3.2 kB
  static: 35.3 kB
  templates: 1.1 kB
foo: 3.5 MB
  .idea: 4.2 kB
    inspectionProfiles: 174.0 B
  .venv: 23.2 MB
    Include: 0.0 B
    Lib: 21.6 MB
    Scripts: 1.7 MB
  data: 3.5 MB
    parsed: 25.2 kB
    raw: 1.6 MB
    verified: 24.1 kB
```
