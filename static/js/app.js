function execBlock(blockId) {
    const params = new URLSearchParams({
        bid: blockId,
    });
    fetch(window.location.pathname+"?"+params.toString(), {
        method: "POST",
    }).then((data) => {
        return data.text()
    }).then((v) => {
        document.getElementById(blockId).innerHTML =
            `<div id="${blockId}"><textarea readonly="true">${v}</textarea><br /><button onclick="execBlock('${blockId}')">Rerun</button><button onclick="clearBlock('${blockId}')">Clear</button></div>`

    })
}

function clearBlock(blockId) {
    document.getElementById(blockId).innerHTML =
        `<div id="${blockId}"><button onclick="execBlock('${blockId}')">Run</button></div>`
}

