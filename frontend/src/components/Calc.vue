<script lang="ts">
import { Calculate } from '../../wailsjs/go/main/App';
import { EventsOn } from "../../wailsjs/runtime/runtime.js";
// import { parse, HtmlGenerator } from 'latex.js';

const allowed_symbols = [
    '1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
    '+', '-', '*', '/', '^', '%', '(', ')', '=', '.', '^^', '^',
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
    ' ',
]

interface Result {
    Equation : string;
    Error : string;
    ErrorID : number;
}

interface Variable {
    Variable: string,
    Equation: string,
    Dependency: string[],
}

export default {
    data() {
        return {
            scrollPos: 0,
            scrollDelay: false,
            equation: "",
            variables: [] as Variable[],
            pos: 0,
            eq_len: 0,
            history: [""],
            historyPos: 0,
            exchangeHistory: [] as Result[],
            checked: false,
        }
    },
    mounted() {
        document.addEventListener("keydown", this.input);
        document.addEventListener("paste", this.paste)
        document.getElementById('keyboard-scroll')?.addEventListener("wheel", this.keyboardScroll);

        EventsOn('app:variable_defined', (event: Variable) => {
            /* console.log("Defining Variable") */
            let a = -1
            for(let i = 0; i < this.variables.length; i++) {
                if(this.variables[i].Variable == event.Variable) {
                    a = i
                    break
                }
            }
            if (a != -1) {
                this.variables[a] = event
            } else {
                this.variables.push(event)
            }
        });

        EventsOn('app:variable_dropped', (event: Variable) => {
            console.log("Dropping Variable")
            for(var i = 0; i < this.variables.length; i++) {
                if(this.variables[i].Variable == event.Variable) {
                    console.log("Variable to delete: ", event.Variable, " Variable found to delete: ", this.variables[i].Variable, " at: ", i)
                    this.variables.splice(i, 1);
                }
            }
        });
    },
    methods: {
        modelLightSwtich(event: Event) {
            if( this.checked ) {
                document.documentElement.setAttribute("data-theme", "dark")
            } else {
                document.documentElement.setAttribute("data-theme", "light")
            }
        },
        keyboardScrollNum(event: Event) {
            this.scrollPos = 0;
        },
        keyboardScrollMath1(event: Event) {
            this.scrollPos = 1;
        },
        keyboardScrollMath2(event: Event) {
            this.scrollPos = 2;
        },
        keyboardScroll(event: WheelEvent) {
            if (!this.scrollDelay) {
                this.scrollDelay = true;
                if (event.deltaY < -0.1) {
                    if (this.scrollPos > 0) {
                        this.scrollPos -= 1;
                    }
                } else if (event.deltaY > 0.1) {
                    if (this.scrollPos < 2) {
                        this.scrollPos += 1;
                    }
                }
                setTimeout(() => {
                    this.scrollDelay = false;
                }, 300);
            }
        },

        input(event: KeyboardEvent) {
            // console.log(event.key);

            if (event.ctrlKey && (event.key == 'Backspace' || event.key == 'Delete')) {
                this.clear();
            } else if (event.key == 'Backspace') {
                this.remove();
            } else if (event.key == 'Delete') {
                if (this.pos < this.equation.length) {
                    this.pos += 1;
                }
                this.remove();
            } else if (event.ctrlKey && event.key == 'ArrowLeft') {
                this.pos = 0;
            } else if (event.ctrlKey && event.key == 'ArrowRight') {
                this.pos = this.eq_len;
            } else if (event.key == 'ArrowDown') {
                if (this.historyPos < this.history.length - 1) {
                    this.historyPos += 1;
                    this.equation = this.history[this.historyPos];
                    this.pos = this.equation.length;
                    this.eq_len = this.equation.length;
                }
            } else if (event.key == 'ArrowUp') {
                if (this.historyPos > 0) {
                    this.historyPos -= 1;
                    this.equation = this.history[this.historyPos];
                    this.pos = this.equation.length;
                    this.eq_len = this.equation.length;
                }
            } else if (event.key == 'ArrowLeft') {
                if (this.pos > 0) {
                    this.pos -= 1;
                }
            } else if (event.key == 'ArrowRight') {
                if (this.pos < this.equation.length) {
                    this.pos += 1;
                }
            } else if (event.key == 'Enter') {
                if (this.equation.length >= 1) {
                    this.exe();
                }
            } else if (event.key == 'v' && event.ctrlKey) {
                return
            } else if (allowed_symbols.includes(event.key)) {
                this.enter(event.key);
            }
        },
        enter(char: String) {
            if (char == '^^') {
                char = '^';
            } else if (char == ' ' && this.equation[this.pos - 1] == ' ') {                
                return
            }
            this.historyPos = this.history.length;
            this.equation = [this.equation.slice(0, this.pos), char, this.equation.slice(this.pos)].join('');
            this.eq_len = this.equation.length
            this.pos += char.length;
        },
        remove() {
            let a = this.equation.substring(0, this.pos - 1);
            let b = this.equation.substring(this.pos, this.equation.length);
            this.equation = a + b;
            if (this.pos > 0) {
                this.pos -= 1;
            }
        },
        clear() {
            this.equation = '';
            this.pos = 0;
            this.eq_len = 0;
        },
        exe() {
            this.exchangeHistory.push({"Equation": this.equation, "Error": "", "ErrorID":200})
            Calculate(this.equation).then(result => {
                this.exchangeHistory.push(result)
                this.$nextTick(() => {
                    this.scroll();
                });
            })
            if (this.history[this.history.length - 1] != this.equation) {
                this.history.push(this.equation);
            }
            this.clear();
            this.historyPos = this.history.length;
        },
        clearWindow() {
            /* this.history = [];
            this.historyPos = 0; */
            this.exchangeHistory = [];
            this.clear();
        },
        scroll() {
            let scroll_obj = document.getElementById("exchange-scroll-container")
            if (scroll_obj != undefined) {
                scroll_obj.scrollBy({ behavior: 'smooth', top: scroll_obj.scrollHeight });
            } else {
                console.log("Not able to scroll");
            }
        },
        paste(event: ClipboardEvent) {
            event.preventDefault();

            const paste = event.clipboardData?.getData("text");

            if (paste == undefined) {
                return
            }

            for (let i = 0; i < paste.length; i++) {
                if (allowed_symbols.includes(paste.charAt(i))) {
                    this.enter(paste.charAt(i));
                }
            }
        }
    }
}
</script>

