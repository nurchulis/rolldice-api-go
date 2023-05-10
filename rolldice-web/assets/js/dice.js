var elDiceOne       = document.getElementById('dice1');
var elDiceTwo       = document.getElementById('dice2');
var elComeOut       = document.getElementById('roll');

elComeOut.onclick   = function () {rollDice();};

function rollDice() {



  // Create the JSON object
  let data = {
    dice: 2,
    userid: "1",
    event: "rolldice",
    bet: "big",
    mode: "singleplayer",
    bet_point: 500,
    session_id: 1
  };

  // Convert the JSON object to a string
  let jsonData = JSON.stringify(data);

  // Send the JSON string over the socket connection
  socket.send(jsonData);
}

let socket = new WebSocket("ws://127.0.0.1:8080/play");

socket.onopen = function(e) {
  alert("Connected To Server");
};

socket.onmessage = function(event) {
  alert(`[message] Data received from server: ${event.data}`);
  let receivedData = JSON.parse(event.data);
  var diceOne   = receivedData.dice_total/2
  var diceTwo   = receivedData.dice_total/2
 
  console.log(diceOne + ' ' + diceTwo);

  for (var i = 1; i <= 6; i++) {
    elDiceOne.classList.remove('show-' + i);
    if (diceOne === i) {
      elDiceOne.classList.add('show-' + diceOne);
    }
  }

  for (var k = 1; k <= 6; k++) {
    elDiceTwo.classList.remove('show-' + k);
    if (diceTwo === k) {
      elDiceTwo.classList.add('show-' + diceTwo);
    }
  } 
};

socket.onclose = function(event) {
  if (event.wasClean) {
    alert(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
  } else {
    // e.g. server process killed or network down
    // event.code is usually 1006 in this case
    alert('[close] Connection died');
  }
};

socket.onerror = function(error) {
  alert(`[error]`);
};
