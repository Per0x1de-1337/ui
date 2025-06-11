#!/bin/bash

# Output file
output_file="All_functions.json"

# Start writing to the output file
echo "[" > "$output_file"

first=1

# Traverse all Go files in the current directory and subdirectories
find . -type f -name "*.go" | while read -r file; do
    # Extract function names using grep and regex
    while IFS= read -r line; do
        # Extract the function name (with or without receiver)
        func_name=$(echo "$line" | sed -nE 's/^func\s+(\([^)]+\)\s*)?([A-Za-z_][A-Za-z0-9_]*)\s*\(.*$/\2/p')
        
        if [[ -n "$func_name" ]]; then
            # Add comma if not the first entry
            if [[ $first -eq 0 ]]; then
                echo "," >> "$output_file"
            fi
            first=0
            echo -n "  {\"function\": \"$func_name\", \"file\": \"${file#./}\"}" >> "$output_file"
        fi
    done < "$file"
done

# Close the JSON array
echo "" >> "$output_file"
echo "]" >> "$output_file"

echo "Output saved to $output_file"
