#!/bin/bash

# Arrays
array1=("Item1" "Item2" "Item3")
array2=("New1" "New2" "New3")

# Index where you want to insert array2 into array1 (assuming zero-based indexing)
insert_index=1

# Insert array2 into array1
result_array=("${array1[@]:0:$insert_index}" "${array2[@]}" "${array1[@]:$insert_index}")

# Print the resulting array
for element in "${result_array[@]}"; do
  echo "Element: $element"
done