param(
    [string]$Version = "v0.2.0",
    [string]$InstallPath = "$env:USERPROFILE\.local\bin"
)

$ErrorActionPreference = "Stop"

$InfoColor = "Cyan"
$SuccessColor = "Green"
$ErrorColor = "Red"
$WarningColor = "Yellow"

$BaseUrl = "https://github.com/hulo-lang/hulo/releases/download/$Version"

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

function Write-Download($msg) {
    Write-Host "[DOWNLOAD] $msg" -ForegroundColor Cyan
}

function Write-Install($msg) {
    Write-Host "[INSTALL] $msg" -ForegroundColor Green
}

function Write-Verify($msg) {
    Write-Host "[VERIFY] $msg" -ForegroundColor Yellow
}

function Write-Header($msg) {
    Write-Host ""
    Write-Host "=== $msg ===" -ForegroundColor White
    Write-Host ""
}

function Get-OperatingSystem {
    if ($IsWindows -or $env:OS -eq "Windows_NT") {
        return "Windows"
    }
    elseif ($IsLinux) {
        return "Linux"
    }
    elseif ($IsMacOS) {
        return "Darwin"
    }
    else {
        Write-ErrorAndExit "Unsupported operating system"
    }
}

function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    if ($arch -eq "AMD64") {
        return "x86_64"
    }
    elseif ($arch -eq "ARM64") {
        return "arm64"
    }
    elseif ($arch -eq "x86") {
        return "i386"
    }
    else {
        Write-ErrorAndExit "Unsupported architecture: $arch"
    }
}

function Get-FileExtension($os) {
    if ($os -eq "Windows") {
        return "zip"
    }
    else {
        return "tar.gz"
    }
}

function Download-File($url, $outputPath) {
    Write-Download "Downloading from: $url"

    try {
        Invoke-WebRequest -Uri $url -OutFile $outputPath -UseBasicParsing
        Write-Success "Download completed successfully"
    }
    catch {
        Write-ErrorAndExit "Failed to download file: $($_.Exception.Message)"
    }
}

function Verify-Checksum($filePath, $expectedChecksum) {
    Write-Verify "Verifying SHA256 checksum..."

    try {
        $actualChecksum = (Get-FileHash -Path $filePath -Algorithm SHA256).Hash.ToLower()

        if ($actualChecksum -eq $expectedChecksum) {
            Write-Success "Checksum verification passed"
        }
        else {
            Write-ErrorAndExit "Checksum verification failed. Expected: $expectedChecksum, Got: $actualChecksum"
        }
    }
    catch {
        Write-Warning "Checksum verification skipped (Get-FileHash not available)"
        exit 1
    }
}

function Extract-File($filePath, $extractDir) {
    Write-Step "Extracting archive: $filePath"

    try {
        if ($filePath -match "\.zip$") {
            Expand-Archive -Path $filePath -DestinationPath $extractDir -Force
        }
        elseif ($filePath -match "\.tar\.gz$") {
            # PowerShell 7+ 支持 tar 命令
            if (Get-Command tar -ErrorAction SilentlyContinue) {
                tar -xzf $filePath -C $extractDir
            }
            else {
                Write-ErrorAndExit "tar command not found. Please install tar or use PowerShell 7+"
            }
        }
        else {
            Write-ErrorAndExit "Unknown file format: $filePath"
        }
        Write-Success "Archive extracted successfully"
    }
    catch {
        Write-ErrorAndExit "Failed to extract file: $($_.Exception.Message)"
    }
}

