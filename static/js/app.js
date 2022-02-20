function execBlock(blockId) {
    const params = new URLSearchParams({
        bid: blockId,
    });
    fetch(window.location.pathname+"?"+params.toString(), {
        method: "POST",
    }).then((data) => {
        return Promise.all([data.status, data.text()])
    }).then((v) => {
        let outputClass = "";
        if (v[0] != 200) {
            outputClass = "class=\"output-error\"";
        }
        document.getElementById(blockId).innerHTML =
            `<div id="${blockId}"><textarea ${outputClass} readonly="true">${v[1]}</textarea><br /><button onclick="execBlock('${blockId}')">Rerun</button><button onclick="clearBlock('${blockId}')">Clear</button></div>`
    })
}

function clearBlock(blockId) {
    document.getElementById(blockId).innerHTML =
        `<div id="${blockId}"><button onclick="execBlock('${blockId}')">Run</button></div>`
}

