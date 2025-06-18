const WS_URL = "../api/rnd/ws";

let socket = null;

const VoiceName = "Microsoft Zira - English (United States)"

function SetNewMessage(message) {
    window.speechSynthesis.cancel();
    let params = new SpeechSynthesisUtterance(message);
    params.voice = window.speechSynthesis.getVoices().find(voice => voice.name === VoiceName);
    window.speechSynthesis.speak(params);
    const MessageBox = document.getElementById("subtitle");
    MessageBox.innerHTML = message;
}

function Disconnect() {
    SetNewMessage(" ");
    document.getElementById("ava_img").src = "";
    document.getElementById("nickname").innerHTML = "";
}

function NewUser(name, pfp) {
    console.log("New user!");
    document.getElementById("nickname").innerHTML = name;
    document.getElementById("ava_img").src = pfp; 
}

function HandleMessage(message) {
    let data = JSON.parse(message.data);
    console.log(data);
    switch(data.type) {
        case "message":
            SetNewMessage(data.message);
            break;
        case "disconnect":
            Disconnect();
            break;
        case "connect":
            NewUser(data.username, data.userpfp);
            break;
    }
}   

function Connect() {
    socket = new WebSocket(WS_URL);
    socket.onerror = OnWSError;
    socket.onclose = OnWSError;
    socket.onmessage = HandleMessage;
    document.getElementById("lost_connection").style.visibility = "hidden"; 

    /*
    console.log("Voices:");
    const voices = speechSynthesis.getVoices();
    for (const voice of voices) {
        console.log(voice.name);
    }
    */

}

function OnWSError(err) {
    document.getElementById("lost_connection").style.visibility = "visible"; 
    console.log("Reconnecting to ws due to error. Error: " + err);
    console.log("Reconnecting...");
    Connect();
}
