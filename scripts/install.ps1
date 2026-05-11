$ErrorActionPreference = "Stop"

$Repo = "ghulammuzz-mit/mit-pl"
$Binary = "envctl"
$Filename = "envctl-windows-amd64.exe"
$Url = "https://github.com/$Repo/releases/latest/download/$Filename"
$InstallDir = "$env:LOCALAPPDATA\Programs\envctl"
$Dest = "$InstallDir\envctl.exe"

# Create install dir
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

Write-Host "Downloading $Filename..."
Invoke-WebRequest -Uri $Url -OutFile $Dest

# Add to PATH (persistent + current session)
$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    $NewPath = "$UserPath;$InstallDir"

    # Persist to registry (future sessions)
    [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")

    # Update current session immediately
    $env:PATH = "$env:PATH;$InstallDir"

    # Broadcast WM_SETTINGCHANGE so Explorer and other processes pick up the change
    $signature = @'
[DllImport("user32.dll", SetLastError=true, CharSet=CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
    IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
    uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
'@
    $type = Add-Type -MemberDefinition $signature -Name NativeMethods -Namespace Win32 -PassThru -ErrorAction SilentlyContinue
    if ($type) {
        $result = [UIntPtr]::Zero
        $type::SendMessageTimeout([IntPtr]0xffff, 0x001A, [UIntPtr]::Zero, "Environment", 2, 5000, [ref]$result) | Out-Null
    }

    Write-Host "Added $InstallDir to PATH"
}

Write-Host "envctl installed to $Dest"
& $Dest --help
