#!/bin/bash

# Example array
my_array=("Item1" "Item2" "Item3" "Item4" "Item5")

# Starting index for the slice
start_index=1

# Ending index for the slice
end_index=3

# Slice the array
sliced_array=("${my_array[@]:$start_index:$((end_index - start_index + 1))}")

# Print the sliced array
for item in "${sliced_array[@]}"; do
  echo "Item: $item"
done
