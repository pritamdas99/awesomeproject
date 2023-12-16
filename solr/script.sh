#!/bin/bash

sudo mkdir .
# Specify the file name
target="target.xml"

# Read the entire file content into a variable
target_xml=$(cat "$target")

#remove new lines
target_xml=$(echo -n "$target_xml" | tr -d '\n')

# Print or manipulate the file content as needed
#echo "$input_content"

#replace the pattern >{multiple space}< by ><
target_xml=$(echo "$target_xml" | sed 's/>\s*</></g')

#replace the pattern >< by >\n<
target_xml=$(echo "$target_xml" | sed 's/></>\n</g')

old_ifs=$IFS

# Set IFS to newline
IFS=$'\n'

# Convert the string to an array
target_xml_lines=($target_xml)

  for item in "${target_xml_lines[@]}"; do
    # Check if the item is not an empty string
   echo "${item}"
  done

IFS=$old_ifs

diff="diff.xml"

# Read the entire file content into a variable
diff_xml=$(cat "$diff")

#remove new lines
diff_xml=$(echo -n "$diff_xml" | tr -d '\n')

# Print or manipulate the file content as needed


#replace the pattern >{multiple space}< by ><
diff_xml=$(echo "$diff_xml" | sed 's/>\s*</></g')

#replace the pattern >< by >\n<
diff_xml=$(echo "$diff_xml" | sed 's/></>\n</g')

old_ifs=$IFS

# Set IFS to newline
IFS=$'\n'

# Convert the string to an array
diff_xml_lines=($diff_xml)

i1=-1
i2=-1


#echo ${diff_xml_lines[@]}
#echo starting...

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

  for item in "${result[@]}"; do
      # Check if the item is not an empty string
      echo "$item"
    done

#  echo ${result[@]}
}