<template>
    <div class="grid">
        <div class="lsidebar">
            <label class="theme-slider">
                <input type="checkbox" class="switch-check" name="color-theme-switch" id="color-theme-switch" v-model="checked" @change="modelLightSwtich">
                <div class="color-theme-switch">
                    <span class="slider">
                        <!-- <img class="switch-img" src="../assets/images/Sun.svg"> -->
                        <div class="switch-img">
                            <svg width="512mm" height="512mm" viewBox="0 0 512 512" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg">
                                <g>  
                                    <circle style="fill:none;paint-order:stroke fill markers;stroke-width:35;stroke-dasharray:none" id="path1" cx="256" cy="256" r="100" />
                                    <rect style="stroke:none;paint-order:stroke fill markers" id="rect1" width="35" height="95" x="238.13339" y="13" rx="0" ry="0" />
                                    <rect style="stroke:none;paint-order:stroke fill markers" id="rect1-8" width="35" height="95" x="238.13339" y="404" rx="0" ry="0" />
                                    <rect style="stroke:none;paint-order:stroke fill markers" id="rect1-5" width="35" height="95" x="-273.5" y="13.000002" rx="0" ry="0" transform="rotate(-90)" />
                                    <rect style="stroke:none;paint-order:stroke fill markers" id="rect1-8-5" width="35" height="95" x="-273.5" y="404" rx="0" ry="0" transform="rotate(-90)" />
                                    <rect style="stroke:none;paint-order:stroke fill markers" id="rect1-5-1" width="35" height="95" x="-17.500008" y="119.0387" rx="0" ry="0" transform="rotate(-45)" />
                                    <rect style="stroke:none;paint-order:stroke fill markers" id="rect1-8-5-5" width="35" height="95" x="-17.500008" y="510.03864" rx="0" ry="0" transform="rotate(-45)" />
                                    <rect style="stroke:none;paint-order:stroke fill markers" id="rect1-5-1-1" width="35" height="95" x="-379.53867" y="-242.99997" rx="0" ry="0" transform="matrix(-0.70710678,-0.70710678,-0.70710678,0.70710678,0,0)" /> 
                                    <rect style="stroke:none;paint-order:stroke fill markers" id="rect1-8-5-5-4" width="35" height="95" x="-379.53867" y="147.99997" rx="0" ry="0" transform="matrix(-0.70710678,-0.70710678,-0.70710678,0.70710678,0,0)" />
                                </g>
                            </svg>
                        </div>
                    </span>
                </div>
            </label>
            <div class="placeholder"></div>
            <div class="sidebar-button clear-button bad-action" @click="clearWindow">
                <label for="">
                    <svg width="512mm" height="512mm" viewBox="0 0 512 512" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg">
                        <g xmlns="http://www.w3.org/2000/svg">
                            <rect style="fill:none;fill-opacity:1;paint-order:stroke fill markers;stroke-width:25;stroke-dasharray:none" id="rect1" width="285.50574" height="344.84955" x="111.17564" y="141.96701" rx="28.489468" ry="28.555338"/>
                            <rect style="fill:none;fill-opacity:1;paint-order:stroke fill markers;stroke-width:25;stroke-dasharray:none" id="rect2" width="416.39383" height="70.546059" x="45.731598" y="70.258781" rx="14.022497" ry="13.172651"/>
                            <rect style="fill:none;fill-opacity:1;paint-order:stroke fill markers;stroke-width:25;stroke-dasharray:none" id="rect3" width="232.2858" height="49.038113" x="137.78563" y="20.803186" ry="11.254227"/>
                        </g>
                    </svg>
                </label>
            </div>
        </div>


        <div class="main">
            <div class="exchange-container" id="exchange-scroll-container">
                    <div class="exchange-obj" v-for="item in exchangeHistory">
                        <div class="exchange-obj-container">
                            <p>
                                <!-- {{ item.ErrorID }} : -->
                                <span v-if="item.Error != '' && item.Equation == ''" class="fatalError">
                                    <svg v-if="item.ErrorID < 200" width="512mm" height="512mm" viewBox="0 0 512 512" version="1.1" id="svg1" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg">
                                        <g>
                                            <circle style="fill:none;fill-opacity:1;paint-order:stroke fill markers;stroke-width:50;stroke-dasharray:none" id="path1" cx="256" cy="256" r="228.60632"/>
                                            <rect style="fill-opacity:1;stroke:none;stroke-width:0;stroke-dasharray:none;paint-order:stroke fill markers" id="rect1" width="65.602287" height="350.21713" x="-32.801144" y="186.9301" ry="35.191429" transform="rotate(-45)" rx="32.801144"/>
                                            <rect style="fill-opacity:1;stroke:none;stroke-width:0;stroke-dasharray:none;paint-order:stroke fill markers" id="rect1-8" width="65.602287" height="350.21713" x="-394.83981" y="-175.10857" ry="35.191429" transform="matrix(-0.70710678,-0.70710678,-0.70710678,0.70710678,0,0)" rx="32.801144"/>
                                        </g>
                                    </svg>
                                    {{ item.Error }}
                                </span>
                                <span class="equation">
                                    {{ item.Equation }}
                                    <span v-if="item.Error != '' && item.Equation != ''" class="errorInfo">
                                        <svg width="512mm" height="512mm" viewBox="0 0 512 512" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg">
                                          <g>
                                            <circle style="fill-opacity:1;paint-order:stroke fill markers;fill:none;stroke-width:50;stroke-dasharray:none" id="path1" cx="256" cy="256" r="228.60632"/>
                                            <rect style="fill-opacity:1;stroke:none;stroke-width:0;stroke-dasharray:none;paint-order:stroke fill markers" id="rect1" width="65.602295" height="222.78252" x="223.19885" y="192.94597" ry="35.191429" rx="32.801144"/>
                                            <circle style="fill-opacity:1;stroke:none;stroke-width:0;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers" id="path2" cx="256" cy="145.55412" r="32.5"/>
                                          </g>
                                        </svg>
                                </span>
                            </span>
                        </p>
                    </div>
                    <div class="infoPop" v-if="item.Error != '' && item.Equation != ''">
                        {{ item.Error }}
                    </div> 
                </div>
                <div class="scroll-buffer"></div>
            </div>
            <div class="input-bar-container">
                <div class="input-bar">
                    <p id="input" class="text-input">{{ equation }}</p>
                </div>
            </div>
        </div>

        <div class="rsidebar">
            <div class="rsidebar-scroll-container">
                <p class="rsidebar-heading">Variables</p>
                <div class="variables">
                    <div class="var-obj" v-for="item in variables">
                        {{ item.Variable }} = {{ item.Equation }}
                    </div>
                </div>
                <p class="rsidebar-heading">Functions</p>
                <div class="functions"></div>
            </div>
        </div>

        <div class="bottombar" id="keyboard-container">
            <div class="switch">
                <button @click="keyboardScrollNum" class="switch-pos">
                    <svg class="switch-pos-img" width="512mm" height="512mm" viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg">
                        <g xmlns="http://www.w3.org/2000/svg" inkscape:label="Layer 1" xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape" inkscape:groupmode="layer" id="layer1">
                            <!-- <rect style="fill-opacity:1;paint-order:stroke fill markers;fill:none;stroke-width:30;stroke-dasharray:none" id="rect1" width="400" height="400" x="56" y="56" ry="48.337254"/> -->
                            <text xml:space="preserve" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:304.898px;line-height:1;font-family:'JetBrainsMono Nerd Font';-inkscape-font-specification:'JetBrainsMono Nerd Font';text-align:start;letter-spacing:0px;word-spacing:3.74007px;writing-mode:lr-tb;direction:ltr;text-anchor:start;fill-opacity:1;stroke:none;stroke-width:384.122;stroke-opacity:1;paint-order:stroke fill markers" x="254.56429" y="282.48392" id="text6"><tspan sodipodi:role="line" style="line-height:1;stroke-width:384.124;fill:var(--maroon)" x="254.56429" y="282.48392" id="tspan8">÷</tspan></text>
                            <text xml:space="preserve" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:315.823px;line-height:1;font-family:'JetBrainsMono Nerd Font';-inkscape-font-specification:'JetBrainsMono Nerd Font';text-align:start;letter-spacing:0px;word-spacing:3.87408px;writing-mode:lr-tb;direction:ltr;text-anchor:start;fill-opacity:1;stroke:none;stroke-width:397.887;stroke-opacity:1;paint-order:stroke fill markers" x="251.4447" y="451.51575" id="text6-8"><tspan sodipodi:role="line" style="line-height:1;stroke-width:397.889;fill:var(--sky)" x="251.4447" y="451.51575" id="tspan8-8">×</tspan></text>
                            <text xml:space="preserve" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:315.823px;line-height:1;font-family:'JetBrainsMono Nerd Font';-inkscape-font-specification:'JetBrainsMono Nerd Font';text-align:start;letter-spacing:0px;word-spacing:3.87408px;writing-mode:lr-tb;direction:ltr;text-anchor:start;fill-opacity:1;stroke:none;stroke-width:397.887;stroke-opacity:1;paint-order:stroke fill markers" x="79.671394" y="285.3269" id="text6-8-1"><tspan sodipodi:role="line" style="line-height:1;stroke-width:397.889;fill:var(--sky)" x="79.671394" y="285.3269" id="tspan11">+</tspan></text>
                            <text xml:space="preserve" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:315.823px;line-height:1;font-family:'JetBrainsMono Nerd Font';-inkscape-font-specification:'JetBrainsMono Nerd Font';text-align:start;letter-spacing:0px;word-spacing:3.87408px;writing-mode:lr-tb;direction:ltr;text-anchor:start;fill-opacity:1;stroke:none;stroke-width:397.887;stroke-opacity:1;paint-order:stroke fill markers" x="79.671394" y="451.19992" id="text6-8-1-6"><tspan sodipodi:role="line" style="line-height:1;stroke-width:397.889;fill:var(--maroon)" x="79.671394" y="451.19992" id="tspan11-1">-</tspan></text>
                        </g>
                    </svg>
                </button>
                <button @click="keyboardScrollMath1" class="switch-pos">
                    <svg class="switch-pos-img" width="512mm" height="512mm" viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg">
                        <g id="layer1">
                            <!-- <rect style="fill-opacity:1;paint-order:stroke fill markers;fill:none;stroke-width:30;stroke-dasharray:none" id="rect1" width="400" height="400" x="56" y="56" ry="48.337254"/> -->
                            <text xmlns="http://www.w3.org/2000/svg" xml:space="preserve" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:197.486px;line-height:1;font-family:'JetBrainsMono Nerd Font';-inkscape-font-specification:'JetBrainsMono Nerd Font';text-align:start;letter-spacing:0px;word-spacing:2.42249px;writing-mode:lr-tb;direction:ltr;text-anchor:start;fill-opacity:1;stroke:none;stroke-width:0;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers" x="76.5839" y="234.04414" id="text12"><tspan sodipodi:role="line" id="tspan12" style="fill-opacity:1;stroke:none;stroke-width:0;stroke-dasharray:none;fill:var(--maroon)" x="76.5839" y="234.04414">def</tspan></text>                        
                            <text xmlns="http://www.w3.org/2000/svg" xml:space="preserve" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:197.486px;line-height:1;font-family:'JetBrainsMono Nerd Font';-inkscape-font-specification:'JetBrainsMono Nerd Font';text-align:start;letter-spacing:0px;word-spacing:2.42249px;writing-mode:lr-tb;direction:ltr;text-anchor:start;fill-opacity:1;stroke:none;stroke-width:0;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers" x="76.5839" y="407.20685" id="text12-7"><tspan sodipodi:role="line" id="tspan12-9" style="fill-opacity:1;stroke:none;stroke-width:0;stroke-dasharray:none;fill:var(--sky)" x="76.5839" y="407.20685">sol</tspan></text>                      
                        </g>
                    </svg>
                </button>
                <button @click="keyboardScrollMath2" class="switch-pos">
                    <svg class="switch-pos-img" width="512mm" height="512mm" viewBox="0 0 512 512" version="1.1" id="svg1" inkscape:version="1.4 (86a8ad7, 2024-10-11)" sodipodi:docname="Calc_1.svg" xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape" xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg">
                      <g id="layer1">
                        <!-- <rect style="display:inline;fill:none;fill-opacity:1;paint-order:stroke fill markers;stroke-width:30;stroke-dasharray:none" id="rect1" width="400" height="400" x="56" y="56" ry="48.337254"/> -->
                        <text xml:space="preserve" transform="matrix(1.015818,0,0,1.015818,-565.1203,-658.28699)" id="text12" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:150px;line-height:1;font-family:'JetBrainsMono Nerd Font';-inkscape-font-specification:'JetBrainsMono Nerd Font';text-align:start;letter-spacing:0px;word-spacing:1.84px;writing-mode:lr-tb;direction:ltr;white-space:pre;shape-inside:url(#rect12);display:inline;fill:var(--sky);fill-opacity:1;stroke: none;stroke-width:0;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers"><tspan x="638.08398" y="879.26065" id="tspan4"><tspan dx="0 0 0 0 -1.84" id="tspan3">fn()</tspan></tspan></text>
                        <text xml:space="preserve" transform="matrix(0.26458333,0,0,0.26458333,122.80388,9.3268567)" id="text14-2" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:150px;line-height:1;font-family:'JetBrainsMono Nerd Font';-inkscape-font-specification:'JetBrainsMono Nerd Font';text-align:start;letter-spacing:0px;word-spacing:1.84px;writing-mode:lr-tb;direction:ltr;white-space:pre;shape-inside:url(#rect14-3);display:inline;fill:#4c4f69;fill-opacity:1;stroke:#4c4f69;stroke-width:0;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers"/>
                        <text xml:space="preserve" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:174.383px;line-height:1;font-family:'Times New Roman';-inkscape-font-specification:'Times New Roman, ';text-align:start;letter-spacing:0px;word-spacing:2.1391px;writing-mode:lr-tb;direction:ltr;text-anchor:start;fill:#e64553;fill-opacity:1;stroke:#4c4f69;stroke-width:0;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers" x="106.66274" y="419.96066" id="text15"><tspan sodipodi:role="line" style="fill:var(--maroon);stroke-width:0;-inkscape-font-specification:'Times New Roman, ';font-family:'Times New Roman';font-weight:normal;font-style:normal;font-stretch:normal;font-variant:normal" x="106.66274" y="419.96066" id="tspan17">λ</tspan></text>
                        <text xml:space="preserve" transform="scale(0.26458333)" id="text19" style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-size:150px;line-height:1;font-family:'JetBrainsMono Nerd Font';-inkscape-font-specification:'JetBrainsMono Nerd Font';text-align:start;letter-spacing:0px;word-spacing:1.84px;writing-mode:lr-tb;direction:ltr;white-space:pre;shape-inside:url(#rect19);display:inline;fill:#4c4f69;fill-opacity:1;stroke:#4c4f69;stroke-width:0;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers"/>
                        <rect style="fill:none;fill-opacity:1;stroke:var(--maroon);stroke-width:15.6235;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers" id="rect20" width="113.27856" height="113.27856" x="288.67621" y="275.55365" ry="0"/>
                        <circle style="fill:none;fill-opacity:1;stroke:var(--sky);stroke-width:15.6235;stroke-dasharray:none;stroke-opacity:1;paint-order:stroke fill markers" id="path19" cx="288.25027" cy="369.89761" r="51.567104"/>
                      </g>
                    </svg>
                </button>
            </div>
            <div class="keyboard" id="keyboard-scroll">
                <div class="keys key-1">
                    <div class="keys-container key-split">
                        <div class="k-container num-container">
                            <div class="key" @click="enter('1')">1</div>
                            <div class="key" @click="enter('2')">2</div>
                            <div class="key" @click="enter('3')">3</div>
                            <div class="key" @click="enter('4')">4</div>
                            <div class="key" @click="enter('5')">5</div>
                            <div class="key" @click="enter('6')">6</div>
                            <div class="key" @click="enter('7')">7</div>
                            <div class="key" @click="enter('8')">8</div>
                            <div class="key" @click="enter('9')">9</div>
                            <div class="key" @click="enter('%')">%</div>
                            <div class="key" @click="enter('0')">0</div>
                            <div class="key" @click="enter('.')">.</div>
                        </div>
                        <div class="k-container symbol-container">
                            <div class="key" @click="enter('+')">+</div>
                            <div class="key" @click="enter('-')">-</div>
                            <div class="key" @click="enter('*')">*</div>
                            <div class="key" @click="enter('/')">/</div>
                            <div class="key" @click="enter('(')">(</div>
                            <div class="key" @click="enter(')')">)</div>
                            <div class="key" @click="enter('sqrt(')">sqrt</div>
                            <div class="key" @click="enter('^')">^</div>
                            <div class="key" @click="enter('log(')">log</div>
                            <div class="key" @click="enter('ln(')">ln</div>
                            <div class="key" @click="enter('=')">=</div>
                            <div class="key" @click="enter('.')">.</div>
                            <div class="key" @click="enter(' ')">_</div>
                            <div class="key" @click="remove"><-</div>
                            <div class="key red-key" @click="clear">clr</div>
                            <div class="key blue-key" @click="exe">exe</div>
                        </div>
                    </div>
                </div>
                <div class="keys key-2">
                    <div class="keys-container">
                        <div class="key" @click="enter('define ')">define</div>
                        <div class="key" @click="enter('drop ')">drop</div>
                        <div class="key" @click="enter('solve ')">solve</div>
                    </div>
                </div>
                    <div class="keys key-3">log cos sin tan</div>
            </div>
        </div>
    </div>
</template>

<style lang="css" scoped>
.grid {
    width: 100%;
    height: 100%;
    display: grid;
    grid-template-columns: 50px 11fr 400px;
    grid-template-rows: 1fr 16rem;
    grid-template-areas:
        "lsidebar main rsidebar"
        "bottombar bottombar rsidebar";
}

@media (max-width: 1080px) {
    .grid {
        grid-template-columns: 50px 11fr 13rem;
        grid-template-areas:
            "lsidebar main rsidebar"
            "bottombar bottombar bottombar";
    }
}

.lsidebar {
    grid-area: lsidebar;
    background-color: var(--mantle);
    border: 1px solid var(--crust);

    padding: 0.3rem;

    gap: 1rem;
    display: grid;
    grid-template-rows: 1fr 1fr auto;
    justify-items: center;

    .sidebar-button {
        position: relative;
        cursor: pointer;
        background-color: var(--base);
        border: 1px solid var(--crust);
        border-radius: 10px;
        width: 100%;
        height: 2.5rem;
        text-align: center;

        svg {
            width: 1.5rem;
            stroke: var(--text);
            fill: var(--text);
        }
    }

    /* .bad-action {
        transition-property: background-color, color;
        transition-duration: 0.3s;
    } */

    .bad-action:hover {
        background-color: var(--maroon);
        color: var(--crust);
        stroke: var(--crust);
        fill: var(--crust);
    }

    .theme-slider {
        position: relative;
        width: 25px;
        height: 50px;

        .switch-check { 
            opacity: 0;
            width: 0;
            height: 0;
        }

        .switch-check:checked + .color-theme-switch .switch-img {
            transform: translateY(25px);
        }

        .color-theme-switch {
            position: absolute;
            cursor: pointer;
            outline: 1px solid var(--crust);
            border-radius: 20px;
            background-color: var(--base);
            top: 0;
            bottom: 0;
            left: 0;
            right: 0;

            display: flex;

            .switch-img {
                position: absolute;
                top: 3px;
                left: 3px;
                width: 19px;
                height: 19px;
                transition: transform .4s;
                outline: 1px solid var(--crust);
                border-radius: 50%;
                background-color: var(--mantle);

                svg {
                    width: 100%;
                    height: 100%;
                    stroke: var(--text);
                    fill: var(--text);
                }
            }  
        }
    }
}

.main {
    grid-area: main;
    background-color: var(--base);
    border: 1px solid var(--crust);

    position: relative;


    .exchange-container {
        position: absolute;
        top: 0.2rem;
        bottom: 0;
        left: 0;
        right: 0.2rem;
        overflow-y: scroll;
        padding: 1rem;

        display: flex;
        flex-direction: column;

        .scroll-buffer {
            margin-bottom: 5rem;
        }

        .exchange-obj {
            align-self: flex-start;
            padding: 1rem;
            margin-block: 0.5rem;
            width: fit-content;
            max-width: 100%;
            border: 1px solid var(--crust);
            border-radius: 10px;
            background-color: var(--mantle);


            span {
                word-wrap: break-word;
            }

            .errorInfo, .fatalError {
                /* display: inline-flex;
                align-items: center; */
                svg {
                    position: relative;
                    top: 2px;
                    width: 1rem;
                    display: inline;
                }
            }

            .fatalError {
                svg {
                    /* margin-right: 0.5rem; */
                    fill: var(--red);
                    stroke: var(--red);
                }
            }

            .errorInfo {
                position: relative;
                fill: var(--yellow);
                stroke: var(--yellow);
                /* margin-left: 0.5rem; */
            }


        }


        .exchange-obj:hover > .infoPop {
            display: block;
            position: absolute;
            background-color: var(--base);
            border: 1px solid var(--crust);
            border-radius: 10px;
            padding: 1rem;
            height: min-content;
            /* max-height: 5rem; */
            width: max-content;
            max-width: 50dvw;
        }

        .infoPop {
            display: none;
        }

        .exchange-obj:nth-child(odd) { 
            align-self: flex-end;
        }
    
    }

    .exchange-container::-webkit-scrollbar {
        margin-left: 1rem;
        width: 0.5rem;
    }

        /* Track */
    .exchange-container::-webkit-scrollbar-track {
        background: transparent;
    }

        /* Handle */
    .exchange-container::-webkit-scrollbar-thumb {
        border-radius: 10px;
        background: var(--crust);
    }


    .input-bar-container {
        position: absolute;
        left: 0;
        right: 0;
        bottom: 1rem;

        display: flex;
        justify-content: center;
    }

    .input-bar {
        min-height: 44px;
        width: 50ch;
        border-radius: 10px;
        border: 1px solid var(--surface0);
        
        background-color: var(--mantle);
        
        .text-input {
            content: "";
            word-wrap: break-word;

            --parent-width: max-width;

            /* min-height: 20px; */
            height: 100%;
            padding: 10px;

            line-break: anywhere;

            font-size: 18px;

            position: relative;
        }
        .text-input:empty:before {
            content:"";
            display:inline-block;
        }

        .text-input::after {
            content: '';
            position: absolute;
            bottom: 12px;
            left: 10px;

            translate: calc(1ch * mod(v-bind(pos), 42)) 
                calc(
                    2ch * (
                        round(
                            down, 
                            (
                                v-bind(pos) - (v-bind(eq_len) - mod(
                                    v-bind(eq_len), 42
                                    ))
                            )/ 42
                        )
                    )
                );
            transition: translate 0.1s ease-out;

            height: 20px;
            width: 1px;
            background-color: var(--blue);
        }
    }

    @media (max-width: 720px) {
        .input-bar {
            width: 30ch;

            .text-input::after {
                translate: calc(1ch * mod(v-bind(pos), 24)) 
                    calc(2ch * (round(down,(v-bind(pos) - (v-bind(eq_len) - mod(v-bind(eq_len), 24)))/ 24)));
            }
        }
    }

    @media (min-width: 1080px) {
        .input-bar {
            width: 70ch;

            .text-input::after {
                translate: calc(1ch * mod(v-bind(pos), 60)) 
                    calc(2ch * (round(down,(v-bind(pos) - (v-bind(eq_len) - mod(v-bind(eq_len), 60)))/ 60)));
            }
        }
    }
}

.rsidebar {
    grid-area: rsidebar;
    background-color: var(--mantle);
    border: 1px solid var(--crust);
    width: 100%;
    overflow-y: scroll;

    .rsidebar-scroll-container {
        width: 100%;
        max-width: inherit;
        display: flex;
        flex-direction: column;
        /* grid-template-rows: auto 1fr auto 1fr; */

        .rsidebar-heading {
            padding: 0.5rem;
            text-align: center;
            border-top: 1px solid var(--crust);
        }

        .variables {
            width: 100%;
            display: flex;
            flex-direction: column;
            padding: 0.2rem;
            gap: 0.2rem;

            
            .var-obj {
                overflow: hidden;
                background-color: var(--base);
                border: 1px solid var(--crust);
                border-radius: 10px;
                padding: 1rem;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
        }
    }
}

.rsidebar::-webkit-scrollbar {
    margin-right: 1rem;
    width: 0.5rem;
}

.rsidebar::-webkit-scrollbar-track {
    background: transparent;
}

.rsidebar::-webkit-scrollbar-thumb {
    border-radius: 10px;
    background: var(--base);
    border: 1px solid var(--crust);
}

.bottombar {
    grid-area: bottombar;
    background-color: var(--mantle);
    border: 1px solid var(--crust);

    display: grid;
    grid-template-columns: 50px 1fr;

    .switch {
        align-self: center;
        justify-self: center;

        height: 120px;
        width: 40px;

        border-radius: 10px;
        border: 1px solid var(--surface0);
        background-color: var(--crust);

        position: relative;

        display: flex;
        flex-direction: column;
        justify-content: space-between;
        align-items: center;
        padding: 2px;

        &::after {
            position: absolute;
            content: '';
            z-index: 1;

            width: 34px;
            height: 34px;

            border-radius: 8px;

            background-color: var(--base);


            translate: 0 calc(40px * v-bind(scrollPos));
            transition: translate 0.5s ease-in-out;

        }

        .switch-pos {
            z-index: 2;
            margin: 2px;
            width: 30px;
            height: 30px;
            padding: 0;
            background-color: transparent;
            border: none;
            outline: none;
        }

        .switch-pos-img {
            /* fill: var(--text);
            stroke: var(--text);; */
            width: 30px;
            height: 30px;
        }
    }

    .keyboard {
        max-width: 100%;
        max-height: fit-content;

        overflow: hidden;

        width: 100%;
        height: 100%;

        .keys {
            width: 100%;
            height: 100%;

            translate: 0 calc(-100% * v-bind(scrollPos));
            transition: translate 0.5s ease-in-out;

            position: relative;

            .keys-container {
                width: 100%;
                height: 100%;
                position: absolute;

                display: grid;
                grid-template-rows: 1fr 1fr 1fr 1fr;
                grid-template-columns: 1fr 1fr 1fr 1fr 1fr;

                padding: 1rem;
                gap: 1rem;
                
                .key {
                    width: 100%;
                    height: 100%;
                    display: flex;
                    align-items: center;
                    justify-content: center;

                    cursor: pointer;

                    background-color: var(--base);

                    border-radius: 10px;
                    border: 1px solid var(--surface0);

                    user-select: none;

                    font-size: 20px;
                }

                .blue-key {
                    background-color: var(--sky);
                    color: var(--base);
                }

                .red-key {
                    background-color: var(--maroon);
                    color: var(--base);
                }
            }

            .key-split {
                grid-template-columns: 3fr 4fr;
                grid-template-rows: 1fr;
                padding: 0;
                gap: 0;

                .num-container {
                    height: 100%;
                    width: 100%;
                    grid-template-rows: 1fr 1fr 1fr 1fr;
                    grid-template-columns: 1fr 1fr 1fr;
                }

                .symbol-container {
                    height: 100%;
                    width: 100%;
                    grid-template-rows: 1fr 1fr 1fr 1fr;
                    grid-template-columns: 1fr 1fr 1fr 1fr;
                }

                .k-container {
                    padding: 1rem;
                    gap: 1rem;
                    display: grid;
                }
            }
        }
    }
}
</style>
