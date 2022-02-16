// var s = new WebSocket(
//     ((window.location.protocol === "https:") ? "wss://" : "ws://") +
//     window.location.host + "/-/ws");

// s.addEventListener('message', function (event) {
//     console.log('Message from server ', event);
// });

function execBlock(blockId) {
    const params = new URLSearchParams({
        bid: blockId,
    });
    fetch(window.location.pathname+"?"+params.toString(), {
        method: "POST",
    }).then((data) => {
        return data.text()
    }).then((v) => {
        alert(v)
    })
}

