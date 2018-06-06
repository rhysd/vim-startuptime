$zip = $Env:APPVEYOR_BUILD_FOLDER + '\vim.zip'
$vim = $Env:APPVEYOR_BUILD_FOLDER + '\vim\'
$redirect = Invoke-WebRequest -URI 'http://vim-jp.org/redirects/koron/vim-kaoriya/latest/win64/'
(New-Object Net.WebClient).DownloadFile($redirect.Links[0].href, $zip)
[Reflection.Assembly]::LoadWithPartialName('System.IO.Compression.FileSystem') > $null
[System.IO.Compression.ZipFile]::ExtractToDirectory($zip, $vim)
$Env:Path = $vim + (Get-ChildItem $vim).Name + ';' + $Env:Path
