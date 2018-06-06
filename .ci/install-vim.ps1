$zip = $Env:APPVEYOR_BUILD_FOLDER + '\vim.zip'
$vim = $Env:APPVEYOR_BUILD_FOLDER + '\vim\'
$redirect = Invoke-WebRequest -URI 'http://vim-jp.org/redirects/koron/vim-kaoriya/latest/win64/'
(New-Object Net.WebClient).DownloadFile($redirect.Links[0].href, $zip)
[Reflection.Assembly]::LoadWithPartialName('System.IO.Compression.FileSystem') > $null
[System.IO.Compression.ZipFile]::ExtractToDirectory($zip, $vim)
$Env:Path = $vim + (Get-ChildItem $vim).Name + ';' + $Env:Path

$zip = $Env:APPVEYOR_BUILD_FOLDER + '\nvim-win64.zip'
$nvim = $Env:APPVEYOR_BUILD_FOLDER + '\nvim-win64\'
$url = 'https://github.com/neovim/neovim/releases/download/nightly/nvim-win64.zip'
(New-Object Net.WebClient).DownloadFile($url, $zip)
[Reflection.Assembly]::LoadWithPartialName('System.IO.Compression.FileSystem') > $null
[System.IO.Compression.ZipFile]::ExtractToDirectory($zip, $nvim)
$Env:PATH = $Env:PATH + ';' + $nvim + '\Neovim\bin\'
