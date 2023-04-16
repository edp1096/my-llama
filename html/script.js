let history = ""

function statusRemoveAllClasses() {
    const status = document.getElementById('status')

    status.classList.remove('status-disconn')
    status.classList.remove('status-ready')
    status.classList.remove('status-running')
}

function statusAddClass(className, message) {
    const status = document.getElementById('status')
    status.classList.add(className)

    status.innerHTML = message
}

function buttonDisableAll() {
    document.querySelector("#send").disabled = true
    document.querySelector("#stop").disabled = true
}

function buttonSendEnable() {
    document.querySelector("#send").disabled = false
    document.querySelector("#stop").disabled = true
}

function buttonStopEnable() {
    document.querySelector("#send").disabled = true
    document.querySelector("#stop").disabled = false
}

function websocketSetup() {
    let param = "model_file=ggml-vicuna-7b-4bit.bin"

    const loc = window.location

    let uri = "ws:"
    if (loc.protocol === 'https:') {
        uri = "wss:"
    }

    uri += `//localhost:1323/ws?${param}`

    ws = new WebSocket(uri)

    ws.onopen = function () {
        console.log('Connected')

        statusRemoveAllClasses()
        statusAddClass('status-ready', 'Ready')

        document.getElementById("send").click()
    }

    ws.onclose = function () {
        console.log('Disconnected')

        statusRemoveAllClasses()
        statusAddClass('status-disconn', 'Disconnected')
        buttonDisableAll()
    }

    ws.onmessage = function (evt) {
        const out = document.getElementById('outputs')
        buttonStopEnable()

        let response = ""
        response = evt.data.replace(/\n/g, "<br />")

        // catch the end of the response
        if (response.includes("$$__RESPONSE_DONE__$$")) {
            response = response.replace(/\$\$__RESPONSE_DONE__\$\$/g, "")
            response = response.slice(0, -12) // remove the last <br /><br />, not \n\n

            statusRemoveAllClasses()
            statusAddClass('status-ready', 'Ready')
            buttonSendEnable()
        }

        out.innerHTML = history + response

        out.scrollTop = out.scrollHeight
    }
}

function send() {
    const input = document.querySelector("#inputs")
    if (input.value === '') {
        return
    }

    history = document.querySelector("#outputs").innerHTML

    const reflection = document.querySelector("#reflection")
    const antiprompt = document.querySelector("#antiprompt")

    const data = `${input.value}\n$$__SEPARATOR__$$\n${reflection.value}\n$$__SEPARATOR__$$\n${antiprompt.value}`

    ws.send(data)

    statusRemoveAllClasses()
    statusAddClass('status-running', 'Running')
    buttonStopEnable()

    input.value = ''
    input.focus()
}

function stopResponse() {
    const input = document.querySelector("#inputs")
    const stopmsg = "$$__STOP__$$"

    ws.send(stopmsg)

    statusRemoveAllClasses()
    statusAddClass('status-running', 'Running')
    buttonSendEnable()

    input.value = ''
    input.focus()
}

function init() {
    websocketSetup()

    let promptTEXT = `
A chat between a curious human and an artificial intelligence assistant. The assistant gives helpful, detailed, and polite answers to the human's questions.

### Human: Hello, Assistant.
### Assistant: Hello. How may I help you today?
### Human: Please tell me the largest city in Europe.
### Assistant: Sure. The largest city in Europe is Moscow, the capital of Russia.
### Human:`
    let antipromptTEXT = `### Human:`

    document.querySelector("#reflection").value = promptTEXT
    document.querySelector("#antiprompt").value = antipromptTEXT

    buttonSendEnable()

    document.getElementById("inputs").addEventListener("keydown", function (event) {
        if (event.ctrlKey && event.keyCode === 13) {
            event.preventDefault();
            document.getElementById("send").click()
        }
    })
}

document.addEventListener('DOMContentLoaded', init())
