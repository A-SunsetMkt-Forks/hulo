Write-Output "Current time: $(Get-Date)"
Write-Output "10 days from now: $((Get-Date).AddDays(10))"
Write-Output "3 years ago: $((Get-Date).AddYears(-3))"
