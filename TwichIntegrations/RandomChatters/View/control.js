const ConnectEndpoint = "../api/rnd/connect"
const DisconnectEndpoint = "../api/rnd/disconnect"
const BanEndpoint = "../api/rnd/ban"
const PardonEndpoint = "../api/rnd/pardon"
const BanTableEndpoint = "../api/rnd/bannedusers"
const CurrentUserEndpoint = "../api/rnd/currnetchatter"

function CallAPI(endpoint) {
    let API_Key = document.getElementById("ApiKeyField").value;
    let uri = `${endpoint}?token=${API_Key}`
    fetch(uri).then(x => {
        if(x.status != 200) {
            console.error("Invalid token")
        }
        FillUsernameSpace();
    });
}

function AddBannedUser(username) {
    let record = document.createElement("tr");
    record.id = `${username}_rec`
    let unbanButtonPlace = document.createElement("td");
    let unbanButton = document.createElement("button");
    unbanButton.innerHTML = "Unban";
    unbanButton.onclick = () => { PardonUser(username) };
    let Username = document.createElement("td");
    Username.classList.add("nicks");
    Username.innerHTML = username;
    record.appendChild(unbanButtonPlace);
    unbanButtonPlace.appendChild(unbanButton);
    record.appendChild(Username);
    document.getElementById("bantable").appendChild(record);
}

function BanUser() {
    let API_Key = document.getElementById("ApiKeyField").value;
    fetch(CurrentUserEndpoint).then(r => {
        if(r.ok) r.json().then(data => {
            let username = data.username;
            if(username == "") return;
            console.log("Current user: " + username);
            fetch(`${BanEndpoint}?user=${username}&token=${API_Key}`).then((r) => {
                if(r.status == 200 && document.getElementById(`${username}_rec`) == undefined) {
                    AddBannedUser(username);
                }
            });
        });
    });
}

function FillBanTable() {
    fetch(BanTableEndpoint)
        .then(res => {
            if (res.status != 200) {
                console.error("Invalid token")
                return;
            }
            res.json().then(contents => {
                console.log(contents);
                for (var i in contents.chatters) {
                    AddBannedUser(contents.chatters[i]);
                }
            });
        });
}

function FillUsernameSpace() {
    let usernameSpace = document.getElementById("CurrentUser");
    fetch(CurrentUserEndpoint).then(r => {
        if(r.ok) r.json().then(data => {
            let username = data.username;
            if (username == "") {
                usernameSpace.innerHTML = "No user connected";
            } else {
                usernameSpace.innerHTML = `Current user: ${username}`;
            }
        });
    });
}

function ClearBanTable() {
    document.getElementById("bantable").innerHTML = '';
}

function PardonUser(username) {
    console.log("Unbanning " + username)
    let API_Key = document.getElementById("ApiKeyField").value;
    fetch(`${PardonEndpoint}?user=${username}&token=${API_Key}`).then(resp => {
        if (resp.status != 200) {
            console.error("Invalid token")
            return
        } else {
            let record = document.getElementById(`${username}_rec`);
            record.remove();
        }
    });
}
