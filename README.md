Better `vim --startuptime`
==========================
[![CI Badge][]][CI]
[![Codecov Status][]][Codecov]

`vim-startuptime` is a small Go program to measure startup time of Vim or Neovim. This program aims
to be an alternative of `--startuptime` option of Vim, which measures the startup time metrics to
allow vimmers to optimize Vim's startup.

After warm-up, `vim-startuptime` runs `vim --startuptime` multiple times internally and collects the
metrics from the results (e.g. average time for loading each plugin's scripts). The summary calculated
from the measurements is output to stdout.

Tested on Linux, Mac and Windows with both Vim and Neovim.



## Installation

Download an executable from [a release page](https://github.com/rhysd/vim-startuptime/releases).

If you want to install the latest stable version, please run following command. Go toolchain is
necessary for running the command.

```sh
go install github.com/rhysd/vim-startuptime@latest
```



## Requirements

- `vim` 7.4.1444 or later (for `--not-a-term` startup option)
- `nvim`



## Usage

Just run the command with no argument.

```sh
vim-startuptime
```

By default, it tries to run `vim` and `:quit` immediately 10 times, collects the results and outputs
a summary of them to stdout.

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

If you want to run with Neovim, please specify `-vimpath` option as follows:

```
$ vim-startuptime -vimpath nvim
```

If you want to give some options or arguments to underlying `vim` command executions, please specify
them after `--` argument in command line as follows:

```
$ vim-startuptime -- --cmd DoSomeCommand path/to/some/file
```

Please see `-help` option to know the command options.



## What's Next after running `vim-startuptime`?

By running `vim-startuptime`, now you know which script file takes time to run. What you should do
next is `:profile`.

```sh
vim --cmd 'profile start profile.log' --cmd 'profile! file /path/to/slow_script.vim' -c quit
```

Profiled results are dumped to `profile.log`. Please check it. In log file, `:set ft=vim` would help
you analyze the results.
Please see `:help profile` for more details.



## (Maybe) TODO

- Add more metrics like median
- Temporarily isolate CPU for running Vim if possible



## License

Distributed under [the MIT License](./LICENSE).



[CI Badge]: https://github.com/rhysd/vim-startuptime/actions/workflows/ci.yml/badge.svg
[CI]: https://github.com/rhysd/vim-startuptime/actions/workflows/ci.yml
[Codecov Status]: https://codecov.io/gh/rhysd/vim-startuptime/branch/master/graph/badge.svg
[Codecov]: https://codecov.io/gh/rhysd/vim-startuptime
