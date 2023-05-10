Roll Dice Game
This is a simple Roll Dice game implemented in Go, using sockets for multiplayer functionality.

Requirements
Go 1.16 or higher
Installation
Clone the repository: git clone https://github.com/your-username/roll-dice-game.git
Navigate to the project directory: cd roll-dice-game
Build the project: go build
Run the server: ./roll-dice-game server
In a separate terminal window, run the client: ./roll-dice-game client
Usage
Once the server and client are running, the game will automatically start. Players will take turns rolling the dice, and the player with the highest total score after a set number of rounds wins.

Configuration
The following environment variables can be set to configure the game:

ROLL_DICE_GAME_PORT: The port number for the server to listen on. Default is 8080.
ROLL_DICE_GAME_ROUNDS: The number of rounds to play in the game. Default is 5.
Contributing
If you'd like to contribute to this project, please follow these steps:

Fork the repository
Create a new branch: git checkout -b my-feature-branch
Make your changes and commit them: git commit -am "Add new feature"
Push to the branch: git push origin my-feature-branch
Create a pull request
License
This project is licensed under the MIT License. See the LICENSE file for details.