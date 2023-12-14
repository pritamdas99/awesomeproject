#!/bin/bash

#function process_array {
#  local my_array=("$@")  # Use "$@" to reference all arguments passed to the function
#  for element in "${my_array[@]}"; do
#    echo "Element: $element"
#  done
#}
#
## Declare an array
#original_array=("Item1" "Item2" "Item3")
#
## Call the function and pass the array as an argument
#process_array "${original_array[@]}"

function string_to_array {
  local input_string=$1
  local IFS='<>/=" '
  read -ra items <<< "$input_string"

  result=()

  # Iterate over the original array
  for item in "${items[@]}"; do
    # Check if the item is not an empty string
    if [[ -n "$item" ]]; then
      # Add non-empty item to the new array
      result+=("$item")
    fi
  done

  echo "${result[@]}"
}

function get_name {
  local items=("$@")

  arr_length=${#items[@]}
  #echo "arrrrrrrrrreln${arr_length}"

  name="name"
  res=""
  for ((i=0; i<arr_length; i++)); do
      # Check if the item is not an empty string

      if [[ "${items[i]}" == "${name}" ]]; then
        # Add non-empty item to the new array
         res=${items[i+1]}
         break
      fi
    done

  echo $res
}




# Your input string with multiple characters as separators
input_string='<str name="host"> ${solr.host:} </str>'
result_array=($(string_to_array "$input_string"))
echo "len ${#result_array[@]}"
for item in "${result_array[@]}"; do
  echo "Item: $item"
done

echo "res array ${result_array[@]} ${#result_array[@]}"

arr_length=${#result_array[@]}
echo $arr_length
srt=($(get_name "${result_array[@]}"))
echo "srts $srt"
echo done