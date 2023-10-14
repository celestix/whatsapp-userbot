# WhatsApp Userbot in Go

This is a WhatsApp userbot written in the Go programming language. A userbot is a script or program that allows you to automate various tasks on WhatsApp. With this userbot, you can perform a wide range of actions such as sending automated responses, and managing your WhatsApp account programmatically, etc.

## Features

- **Message Automation:** You can use this userbot to automate the process of sending and receiving messages. For example, you can set AFK messages to be sent automatically when you are away from your phone.

- **Auto Replies:** Set up auto-replies for incoming messages. You can define custom responses for specific keywords or phrases.

- **Group Management:** Manage groups you are a part of, including adding or removing members, changing group settings, and more.

- **Custom Commands:** You can use this userbot to execute code snippets which are run on piston. This allows you to run code on your phone without having to install any additional software. 

## Getting Started

### Prerequisites

Before you can use this WhatsApp userbot, you will need to meet the following requirements:

- **WhatsApp Account:** You must have a WhatsApp account with a phone number.

- **Go Programming Language:** You need to have Go installed on your system. You can download and install it from the official Go website: [Go Downloads](https://golang.org/dl/)

### Installation

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/celestix/whatsapp-userbot
   ```

2. Navigate to the project directory:

   ```bash
   cd whatsapp-userbot
   ```

3. Install the dependencies:

    ```bash
    go mod tidy
    ```

4. Build the project:

    ```bash
    go build .
    ```

5. Run the the userbot:

    ```bash
    ./whatsapp-userbot
    ```

6. Scan the QR code displayed in the console with your phone to log in to your WhatsApp account.

## Usage

Once the userbot is running, you can interact with it through WhatsApp. Send commands and messages to your WhatsApp account, and the userbot will respond accordingly.

Some example commands:

`.ping`: Check if the userbot is alive or not.

`.help`: Display a help message with available commands.

`.add <note_name> <reply_to_message>`: Save a text note in the database.

`.afk <<on {reason}>/off>`: Set AFK status.

Please refer to the the help command for a complete list of available commands and their usage.

## Contributing

Contributions are welcome! If you would like to contribute to this project, please open an issue to discuss the changes you would like to make.

## License
[![GNU Affero General Public License v3.0](https://www.gnu.org/graphics/agplv3-155x51.png)](https://www.gnu.org/licenses/agpl-3.0.en.html#header)    
Licensed under GNU AGPL v3.0.   
Selling the codes to other people for money is *strictly prohibited*.