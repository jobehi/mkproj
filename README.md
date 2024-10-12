# mkprj

Welcome to mkprj! This script helps you create and navigate a project directory structure interactively. It allows you to create directories and files, move between directories, and visualize the current project structure.

## Features

- **Interactive Navigation**: Navigate through your project directories interactively.
- **Directory and File Creation**: Create directories and files within your project structure.
- **Visual Feedback**: Display the current project structure using the `tree` command or a fallback method.
- **Color-Coded Output**: Enhanced readability with color-coded messages.

## Prerequisites

- Bash shell
- `tree` command (optional, for better visualization)

## Usage

1. **Clone the Repository**:
    ```sh
    git clone <repository-url>
    cd <repository-directory>
    ```

2. **Run the Script**:
    ```sh
    ./mkprj.bash
    ```

3. **Follow the Interactive Prompts**:
    - Enter the root directory name for your project.
    - Use the menu options to create directories, create files, move between directories, and finish the directory setup.

## Menu Options

- **1. Create a directory**: Prompts for a directory name and creates it within the current path.
- **2. Create a file**: Prompts for a file name and creates it within the current path.
- **3. Move back to the previous directory**: Moves back to the previously navigated directory.
- **4. Switch to an existing directory**: Prompts for a directory path and switches to it if it exists.
- **5. Finish this directory**: Exits the current directory and completes the setup.

## Example

```sh
Welcome to the Interactive Project Structure Builder!
Enter the root directory name for your project: my_project
Root directory 'my_project' created.
Current project structure:
my_project
What would you like to do?
1. Create a directory
2. Create a file
3. Move back to the previous directory
4. Switch to an existing directory
5. Finish this directory
Enter your choice (1/2/3/4/5): 1
Enter directory name: src
Directory 'src' created.