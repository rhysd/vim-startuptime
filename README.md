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

Just run the command with no argument.

```
$ vim-startuptime
```

By default, it tries to run `vim` and `:quit` immediately 10 times, collects the result and show it
to stdout.

```
Extra options: []
Measured: 10 times

Total: 197.416000 msec

   AVERAGE
----------
109.252600: $HOME/.vimrc
 50.508600: opening buffers
 15.871500: /Users/rhysd/.vim/bundle/vim-color-spring-night/colors/spring-night.vim
 10.967500: /Users/rhysd/.vim/bundle/vim-smartinput/autoload/smartinput.vim
  8.736300: /usr/local/Cellar/macvim/HEAD-0db36ff_1/MacVim.app/Contents/Resources/vim/runtime/filetype.vim

...(snip)

  0.002600: clipboard setup
  0.002400: editing files in windows
```

Please see `-help` option to know the command options. If you want to give some options to underlying
`vim` command executions, please specify them after `--` argument in command line as follows:

```
$ vim-startuptime -- --cmd DoSomeCommand
```

## TODO

- Add more metrics like max/min/median
- Temporarily isolate CPU for running Vim if possible

## License

Distributed under [the MIT License](./LICENSE).
