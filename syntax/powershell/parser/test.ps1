# 测试 PowerShell 脚本
function Get-Hello {
    param([string]$Name)
    Write-Host "Hello, $Name!"
}

$name = "World"
$age = 25

if ($age -gt 18) {
    Write-Host "Adult"
} else {
    Write-Host "Minor"
}

for ($i = 1; $i -le 5; $i++) {
    Write-Host "Count: $i"
}

Get-Hello -Name $name
