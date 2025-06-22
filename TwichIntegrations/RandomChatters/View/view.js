const WS_URL = "../api/rnd/wsconn";

let socket = null;

const VoiceName = "Microsoft Mark - English (United States)"
let currentUtteracne = null

function SetNewMessage(message) {
    window.speechSynthesis.cancel();
    let params = new SpeechSynthesisUtterance(message);
    params.voice = window.speechSynthesis.getVoices().find(voice => voice.name === VoiceName);
    params.lang = "en-US"
    currentUtteracne = params;
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
    SetNewMessage(" ");
    document.getElementById("nickname").innerHTML = name;
    document.getElementById("ava_img").src = pfp; 
}

function HandleMessage(message) {
    try {
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
            case "heartbeat":
                break;
        }
    } catch(err) {
        console.log(err);
    }
}   

function Connect() {
    socket = new WebSocket(WS_URL);
    socket.onopen = () => {  document.getElementById("lost_connection").style.visibility = "hidden"; }
    socket.onerror = OnWSError;
    socket.onclose = OnWSClose;
    socket.onmessage = HandleMessage;

    // Keepalive-messages
    setInterval(() => {
        if(socket != undefined && socket != null && socket.readyState == 1)
            socket.send(`{"type":"heartbeat"}`);
    }, 5000);

    // Keepalive
    window.speechSynthesis.addEventListener("onend", () => {
        console.log("ended");
    });
}

function OnWSClose(err) {
    console.log("Reconnecting to ws due to error. Error: " + err);
    console.log("Reconnecting...");
    document.getElementById("lost_connection").style.visibility = "visible"; 
    setTimeout(Connect, 1000)
}

function OnWSError(err) {
    console.log("Disconnected from ws due to error. Error: " + err);
    socket.close()
}

var roatationIndex = 0;
const roatationBoundary = 2; // in degrees
var prevSpeechValue = false

setInterval(() => {
    SetPfpRotation();
}, 1000 / 30);
function SetPfpRotation() {
    if(!window.speechSynthesis.speaking) {
        if(prevSpeechValue) 
            document.getElementById("ava_img").style.transform = `rotate(0deg)`;
        prevSpeechValue = false;
        return;
    };
    var coeff = Math.sin(roatationIndex);
    roatationIndex = (roatationIndex + 5) % (2 * Math.PI);
    document.getElementById("ava_img").style.transform = `rotate(${roatationBoundary * coeff}deg)`;
    prevSpeechValue = true;
    // console.log(`${roatationIndex} ${roatationBoundary * coeff}`);
}
