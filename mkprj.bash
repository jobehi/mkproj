#!/bin/bash

# Define color constants for visual clarity
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# ANSI escape sequences for clearing the screen and moving the cursor
CLEAR_SCREEN='\033[2J'
MOVE_CURSOR_HOME='\033[H'

# Track the previous directory for navigation
previous_dir=""

# Function to display the current project structure using the `tree` command
display_tree() {
  # Clear the screen and move the cursor to the top
  echo -e "${CLEAR_SCREEN}${MOVE_CURSOR_HOME}"
  echo -e "${CYAN}Current project structure:${NC}"

  # Display tree structure, fallback to find if tree is unavailable
  if command -v tree &> /dev/null; then
    tree "$1"
  else
    echo -e "${RED}Warning: 'tree' command not found. Using fallback.${NC}"
    find "$1" | sed 's/[^-][^\/]*\//   |/g;s/|\([^ ]\)/|-- \1/'
  fi
}

# Function for creating directories, files, and navigating directories
navigate_directories() {
  local current_path=$1

  while true; do
    display_tree "$root_dir"  # Display the current structure
    echo -e "${BLUE}Current path: $current_path${NC}"
    echo -e "${GREEN}What would you like to do?${NC}"
    echo "1. Create a directory"
    echo "2. Create a file"
    echo "3. Move back to the previous directory"
    echo "4. Switch to an existing directory"
    echo "5. Finish this directory"
    read -p "Enter your choice (1/2/3/4/5): " choice
    
    case "$choice" in
      1)
        # Create a directory
        read -p "Enter directory name: " dir_name
        mkdir -p "$current_path/$dir_name"
        echo -e "${GREEN}Directory '$dir_name' created.${NC}"
        previous_dir="$current_path"  # Store current directory before changing
        current_path="$current_path/$dir_name"  # Move into the new directory
        ;;
      
      2)
        # Create a file
        read -p "Enter file name (with extension, like file.py): " file_name
        touch "$current_path/$file_name"
        echo -e "${RED}File '$file_name' created.${NC}"
        ;;
      
      3)
        # Move back to the previous directory
        if [[ -n "$previous_dir" && -d "$previous_dir" ]]; then
          echo -e "${GREEN}Moving back to the previous directory: $previous_dir${NC}"
          current_path="$previous_dir"
          previous_dir=""
        else
          echo -e "${RED}No previous directory to go back to.${NC}"
        fi
        ;;
      
      4)
        # Switch to an existing directory
        read -p "Enter the path of the directory you want to switch to: " new_path
        if [[ -d "$new_path" ]]; then
          previous_dir="$current_path"  # Store current directory before switching
          current_path="$new_path"
          echo -e "${GREEN}Switched to directory: $new_path${NC}"
        else
          echo -e "${RED}Invalid path. Please enter a valid directory.${NC}"
        fi
        ;;
      
      5)
        # Finish and exit the current directory
        break
        ;;
      
      *)
        # Handle invalid choices
        echo -e "${RED}Invalid choice. Please try again.${NC}"
        ;;
    esac
  done
}

# Main program flow
echo -e "${CYAN}Welcome to the Interactive Project Structure Builder!${NC}"
read -p "Enter the root directory name for your project: " root_dir

# Create the root directory if it doesn't exist
mkdir -p "$root_dir"
echo -e "${GREEN}Root directory '$root_dir' created.${NC}"

# Start navigating directories
navigate_directories "$root_dir"

echo -e "${GREEN}Project structure created successfully!${NC}"
