$today = Get-Date
$end = Get-Date "2025-12-31"
$diff = ($end - $today).Days
Write-Output "There are $diff days left until the end of 2025!"
