Better `vim --startuptime`
==========================
[![TravisCI Status][]][TravisCI]
[![AppVeyor Status][]][AppVeyor]
[![Codecov Status][]][Codecov]

`vim-startuptime` is a small Go program to measure startup time of Vim. This program aims to be an
alternative of `--startuptime` option of Vim, which measures the startup time metrics to allow vimmers
to optimize Vim's startup.

`vim-startuptime` runs `vim --startuptime` multiple times internally and collects the metrics from
the results (e.g. average time for loading each plugin's scripts).

## Installation

```
$ go get github.com/rhysd/vim-startuptime
```

## Requirements

- Vim 7.4.1444 or later (for `--not-a-term` startup option)

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

Total Average: 189.954400 msec
Total Max:     198.062000 msec
Total Min:     183.966000 msec

  AVERAGE        MAX       MIN
-------------------------------
98.532900 102.605000 94.275000: $HOME/.vimrc
51.859600  56.937000 49.897000: opening buffers
17.027900  18.810000 16.277000: /Users/rhysd/.vim/bundle/vim-color-spring-night/colors/spring-night.vim
11.878900  13.153000 10.567000: /Users/rhysd/.vim/bundle/vim-smartinput/autoload/smartinput.vim
 9.407600  11.710000  8.606000: /usr/local/Cellar/macvim/HEAD-0db36ff_1/MacVim.app/Contents/Resources/vim/runtime/filetype.vim

...(snip)

 0.009100   0.012000  0.007000: window checked
 0.009000   0.012000  0.008000: inits 3
 0.003000   0.005000  0.002000: clipboard setup
 0.002600   0.004000  0.002000: editing files in windows
```

Please see `-help` option to know the command options. If you want to give some options to underlying
`vim` command executions, please specify them after `--` argument in command line as follows:

```
$ vim-startuptime -- --cmd DoSomeCommand
```

## TODO

- Add more metrics like median
- Temporarily isolate CPU for running Vim if possible

## License

Distributed under [the MIT License](./LICENSE).



[TravisCI Status]: https://travis-ci.org/rhysd/vim-startuptime.svg?branch=master
[TravisCI]: https://travis-ci.org/rhysd/vim-startuptime
[AppVeyor Status]: https://ci.appveyor.com/api/projects/status/1tpyd9q9tw3ime5u/branch/master?svg=true
[AppVeyor]: https://ci.appveyor.com/project/rhysd/vim-startuptime/branch/master
[Codecov Status]: https://codecov.io/gh/rhysd/vim-startuptime/branch/master/graph/badge.svg
[Codecov]: https://codecov.io/gh/rhysd/vim-startuptime
