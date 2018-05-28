Better `vim --startuptime`
==========================

`vim-startuptime` is a small Go program to measure startup time of Vim. This program aims to be an
alternative of `--startuptime` option of Vim, which measures the startup time metrics to allow vimmers
to optimize Vim's startup.

`vim-startuptime` runs `vim --startuptime` multiple times internally and collects the metrics from
the results (e.g. average time for loading each plugin's scripts).

## Installation

```
$ go get github.com/rhysd/vim-startuptime
```

## Usage

```
$ vim-startuptime
```

Please see `-help` option to know the command options. If you want to give some options to underlying
`vim` command executions, please specify them after `--` argument in command line as follows:

```
$ vim-startuptime -- --cmd DoSomeCommand
```

## TODO

- Add tests
- Add more metrics like max/min/median
- Add option to show spent time of sourced scripts only

## License

Please read [LICENSE](./LICENSE).
