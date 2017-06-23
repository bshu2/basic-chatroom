new Vue({
    el: "#app",

    data: {
        ws: null, //websocket
        username: null,
        message: "", 
        joined: false 
    },
    created: function() {
        if (!window["WebSocket"]) {
            var new_message = document.createElement("div");
            new_message.innerText = "Websockets not supported";
            new_message.className = "system-message";
            chatbox.appendChild(new_message);
            console.log("Websockets not supported");
            return;
        }  
        var that = this;
        this.ws = new WebSocket("wss://" + window.location.host + "/ws");

        //handle messages
        this.ws.onmessage = function(e) {
            var chatbox = document.getElementById("chat-messages");
            var json_message = JSON.parse(e.data);
            var new_message = document.createElement("div");
            switch (json_message.message_type) {
                case 0://CHAT_MESSAGE
                    new_message.innerText = json_message.username + ": " + json_message.message_text;
                    chatbox.appendChild(new_message);
                    break;
                case 1://SYSTEM_MESSAGE
                    new_message.innerText = json_message.message_text;
                    new_message.className = "system-message";
                    chatbox.appendChild(new_message);
                    break;
                default:
                    break;
            }
            chatbox.scrollTop = chatbox.scrollHeight - chatbox.clientHeight; //automatic scroll
        };

        //message when websocket closes
        this.ws.onclose = function(e) {
            var chatbox = document.getElementById("chat-messages");
            var disconnected_message = document.createElement("div");
            disconnected_message.innerText = "connection lost";
            disconnected_message.className = "system-message";
            chatbox.appendChild(disconnected_message);
        }
    },
    methods: {
        //send a message
        send: function () {
            if (this.message != "") {
                this.ws.send(JSON.stringify({
                        username: this.username,
                        message_text: this.message
                    }
                ));
                this.message = "";
            }
        },
        //join the room
        join: function () {
            if (!this.username) {
                return
            }
            this.joined = true;
            this.ws.send(JSON.stringify({
                    username: this.username,
                    message_text: ""
                }
            ));
        }
    }
});