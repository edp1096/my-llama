const defaultPromptTEXT = `
A chat between a curious human and an artificial intelligence assistant. The assistant gives helpful, detailed, and polite answers to the human's questions.

### Human: Hello, Assistant.
### Assistant: Hello. How may I help you today?
### Human: Please tell me the largest city in Europe.
### Assistant: Sure. The largest city in Europe is Moscow, the capital of Russia.
### Human:`
const defaultAntipromptTEXT = `### Human:`
const defaultFirstInputTEXT = `Please tell me the largest city in Asia.`

let preferences = {}

function toggleDarkMode(isSave = true) {
    const body = document.body
    const darkmodeSwitch = document.getElementById("switch-shade")

    switch (darkmodeSwitch.checked) {
        case true:
            document.getElementById("switch-shade-label").innerHTML = "Dark Mode"
            body.classList.add("dark-mode")
            break
        case false:
            document.getElementById("switch-shade-label").innerHTML = "Light Mode"
            body.classList.remove("dark-mode")
            break
    }

    if (isSave) {
        preferences["darkmode"] = darkmodeSwitch.checked
        savePreferences()
    }
}

function openPanel() {
    buttonDisableAll()
    const panel = document.getElementById("preferences")
    panel.style.display = "block"
}

function closePanel() {
    buttonSendEnable()
    const panel = document.getElementById("preferences")
    panel.style.display = "none"
}

function loadPreferences() {
    localStorage.getItem('preference') ? preferences = JSON.parse(localStorage.getItem('preference')) : preferences = {}

    document.querySelector("#pref_threads").value = preferences["threads"] ? preferences["threads"] : 1

    document.querySelector("#pref_top_k").value = preferences["TOP_K"] ? preferences["TOP_K"] : 40
    document.querySelector("#pref_top_p").value = preferences["TOP_P"] ? preferences["TOP_P"] : 0.8
    document.querySelector("#pref_temperature").value = preferences["TEMPERATURE"] ? preferences["TEMPERATURE"] : 0.5
    document.querySelector("#pref_repeat_penalty").value = preferences["REPEAT_PENALTY"] ? preferences["REPEAT_PENALTY"] : 0.75

    return preferences
}

function savePreferences() {
    localStorage.setItem('preference', JSON.stringify(preferences))
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
    document.querySelector("#open-panel").disabled = true
}

function buttonSendEnable() {
    buttonDisableAll()
    document.querySelector("#send").disabled = false
    document.querySelector("#open-panel").disabled = false
}

function buttonStopEnable() {
    buttonDisableAll()
    document.querySelector("#stop").disabled = false
}

async function websocketSetup() {
    // let param = "model_file=ggml-vicuna-7b-4bit.bin"
    let modelFile = "not_use_yet"

    if (preferences["MODEL_FILE"] != undefined && preferences["MODEL_FILE"] != "") {
        modelFile = preferences["MODEL_FILE"]
    }

    const loc = window.location

    let base = "ws:"
    if (loc.protocol === 'https:') {
        base = "wss:"
    }

    base += `//localhost:1323/ws`
    const params = new URLSearchParams({
        model_file: modelFile,
        n_ctx: preferences["N_CTX"] ? preferences["N_CTX"] : 512,
        n_batch: preferences["N_BATCH"] ? preferences["N_BATCH"] : 32,
    })

    const uri = new URL(`${base}?${params}`)

    ws = new WebSocket(uri)

    ws.onopen = function () {
        console.log('Connected')

        requestModelFiles()

        requestMaxPhycalCPU()
        requestMaxLogicalCPU()
        requestThreads()

        statusRemoveAllClasses()
        statusAddClass('status-ready', 'Ready')

        applyParameters()

        // document.getElementById("send").click()
        sendPrompt(true)
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

                // -28 for cursor, `<span class="cursor"></span>` removing.
                out.innerHTML = out.innerHTML.slice(0, -28) + response + `<span class="cursor"></span>`
                out.scrollTop = out.scrollHeight

                break
            // case response.startsWith("$$__ERROR__$$"):
            case responses[0].includes("$$__RESPONSE_INFO__$$"): // Response info to console log
                // 0: $$__RESPONSE_INFO__$$, 1: $$__MAX_CPU_PHYSICAL__$$ or $$__MAX_CPU_LOGICAL__$$, 2: number
                switch (responses[1]) {
                    case "$$__MODEL_FILES__$$":
                        for (let i = 2; i < responses.length; i++) {
                            const option = document.createElement("option")
                            option.text = responses[i]
                            option.value = responses[i]
                            document.querySelector("select[name=model_files]").add(option)

                            if (i == 2) {
                                document.querySelector("select[name=model_files]").value = responses[i]
                            }
                        }

                        if (preferences["MODEL_FILE"] != undefined && preferences["MODEL_FILE"] != "") {
                            document.querySelector("select[name=model_files]").value = preferences["MODEL_FILE"]
                        }

                        break
                    case "$$__MAX_CPU_PHYSICAL__$$":
                        // console.log(`Max physical CPU: ${responses[2]}`)
                        preferences["maxcpu-physical"] = responses[2]
                        break
                    case "$$__MAX_CPU_LOGICAL__$$":
                        // console.log(`Max logical CPU: ${responses[2]}`)
                        preferences["maxcpu-logical"] = responses[2]
                        break
                    case "$$__THREADS__$$":
                        // console.log(`Thread count: ${responses[2]}`)
                        preferences["threads"] = responses[2]
                        document.querySelector("#pref_threads").value = preferences["threads"] ? preferences["threads"] : 1
                        break
                }

                savePreferences()

                document.querySelector("#max-cpu-count").innerHTML = `${preferences["maxcpu-physical"]} / ${preferences["maxcpu-logical"]}`
                document.querySelector("#pref_threads").value = preferences["threads"]

                break
            default:
                console.log(`Unknown response: ${evt.data}`)
        }
    }
}