function Install-Binary($binaryFiles, $installDir) {
    Write-Install "Installing hulo binaries to: $installDir"

    try {
        # 创建安装目录
        if (-not (Test-Path $installDir)) {
            New-Item -ItemType Directory -Path $installDir -Force | Out-Null
            Write-Info "Created installation directory: $installDir"
        }

        foreach ($file in $binaryFiles) {
            $installPath = Join-Path $installDir $file.Name
            Copy-Item -Path $file.FullName -Destination $installPath -Force
            Write-Success "Binary installed successfully: $installPath"
        }

        # 添加到 PATH 环境变量
        $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
        if ($currentPath -notlike "*$installDir*") {
            $newPath = "$currentPath;$installDir"
            [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
            Write-Info "Added to PATH: $installDir"
        }
        else {
            Write-Info "PATH already contains: $installDir"
        }

    }
    catch {
        Write-ErrorAndExit "Failed to install binaries: $($_.Exception.Message)"
    }
}

function Install-StdLib($extractDir, $installDir) {
    Write-Install "Installing standard library and modules"

    try {
        # 创建 HULO_MODULES 目录（在用户目录下）
        $huloModulesDir = Join-Path $env:USERPROFILE "HULO_MODULES"
        if (-not (Test-Path $huloModulesDir)) {
            New-Item -ItemType Directory -Path $huloModulesDir -Force | Out-Null
            Write-Info "Created modules directory: $huloModulesDir"
        }

        # 复制除了可执行文件和压缩包之外的所有文件到 HULO_MODULES 目录
        Write-Step "Copying standard library files to: $huloModulesDir"
        $fileCount = 0
        Get-ChildItem -Path $extractDir -Exclude "hulo.exe", "hulo", "*.zip", "*.tar.gz" | ForEach-Object {
            if ($_.PSIsContainer) {
                # 如果是目录，复制整个目录
                Copy-Item -Path $_.FullName -Destination $huloModulesDir -Recurse -Force
                $fileCount++
            }
            else {
                # 如果是文件，复制文件
                Copy-Item -Path $_.FullName -Destination $huloModulesDir -Force
                $fileCount++
            }
        }
        Write-Success "Standard library installed successfully ($fileCount items)"

        # 设置 HULO_PATH 环境变量指向 HULO_MODULES 目录
        Write-Step "Configuring HULO_PATH environment variable"

        # 直接设置为用户环境变量（不需要管理员权限）
        $result = cmd /c "setx HULO_PATH `"$huloModulesDir`"" 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Info "HULO_PATH set to: $huloModulesDir"
        }
        else {
            # 如果 setx 失败，使用 PowerShell 方法
            [Environment]::SetEnvironmentVariable("HULO_PATH", $huloModulesDir, "User")
            Write-Info "HULO_PATH set to: $huloModulesDir (fallback method)"
        }

        # 设置当前会话的环境变量
        $env:HULO_PATH = $huloModulesDir

    }
    catch {
        Write-ErrorAndExit "Failed to install standard library: $($_.Exception.Message)"
    }
}

function Main {
    Write-Host "Hulo Installer $Version"
    Write-Host "================================"

    $os = Get-OperatingSystem
    $arch = Get-Architecture
    $extension = Get-FileExtension $os

    $filename = "hulo_${os}_${arch}.${extension}"
    $downloadUrl = "$BaseUrl/$filename"

    Write-Info "Downloading from: $downloadUrl"
    $tempDir = Join-Path $env:TEMP "hulo-install-$(Get-Random)"
    $downloadPath = Join-Path $tempDir $filename

    # 动态获取 checksums.txt 文件
    $checksumsUrl = "$BaseUrl/checksums.txt"
    Write-Info "Downloading checksums from: $checksumsUrl"

    try {
        $response = Invoke-WebRequest -Uri $checksumsUrl -UseBasicParsing
        $checksumsContent = [System.Text.Encoding]::UTF8.GetString($response.Content)
        $checksums = @{}

        # Debug: Show the first 200 characters of the content
        Write-Info "Checksums content preview: $($checksumsContent.Substring(0, [Math]::Min(200, $checksumsContent.Length)))"

        # 解析 checksums.txt 文件
        $checksumsContent -split "`n" | ForEach-Object {
            $line = $_.Trim()
            if ($line -and -not $line.StartsWith("#")) {
                $parts = $line -split '\s+'
                if ($parts.Length -ge 2) {
                    $checksum = $parts[0]
                    $filename = $parts[1]
                    $checksums[$filename] = $checksum
                    Write-Info "Found checksum: $checksum for file: $filename"
                }
            }
        }

        Write-Info "Successfully loaded checksums for $($checksums.Count) files"
        Write-Info "Available files: $($checksums.Keys -join ', ')"
    }
    catch {
        Write-ErrorAndExit "Failed to download checksums: $($_.Exception.Message)"
    }

    $expectedChecksum = $checksums[$filename]
    if (-not $expectedChecksum) {
        Write-ErrorAndExit "No checksum found for $filename in checksums.txt"
    }

    try {
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

        Download-File $downloadUrl $downloadPath

        Verify-Checksum $downloadPath $expectedChecksum

        Extract-File $downloadPath $tempDir

        $binDir = Join-Path $tempDir "bin"
        $exeFiles = @()
        if (Test-Path $binDir) {
            $exeFiles = Get-ChildItem -Path $binDir -Filter "*.exe" -File
            if ($exeFiles.Count -eq 0) {
                Write-ErrorAndExit "No .exe files found in bin directory."
            }
        } else {
            Write-ErrorAndExit "bin directory not found in extracted archive."
        }

        Install-Binary $exeFiles $InstallPath

        Write-Info "Binary installation completed, proceeding to install standard library..."
        Write-Info "About to call Install-StdLib with tempDir: $tempDir, InstallPath: $InstallPath"
        Install-StdLib $tempDir $InstallPath
        Write-Info "Install-StdLib completed successfully"

        Write-Host ""
        Write-Success "Installation completed successfully!"
        Write-Info "You can now use 'hulo' command"
        # 获取环境变量值（优先检查用户，然后检查当前会话）
        $HULO_PATH = [Environment]::GetEnvironmentVariable('HULO_PATH', 'User')
        if (-not $HULO_PATH) {
            $HULO_PATH = $env:HULO_PATH
        }

        Write-Info "Standard library and modules are available at: $HULO_PATH"
        Write-Info "HULO_PATH environment variable: $HULO_PATH"
    }
    finally {
        if (Test-Path $tempDir) {
            Remove-Item -Path $tempDir -Recurse -Force
        }
    }
}

Main
