Dim today, endDate, diff
today = Date
endDate = DateSerial(2025, 12, 31)
diff = DateDiff("d", today, endDate)
WScript.Echo "There are " & diff & " days left until the end of 2025!"
