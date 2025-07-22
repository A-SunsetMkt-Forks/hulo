param(
    [string]$InstallPath = "$env:USERPROFILE\.local\bin"
)

$ErrorActionPreference = "Stop"

$InfoColor = "Cyan"
$SuccessColor = "Green"
$ErrorColor = "Red"
$WarningColor = "Yellow"

function Write-Info($msg) {
    Write-Host "[INFO] $msg" -ForegroundColor $InfoColor
}

function Write-Success($msg) {
    Write-Host "[SUCCESS] $msg" -ForegroundColor $SuccessColor
}

function Write-ErrorAndExit($msg) {
    Write-Host "[ERROR] $msg" -ForegroundColor $ErrorColor
    exit 1
}

function Write-Warning($msg) {
    Write-Host "[WARNING] $msg" -ForegroundColor $WarningColor
}

function Write-Step($msg) {
    Write-Host "[STEP] $msg" -ForegroundColor Magenta
}

function Write-Remove($msg) {
    Write-Host "[REMOVE] $msg" -ForegroundColor Red
}

function Write-Header($msg) {
    Write-Host ""
    Write-Host "=== $msg ===" -ForegroundColor White
    Write-Host ""
}

function Remove-Binary($installDir) {
    Write-Remove "Removing hulo binaries from: $installDir"

    $exeNames = @("hulo.exe", "hlpm.exe", "hulo-repl.exe")
    foreach ($exe in $exeNames) {
        $binaryPath = Join-Path $installDir $exe
        if (Test-Path $binaryPath) {
            Remove-Item -Path $binaryPath -Force
            Write-Success "Binary removed successfully: $binaryPath"
        }
        else {
            Write-Warning "Binary not found: $binaryPath"
        }
    }

    # 从 PATH 环境变量中移除安装目录
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -like "*$installDir*") {
        $newPath = ($currentPath -split ';' | Where-Object { $_ -ne $installDir }) -join ';'
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Info "Removed from PATH: $installDir"
    }
    else {
        Write-Info "PATH does not contain: $installDir"
    }
}

function Remove-HuloModules() {
    Write-Remove "Removing HULO_MODULES directory"

    # 获取 HULOPATH 环境变量
    $huloPath = [Environment]::GetEnvironmentVariable("HULO_PATH", "User")
    if (-not $huloPath) {
        $huloPath = $env:HULOPATH
    }

    if ($huloPath -and (Test-Path $huloPath)) {
        Write-Info "Found modules directory: $huloPath"

        # 确认删除
        $confirm = Read-Host "Are you sure you want to delete the HULO_MODULES directory? (y/N)"
        if ($confirm -eq 'y' -or $confirm -eq 'Y') {
            Remove-Item -Path $huloPath -Recurse -Force
            Write-Success "Modules directory removed successfully: $huloPath"
        }
        else {
            Write-Info "Skipped deletion of modules directory"
        }
    }
    else {
        Write-Warning "HULO_MODULES directory not found"
    }
}

function Remove-EnvironmentVariable() {
    Write-Remove "Removing HULO_PATH environment variable"

    # 检查用户环境变量
    $userHuloPath = [Environment]::GetEnvironmentVariable("HULO_PATH", "User")
    if ($userHuloPath) {
        [Environment]::SetEnvironmentVariable("HULO_PATH", $null, "User")
        Write-Info "Removed HULO_PATHv from user environment variables"
    }
    else {
        Write-Info "HULO_PATH not found in user environment variables"
    }

    # 检查全局环境变量
    $machineHuloPath = [Environment]::GetEnvironmentVariable("HULO_PATH", "Machine")
    if ($machineHuloPath) {
        Write-Warning "Found global HULO_PATH environment variable. Administrator privileges may be required."
        $confirm = Read-Host "Do you want to try to remove the global HULO_PATH? (y/N)"
        if ($confirm -eq 'y' -or $confirm -eq 'Y') {
            # 尝试删除全局环境变量
            $globalRemoved = $false
            try {
                [Environment]::SetEnvironmentVariable("HULO_PATH", $null, "Machine")
                $globalRemoved = $true
            }
            catch {
                $globalRemoved = $false
            }

            if ($globalRemoved) {
                Write-Success "Removed HULO_PATH from global environment variables"
            }
            else {
                Write-Warning "Failed to remove global HULO_PATH (requires administrator privileges)"
            }
        }
        else {
            Write-Info "Skipped removal of global HULO_PATH"
        }
    }
    else {
        Write-Info "HULO_PATH not found in global environment variables"
    }
}

# 主函数
Write-Header "Hulo Uninstaller"

try {
    Write-Info "Starting uninstallation process"
    Write-Info "Install path: $InstallPath"

    # 删除可执行文件
    Write-Step "Step 1: Removing binary file"
    Remove-Binary $InstallPath

    # 删除 HULO_MODULES 目录
    Write-Step "Step 2: Removing modules directory"
    Remove-HuloModules

    # 删除环境变量
    Write-Step "Step 3: Removing environment variables"
    Remove-EnvironmentVariable

    Write-Header "Uninstallation Summary"
    Write-Success "Uninstallation completed successfully"
    Write-Info "Hulo has been completely removed from your system"

}
catch {
    Write-ErrorAndExit "Uninstallation failed: $($_.Exception.Message)"
}
