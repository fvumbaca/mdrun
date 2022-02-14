var s = new WebSocket(((window.location.protocol === "https:") ? "wss://" : "ws://") + window.location.host + "/-/ws");

s.addEventListener('message', function (event) {
    console.log('Message from server ', event);
});

function execBase64(based) {
	// alert("Going to execute " + based);
    s.send(JSON.stringify({
		lang: "sh",
		script: "Hello World!",
	}));

}
