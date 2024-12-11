# protogui

**protogui** is a simple and intuitive graphical interface designed to encode and decode Protobuf messages. This tool makes it easy to work with binary messages in the Protocol Buffers format, offering a practical and efficient way to handle such data.

## Features

- **Encode**: Convert structured data into Protobuf binary messages.
- **Decode**: Transform Protobuf binary messages into readable formats.
- User-friendly and easy-to-navigate interface.

## Requirements

- [Golang](https://golang.org/) installed on your system (minimum recommended version: 1.18).
- This tool was developed using the Fyne library, please check prerequisites [here](https://docs.fyne.io/started/) 
   

## Installation

To install the tool, follow these steps:

1. Clone the **protogui** repository:
    ```bash
    git clone https://github.com/lawmatsuyama/protogui.git
    ```
2. Navigate to the repository's directory:
    ```bash
    cd protogui
    ```
3. Run the `go install` command to install the tool:
    ```bash
    go install
    ```

Once installed, the **protogui** executable will be available in your `$GOPATH/bin`.

## Usage

1. Open a terminal prompt.
2. Run the `protogui` command:
    ```bash
    protogui
    ```

The graphical interface will launch, allowing you to quickly encode and decode Protobuf messages.


