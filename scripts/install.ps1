$ErrorActionPreference = "Stop"

$Repo = "ghulammuzz-mit/mit-platform"
$Binary = "envctl"
$Filename = "envctl-windows-amd64.exe"
$Url = "https://github.com/$Repo/releases/latest/download/$Filename"
$InstallDir = "$env:LOCALAPPDATA\Programs\envctl"
$Dest = "$InstallDir\envctl.exe"

# Create install dir
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

Write-Host "Downloading $Filename..."
Invoke-WebRequest -Uri $Url -OutFile $Dest

# Add to PATH for current user if not already there
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
    Write-Host "Added $InstallDir to PATH (restart terminal to take effect)"
}

Write-Host "envctl installed to $Dest"
& $Dest --help
