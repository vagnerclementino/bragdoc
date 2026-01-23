#!/usr/bin/env bash

add_empty_line() {
  echo "" >> $1
}

get_available_commands() {
  command_path=$1

  result=$(cat "$command_path" | awk '/Available Commands:/,EOF' | sed -n '/Flags:/q;p' | grep -v "Available Commands:" | awk '{print $1}' | sed '/^[[:space:]]*$/d' | grep -v help)
}

process_command() {
  command=$1
  path=$2

  command_path_out=$command.out

  go run ./cmd/cli $command --help > "$command_path_out"
  has_commands=$(cat "$command_path_out" | grep "Available Commands")

  add_empty_line $path
  if [[ -n "$has_commands" ]]; then
    echo "## Available Commands" >> $path
    add_empty_line $path

    echo '```bash' >> $path
    cat "$command_path_out" | awk '/Available Commands:/,EOF' | sed -n '/Flags:/q;p' | grep -v "Available Commands:" | sed '/^[[:space:]]*$/d' >> $path
    echo '```' >> $path

    get_available_commands $command_path_out
    for inner_command in $result ; do
      add_empty_line $path
      echo "### $inner_command" >> $path
      add_empty_line $path
      go run ./cmd/cli $command $inner_command --help > $inner_command.out
      cat $inner_command.out | sed 's/Usage:/```bash\nUsage:/' >> $path
      echo '```' >> $path

      has_commands=$(cat $inner_command.out | grep "Available Commands")
      if [[ -n "$has_commands" ]]; then
        get_available_commands $inner_command.out
        for double_inner_command in $result ; do
          add_empty_line $path
          echo "#### $double_inner_command" >> $path
          add_empty_line $path
          go run ./cmd/cli $command $inner_command $double_inner_command --help > $double_inner_command.out
          cat $double_inner_command.out | sed 's/Usage:/```bash\nUsage:/' >> $path
          echo '```' >> $path
          rm $double_inner_command.out
        done
      fi
      rm $inner_command.out
    done
  else
    echo "### $command" >> $path
    add_empty_line $path
    go run ./cmd/cli $command --help | sed 's/Usage:/```bash\nUsage:/' >> $path
    echo '```' >> $path
  fi

  rm "$command_path_out"
}

rm -rf docs/commands/*.md

go run ./cmd/cli --help > help.out

mkdir -p docs/commands

get_available_commands help.out
for command in $result ; do
  path="./docs/commands/$command.md"

  all_caps=$(echo $command | tr [a-z] [A-Z])
  capitalize="$(tr '[:lower:]' '[:upper:]' <<< ${command:0:1})${command:1}"
  if [[ ${#command} -le 3 ]]; then
    echo "# $all_caps" > $path
  else
    echo "# $capitalize" > $path
  fi

  process_command $command $path
done

rm help.out