function get_name {
  local items=("$@")

  arr_length=${#items[@]}
  name="name"
  res=""
  for ((i=0; i<arr_length; i++)); do
      if [[ "${items[i]}" == "name" ]]; then
        # Add non-empty item to the new array
         res=${items[i+1]}

         break
      fi
    done

  echo $res
}


diff_length=${#diff_xml_lines[@]}
target_length=${#target_xml_lines[@]}

i1=-1
i2=-1

for ((i=0; i<diff_length; i++)); do
  result_array=($(string_to_array "${diff_xml_lines[i]}"))
 # echo "echoing $i ${result_array[0]}"
  diffName="${result_array[0]}"
 # echo "check ${diffName}"
  if [[ "${diffName}" == "shardHandlerFactory" ]]; then
    i1=$i
    #echo "how enetred here if diffname is ${diffName}"
    break
  fi
done
#echo "i1 $i1"

if [[ "$i1" -ne -1 ]]; then
  for((i=i1+1; i<diff_length; i++)); do
    result_array=($(string_to_array "${diff_xml_lines[i]}"))
   # echo "${result_array[0]}"
   # echo "echoing $i ${result_array[0]}"
    diffName="${result_array[0]}"
    if [[ ${diffName} == "shardHandlerFactory" ]]; then
      i2=$i
      break
    fi
  done
  sliced_array=("${diff_xml_lines[@]:$i1:$((i2 - i1 + 1))}")
  diff_xml_lines=("${diff_xml_lines[@]:0:$i1}" "${diff_xml_lines[@]:$((i2+1))}")
  f=0
  for item in ${sliced_array[@]}; do
    result_array=($(string_to_array ${item}))
    ln=${#result_array[@]}
    diffName=($(get_name "${result_array[@]}"))
    if [[ ${diffName} == "connTimeout" ]]; then
      f=1
      break
    fi
  done

  if [[ $f == 0 ]]; then
    element='<int name="connTimeout">${connTimeout:60000}</int>'
    sliced_array=("${sliced_array[@]:0:1}" "${element}" "${sliced_array[@]:1}")
  fi

  f=0
  for item in ${sliced_array[@]}; do
    result_array=($(string_to_array ${item}))
    ln=${#result_array[@]}
    diffName=($(get_name "${result_array[@]}"))
    echo "diffname ${diffName}"
    if [[ ${diffName} == "socketTimeout" ]]; then
      echo "socket timeout present"
      f=1
      break
    fi
  done

  if [[ $f == 0 ]]; then
    element='<int name="socketTimeout">${socketTimeout:600000}</int>'
    sliced_array=("${sliced_array[@]:0:1}" "${element}" "${sliced_array[@]:1}")
  fi


    i1=-1
    i2=-1

    for ((i=0; i<target_length; i++)); do
      result_array=($(string_to_array "${target_xml_lines[i]}"))
      diffName="${result_array[0]}"
      if [[ "${diffName}" == "shardHandlerFactory" ]]; then
        i1=$i
        break
      fi
    done
    if [[ "$i1" -ne -1 ]]; then
      for((i=i1+1; i<target_length; i++)); do
        result_array=($(string_to_array "${diff_xml_lines[i]}"))
        diffName="${result_array[0]}"
        if [[ ${diffName} == "shardHandlerFactory" ]]; then
          i2=$i
          break
        fi
      done
      target_xml_lines=("${target_xml_lines[@]:0:$i1}" "${target_xml_lines[@]:$((i2+1))}")
      target_xml_lines=("${target_xml_lines[@]:0:2}" "${sliced_array[@]}" "${target_xml_lines[@]:2}")

    else
      target_xml_lines=("${target_xml_lines[@]:0:2}" "${sliced_array[@]}" "${target_xml_lines[@]:2}")
    fi
fi

i1=-1
i2=-1


diff_length=${#diff_xml_lines[@]}
target_length=${#target_xml_lines[@]}

for ((i=0; i<diff_length; i++)); do
  result_array=($(string_to_array "${diff_xml_lines[i]}"))
 # echo "echoing $i ${result_array[0]}"
  diffName="${result_array[0]}"
 # echo "check ${diffName}"
  if [[ "${diffName}" == "metrics" ]]; then
    i1=$i
    #echo "how enetred here if diffname is ${diffName}"
    break
  fi
done

if [[ "$i1" -ne -1 ]]; then
  i2=$i1
  for((i=i1+1; i<diff_length; i++)); do
    result_array=($(string_to_array "${diff_xml_lines[i]}"))
   # echo "${result_array[0]}"
   # echo "echoing $i ${result_array[0]}"
    diffName="${result_array[0]}"
    if [[ ${diffName} == "metrics" ]]; then
      i2=$i
      break
    fi
  done

  sliced_array=("${diff_xml_lines[@]:$i1:$((i2 - i1 + 1))}")
  diff_xml_lines=("${diff_xml_lines[@]:0:$i1}" "${diff_xml_lines[@]:$((i2+1))}")
  i1=-1
  for ((i=0; i<target_length; i++)); do
    result_array=($(string_to_array "${target_xml_lines[i]}"))
    diffName="${result_array[0]}"
    if [[ "${diffName}" == "metrics" ]]; then
      i1=$i
      break
    fi
  done
  i2=i1 ;

  target_xml_lines=("${target_xml_lines[@]:0:$i1}" "${target_xml_lines[@]:$((i2+1))}")
  target_xml_lines=("${target_xml_lines[@]:0:2}" "${sliced_array[@]}" "${target_xml_lines[@]:2}")

fi

i1=-1
i2=-1
diff_length=${#diff_xml_lines[@]}
target_length=${#target_xml_lines[@]}
for ((i=0; i<diff_length; i++)); do
  result_array=($(string_to_array "${diff_xml_lines[i]}"))
 # echo "echoing $i ${result_array[0]}"
  diffName="${result_array[0]}"
 # echo "check ${diffName}"
  if [[ "${diffName}" == "caches" ]]; then
    i1=$i
#    echo "how enetred here if diffname is ${diffName} $i"
    break
  fi
done
if [[ "$i1" -ne -1 ]]; then
  i2=$i1
  for((i=i1+1; i<diff_length; i++)); do
    result_array=($(string_to_array "${diff_xml_lines[i]}"))
   # echo "${result_array[0]}"
   # echo "echoing $i ${result_array[0]}"
    diffName="${result_array[0]}"
    if [[ ${diffName} == "caches" ]]; then
      i2=$i
      break
    fi
  done
  sliced_array=("${diff_xml_lines[@]:$i1:$((i2 - i1 + 1))}")
  diff_xml_lines=("${diff_xml_lines[@]:0:$i1}" "${diff_xml_lines[@]:$((i2+1))}")
#  echo "i1 and i2--------------------------------------> $i1 $i2 ${sliced_array[@]}"
  target_xml_lines=("${target_xml_lines[@]:0:2}" "${sliced_array[@]}" "${target_xml_lines[@]:2}")
fi

i1=-1
i2=-1
diff_length=${#diff_xml_lines[@]}
target_length=${#target_xml_lines[@]}
for ((i=0; i<diff_length; i++)); do
  result_array=($(string_to_array "${diff_xml_lines[i]}"))
 # echo "echoing $i ${result_array[0]}"
  diffName="${result_array[0]}"
 # echo "check ${diffName}"
  if [[ "${diffName}" == "logging" ]]; then
    i1=$i
#    echo "how enetred here if diffname is ${diffName} $i"
    break
  fi
done
if [[ "$i1" -ne -1 ]]; then
  i2=$i1
  for((i=i1+1; i<diff_length; i++)); do
    result_array=($(string_to_array "${diff_xml_lines[i]}"))
   # echo "${result_array[0]}"
   # echo "echoing $i ${result_array[0]}"
    diffName="${result_array[0]}"
    if [[ ${diffName} == "logging" ]]; then
      i2=$i
      break
    fi
  done
  sliced_array=("${diff_xml_lines[@]:$i1:$((i2 - i1 + 1))}")
  diff_xml_lines=("${diff_xml_lines[@]:0:$i1}" "${diff_xml_lines[@]:$((i2+1))}")
#  echo "i1 and i2--------------------------------------> $i1 $i2 ${sliced_array[@]}"
  target_xml_lines=("${target_xml_lines[@]:0:2}" "${sliced_array[@]}" "${target_xml_lines[@]:2}")
fi


diff_length=${#diff_xml_lines[@]}
target_length=${#target_xml_lines[@]}

echo "$target_length target len "
for ((i=0; i<diff_length; i++)); do
  result_array=($(string_to_array ${diff_xml_lines[i]}))
  ln=${#result_array[@]}
  diffName=($(get_name "${result_array[@]}"))
  type=${result_array[0]}
  echo $type
  if [[ $type != "bool" && $type != "int" && $type != "str" ]]; then
    continue
  fi
  echo "line no. 106 type name of diff ${diffName}"

  # element of diff has name property
  if [[ "$diffName" != ""  ]] ; then
    ff=0
    # check if same element exist in target xml or not
    for ((j=0; j<target_length; j++)); do
      result_array=($(string_to_array "${target_xml_lines[j]}"))
      echo "line no. 114. result array for $j th element of target ${result_array[@]} ${#result_array[@]}"
      targetName=($(get_name "${result_array[@]}"))
      echo "line no. 116. target name for result array is    ${targetName}"
      if [[ "$targetName" != "" ]]; then
        echo "line no. 118 targetname is ${targetName} && diffname is ${diffName}"
        if [[ "$targetName" == "$diffName" ]]; then
          # if exists then break
          target_xml_lines[$j]=${diff_xml_lines[$i]}
          echo "line no. 122. checking if target xml line assigned or not ${target_xml_lines[$j]} haha ${diff_xml_lines[$i]}"
          ff=1
          break
        fi
      fi
    done

    echo "line no. 129 value of ff ${ff} $ff"

    if [[ "$ff" == 1 ]]; then
      continue
    fi

    pre=""
    post=""
    i1=0
    i2=0

    #check previous element headers which has a type different of int, str, bool
    for (( j=i; j>=0; j-- )); do
      result_array=($(string_to_array "${diff_xml_lines[j]}"))
      echo "${result_array[0]}"
      diffName="${result_array[0]}"
      echo "line no. 145 $j th element of diff_xml_lines is ${diff_xml_lines[j]} && first elemement of result array is ${result_array[0]}"
      if [[ "$diffName" != "int" && "$diffName" != "str" && "$diffName" != "bool" ]]; then
        pre="$diffName"
        i1=$j
        break
      fi
    done

    #check post element headers which has a type different of int, str, bool
    for (( j=i; j<diff_length; j++ )); do
      result_array=($(string_to_array "${diff_xml_lines[j]}"))
      echo "${result_array[0]}"
      diffName="${result_array[0]}"
      echo "line no. 157 $j th element of diff_xml_lines is ${diff_xml_lines[j]} && first elemement of result array is ${result_array[0]}"
      if [[ "$diffName" != "int" && "$diffName" != "str" && "$diffName" != "bool" ]]; then
        post="$diffName"
        i2=$j
        break
      fi
    done

    echo "line no. 164. pre element ${pre} and post element ${post}"

    if [[ "$pre" != "$post" ]]; then
      pre="solr"
      post="solr"
    fi
    insert_index=-1
    for (( j=0; j<target_length; j++ )); do
      result_array=($(string_to_array "${target_xml_lines[j]}"))
      echo "line no. 173 result array for $j th element of target_xml_lines ${target_xml_lines[j]} is ** ${result_array[@]} ** and first element is **${result_array[0]}"
      type="${result_array[0]}"
      echo "line no. 175 type is $type && pre is $pre insert index is $insert_index"
      if [[ "$type" == "$pre" ]]; then
        insert_index=$j
        break
      fi
    done

    if [[ "$insert_index" -eq -1 ]]; then
      pre="solr"
      for (( j=0; j<target_length; j++ )); do
        result_array=($(string_to_array "${target_xml_lines[j]}"))
        echo "line no. 186 result array for $j th element of target_xml_lines ${target_xml_lines[j]} is ** ${result_array[@]} ** and first element is **${result_array[0]}"
        type="${result_array[0]}"
        echo "line no. 188 type is $type && pre is $pre insert index is $insert_index"
        if [[ "$type" == "$pre" ]]; then
          insert_index=$j
          break
        fi
      done

      ((insert_index++))
      echo "line no. 196. insert index number if $insert_index"
      echo "line no. 197. $i1 th, $i thm $i2 th element of diff_xml_lines recpectively ${diff_xml_lines[i1]} ${diff_xml_lines[i]} ${diff_xml_lines[i2]}"
      my_array=(${diff_xml_lines[i1]} ${diff_xml_lines[i]} ${diff_xml_lines[i2]})
      echo "line no. 199. my array ${my_array[@]}"
      target_xml_lines=("${target_xml_lines[@]:0:$insert_index}" "${my_array[@]}" "${target_xml_lines[@]:$insert_index}")
      continue
    fi



    new_element="${diff_xml_lines[i]}"
    ((insert_index++))
    echo "line no. 208. $i th element of diff_xml_lines is ${diff_xml_lines[i]} so new element is $new_element insert index is $insert_index"


    for ((j=0; j<target_length; j++)); do
      result_array=($(string_to_array "${target_xml_lines[j]}"))
      type="${result_array[0]}"
      echo "line no. 214 type is $type && pre is $pre insert index is $insert_index"
      if [[ $type == $pre ]]; then
        target_xml_lines=("${target_xml_lines[@]:0:$insert_index}" "$new_element" "${target_xml_lines[@]:$insert_index}")
        break
      fi
    done
  fi
done

#for((i=0; i<10; i++));do
#  echo ""
#done

echo done

echo "${target_xml_lines[@]}" > output.xml

xmllint --format output.xml > format.xml

sudo apt-get update
sudo apt-get install -y libxml2-utils

for item in "${target_xml_lines[@]}"; do
    # Check if the item is not an empty string
    echo "${item}"
  done