function sendPrompt(sendFirstInput = false) {
    let input = document.querySelector("#inputs")
    const focusTarget = document.querySelector("#inputs")

    if (sendFirstInput) {
        input = document.querySelector("#first-input")
    }

    if (input.value === '') {
        return
    }

    const reflection = document.querySelector("#reflection")
    const antiprompt = document.querySelector("#antiprompt")

    const data = `$$__PROMPT__$$\n$$__SEPARATOR__$$\n${input.value}\n$$__SEPARATOR__$$\n${reflection.value}\n$$__SEPARATOR__$$\n${antiprompt.value}`

    ws.send(data)

    statusRemoveAllClasses()
    statusAddClass('status-running', 'Running')
    buttonStopEnable()

    focusTarget.value = ''
    focusTarget.focus()
}

function applyParameters(sendRequired = true) {
    const datas = []

    datas["N_CTX"] = document.querySelector("#pref_n_ctx").value
    datas["N_BATCH"] = document.querySelector("#pref_n_batch").value

    datas["SAMPLING_METHOD"] = document.querySelector("input[name='pref_sampling_method']:checked").value

    datas["TOP_K"] = document.querySelector("#pref_top_k").value
    datas["TOP_P"] = document.querySelector("#pref_top_p").value
    datas["TEMPERATURE"] = document.querySelector("#pref_temperature").value
    datas["REPEAT_PENALTY"] = document.querySelector("#pref_repeat_penalty").value

    for (const key in datas) {
        preferences[key] = datas[key]
        const data = `$$__PARAMETER__$$\n$$__SEPARATOR__$$\n$$__${key}__$$\n$$__SEPARATOR__$$\n${datas[key]}\n`
        if (sendRequired) {
            ws.send(data)
        }
    }

    const modelFiles = document.querySelector("select[name=model_files]")
    if (modelFiles.length > 0) {
        preferences["MODEL_FILE"] = modelFiles.value
    }

    savePreferences()

    closePanel()

    const input = document.querySelector("#inputs")
    input.focus()
}

function requestModelFiles() {
    const payload = "$$__COMMAND__$$\n$$__SEPARATOR__$$\n$$__MODEL_FILES__$$"
    ws.send(payload)
}

function requestMaxPhycalCPU() {
    const cpuNummsg = "$$__COMMAND__$$\n$$__SEPARATOR__$$\n$$__MAX_CPU_PHYSICAL__$$"
    ws.send(cpuNummsg)
}

function requestMaxLogicalCPU() {
    const cpuNummsg = "$$__COMMAND__$$\n$$__SEPARATOR__$$\n$$__MAX_CPU_LOGICAL__$$"
    ws.send(cpuNummsg)
}

function requestThreads() {
    const threadsmsg = "$$__COMMAND__$$\n$$__SEPARATOR__$$\n$$__THREADS__$$"
    ws.send(threadsmsg)
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
    let promptTEXT = defaultPromptTEXT
    let antipromptTEXT = defaultAntipromptTEXT
    let firstInputTEXT = defaultFirstInputTEXT

    preferences = loadPreferences()

    document.querySelector("#first-input").value = firstInputTEXT

    if (preferences["model_files"] != undefined) {
        document.querySelector("select[name=model_files]").value = preferences["model_files"]
    }

    if (preferences["SAMPLING_METHOD"] != undefined) {
        document.querySelector("input[name='pref_sampling_method'][value='" + preferences["SAMPLING_METHOD"] + "']").checked = true
    }

    keys = ["N_CTX", "N_BATCH", "TOP_K", "TOP_P", "TEMPERATURE", "REPEAT_PENALTY"]
    for (const key of keys) {
        if (preferences[key] != undefined) {
            document.querySelector(`#pref_${key.toLowerCase()}`).value = preferences[key]
            document.querySelector(`#sl_${key.toLowerCase()}`).value = preferences[key]
        }
    }

    if (preferences["darkmode"]) {
        document.querySelector("#switch-shade").click()
    }

    websocketSetup()

    document.querySelector("#reflection").value = promptTEXT
    document.querySelector("#antiprompt").value = antipromptTEXT

    buttonSendEnable()

    document.getElementById("inputs").addEventListener("keydown", function (event) {
        if (event.ctrlKey && event.key == "Enter") {
            event.preventDefault();
            document.getElementById("send").click()
        }
    })

    document.querySelectorAll(".params_slider").forEach(function (element) {
        element.addEventListener("input", function (event) {
            const id = event.target.id
            const value = event.target.value

            const idPrefix = id.split("_")[0]
            const idName = id.split(idPrefix + "_")[1]

            const inputPrefix = "pref_"
            const inputID = inputPrefix + idName

            document.querySelector(`#${inputID}`).value = value
        })
    })

    document.querySelectorAll(".params_value").forEach(function (element) {
        element.addEventListener("input", function (event) {
            const id = event.target.id
            const value = event.target.value

            const idPrefix = id.split("_")[0]
            const idName = id.split(idPrefix + "_")[1]

            const sliderPrefix = "sl_"
            const sliderID = sliderPrefix + idName

            document.querySelector(`#${sliderID}`).value = value
        })
    })
}

document.addEventListener('DOMContentLoaded', init())
