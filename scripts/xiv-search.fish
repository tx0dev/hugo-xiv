#!/usr/bin/fish

# SPDX-FileCopyrightText: 2025 Mathieu C. <mathieu@tx0.dev>
#
# SPDX-License-Identifier: MIT

function xivapi_search_map
  set -l sheet $argv[1]
  set -l query $argv[2]
  set -l fields $argv[3]

  # Robust argument checks, also checking for EMPTY strings
  if test -z "$sheet"
    echo "Usage: xivapi_search_map <sheet> '<query>' '<fields>' "
    echo "       <sheet>:   The sheet name to search (e.g., Map)."
    echo "       <query>:   The search query string (e.g., 'TerritoryType.ContentFinderCondition.ContentType.value=0 -IsPvpZone=true'). Enclose in single quotes."
    echo "       <fields>:  Comma-separated list of fields to display (e.g., row_id,PlaceName.Name,Id). Enclose in single quotes."
    return 1
  end

  if test -z "$query"
    # Refined query to filter for travelable maps/zones
    echo "Query is missing"
    return 1
    #set query "(TerritoryType.TerritoryIntendedUse=0 TerritoryType.TerritoryIntendedUse=1) TerritoryType.ContentFinderCondition.ContentType=0 TerritoryType.ExclusiveType=0 -TerritoryType.TerritoryIntendedUse=14 -IsEvent=true"
  end
  if test -z "$fields"
    # Updated fields to include desired information
    set fields 'row_id,PlaceName.Name,PlaceNameRegion.Name,Id,TerritoryType.ExVersion.Name,PlaceNameIcon.id,TerritoryType.TerritoryIntendedUse,TerritoryType.ExclusiveType'
  end

  set -l base_url "https://v2.xivapi.com/api/search"
  set -l encoded_query (string replace -a '=' '%3D' (string replace -a ' ' '%20' (string replace -a '>' '%3E' (string replace -a '<' '%3C' (string replace -a '    ' '' "$query")))))
  set -l search_url "$base_url?sheets=$sheet&query=$encoded_query&fields=$fields"

  echo "Querying XIVAPI:" $search_url


  set -l json_output
  if not jq --version > /dev/null 2>&1
    echo "Error: jq is not installed. Please install jq to use this script."
    return 1
  end

  set json_output (curl -s "$search_url")

  if test -z "$json_output"
    echo "Error: No response from XIVAPI. Check your internet connection and URL."
    return 1
  end

  # Error handling for invalid JSON response (using a temp file and then checking exit code)
  set -l temp_json_file (mktemp)
  echo "$json_output" > $temp_json_file
  if not jq -e '.' $temp_json_file > /dev/null 2>&1
    echo "Error: Invalid JSON response from XIVAPI."
    echo "Raw response:"
    cat $temp_json_file
    # Try to parse error message from XIVAPI if possible
    if jq '.message' $temp_json_file > /dev/null 2>&1
      set -l api_error_message (jq -r '.message' $temp_json_file)
      echo "XIVAPI Error Message:" $api_error_message
    end
    rm $temp_json_file
    return 1
  end
  rm $temp_json_file


  echo "" # Add an empty line for better readability before the table

  echo "Id | MapId     | Expansion     | Region              | Name                          |Use|Type"
  #     859  y6f3/00     Dawntrail       Yok Tural             Yak T'el                        1   0

  echo "$json_output" | jq -r '.results[] |
      [
        .row_id,
        .fields.Id,
        .fields.TerritoryType.fields.ExVersion.fields.Name,
        .fields.PlaceNameRegion.fields.Name,
        .fields.PlaceName.fields.Name,
        .fields.TerritoryType.fields.TerritoryIntendedUse.value,
        .fields.TerritoryType.fields.ExclusiveType
      ] |
      join("|") ' | column -s '|' -t
end

# Correct function call (using defaults from inside the function now)
# set q "(TerritoryType.TerritoryIntendedUse=0 TerritoryType.TerritoryIntendedUse=1) \
#     TerritoryType.ContentFinderCondition.ContentType=0 \
#     TerritoryType.ExclusiveType=0 \
#     -TerritoryType.TerritoryIntendedUse=14 \
#     -IsEvent=true"
set q "Hierarchy=1 \
    +TerritoryType.ContentFinderCondition.ContentType=0 \
    -TerritoryType.ContentFinderCondition.ContentType=2 \
    +TerritoryType.TerritoryIntendedUse<2 \
    +TerritoryType.IsPvpZone=false"
xivapi_search_map 'Map' "$q"
