let preferences = {}
let history = ""

function toggleDarkMode() {
    const body = document.body
    body.classList.toggle('dark-mode')
}

function loadPreferences() {
    localStorage.getItem('preference') ? preferences = JSON.parse(localStorage.getItem('preference')) : preferences = {}

    return preferences   
}

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

async function websocketSetup() {
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

        requestMaxPhycalCPU()
        requestMaxLogicalCPU()

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

        const responses = evt.data.split("\n$$__SEPARATOR__$$\n") // Split by separator
        let response = ""

        switch (true) {
            case responses[0].includes("$$__RESPONSE_PREDICT__$$"): // Response prediction to output screen
                response = responses[1].replace(/\n/g, "<br />")
                // Catch the end of the response
                if (response.includes("$$__RESPONSE_DONE__$$")) {
                    response = response.replace(/\$\$__RESPONSE_DONE__\$\$/g, "")
                    response = response.slice(0, -12) // remove the last <br /><br />, not \n\n

                    statusRemoveAllClasses()
                    statusAddClass('status-ready', 'Ready')
                    buttonSendEnable()

                    break
                }

                out.innerHTML = history + response + `<span class="cursor"></span>`
                out.scrollTop = out.scrollHeight

                break
            // case response.startsWith("$$__ERROR__$$"):
            case responses[0].includes("$$__RESPONSE_INFO__$$"): // Response info to console log
                // 0: $$__RESPONSE_INFO__$$, 1: $$__MAX_CPU_PHYSICAL__$$ or $$__MAX_CPU_LOGICAL__$$, 2: number
                switch (responses[1]) {
                    case "$$__MAX_CPU_PHYSICAL__$$":
                        console.log(`Max physical CPU: ${responses[2]}`)
                        break
                    case "$$__MAX_CPU_LOGICAL__$$":
                        console.log(`Max logical CPU: ${responses[2]}`)
                        break
                }

                break
            default:
                console.log(`Unknown response: ${evt.data}`)
        }
    }
}

function send() {
    const input = document.querySelector("#inputs")
    if (input.value === '') {
        return
    }

    history = document.querySelector("#outputs").innerHTML
    history = history.slice(0, -28)

    const reflection = document.querySelector("#reflection")
    const antiprompt = document.querySelector("#antiprompt")

    const data = `$$__PROMPT__$$\n$$__SEPARATOR__$$\n${input.value}\n$$__SEPARATOR__$$\n${reflection.value}\n$$__SEPARATOR__$$\n${antiprompt.value}`

    ws.send(data)

    statusRemoveAllClasses()
    statusAddClass('status-running', 'Running')
    buttonStopEnable()

    input.value = ''
    input.focus()
}

function requestMaxPhycalCPU() {
    const input = document.querySelector("#inputs")
    const cpuNummsg = "$$__COMMAND__$$\n$$__SEPARATOR__$$\n$$__MAX_CPU_PHYSICAL__$$"

    ws.send(cpuNummsg)
}

function requestMaxLogicalCPU() {
    const input = document.querySelector("#inputs")
    const cpuNummsg = "$$__COMMAND__$$\n$$__SEPARATOR__$$\n$$__MAX_CPU_LOGICAL__$$"

    ws.send(cpuNummsg)
}

function stopResponse() {
    const input = document.querySelector("#inputs")
    const stopmsg = "$$__COMMAND__$$\n$$__SEPARATOR__$$\n$$__STOP__$$"

    ws.send(stopmsg)

    statusRemoveAllClasses()
    statusAddClass('status-running', 'Running')
    buttonSendEnable()

    input.value = ''
    input.focus()
}

function init() {
    preferences = loadPreferences()

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
