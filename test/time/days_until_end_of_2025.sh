#!/bin/bash
today=$(date +%s)
end_of_year=$(date -d "2025-12-31" +%s)
diff_days=$(((end_of_year - today) / 86400))
echo "There are $diff_days days left until the end of 2025!"
