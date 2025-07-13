<script lang="ts">
import { Calculate } from '../../wailsjs/go/main/App';
// import { parse, HtmlGenerator } from 'latex.js';

export default {
    data() {
        return {
            scrollPos: 0,
            scrollDelay: false,
            equation: "",
            pos: 0,
            history: [""],
            historyPos: 0,
            resultText: "",
            checked: false,
        }
    },
    mounted() {
        document.addEventListener("keydown", this.input);
        document.getElementById('keyboard-scroll')?.addEventListener("wheel", this.keyboardScroll);
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
            } else if (event.key == 'ArrowUp') {
                if (this.historyPos < this.history.length - 1) {
                    this.historyPos += 1;
                }
                this.equation = this.history[this.historyPos];
                this.pos = this.equation.length;
            } else if (event.key == 'ArrowDown') {
                if (this.historyPos > 0) {
                    this.historyPos -= 1;
                }
                this.equation = this.history[this.historyPos];
                this.pos = this.equation.length;
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
            } else if ([
                '1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
                '+', '-', '*', '/', '^', '%', '(', ')', '=', '.', '^^',
                'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
                ' ',
            ].includes(event.key)) {
                this.enter(event.key);
            }
        },
        enter(char: String) {
            if (char == '^^') {
                char = '^';
            }
            this.equation = [this.equation.slice(0, this.pos), char, this.equation.slice(this.pos)].join('');
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
        },
        exe() {
            Calculate(this.equation).then(result => {
                this.resultText = result
            })
            console.log(this.equation);
            if (this.history[this.history.length - 1] != this.equation) {
                this.history.push(this.equation);
            }
            this.clear();
            this.historyPos = this.history.length - 1;
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
                            <svg width="512mm" height="512mm" viewBox="0 0 512 512" version="1.1" xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg">
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
        </div>


        <div class="main">
            <p>{{ resultText }}</p>
            <div class="input-bar">
                <p id="input" class="text-input">{{ equation }}</p>
                <div class="pos-line"></div>
            </div>
        </div>

        <div class="rsidebar">f(x)=x^2-4</div>

        <div class="bottombar" id="keyboard-container">
            <div class="switch">
                <button @click="keyboardScrollNum" class="switch-pos"><img class="switch-pos-img"
                        src="../assets/images/tauri.svg"></button>
                <button @click="keyboardScrollMath1" class="switch-pos"><img class="switch-pos-img"
                        src="../assets/images/tauri.svg"></button>
                <button @click="keyboardScrollMath2" class="switch-pos"><img class="switch-pos-img"
                        src="../assets/images/tauri.svg"></button>
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
                            <div class="key" @click="enter('sqrt()')">sqrt</div>
                            <div class="key" @click="enter('^')">^</div>
                            <div class="key" @click="enter('log()')">log</div>
                            <div class="key" @click="enter('ln()')">ln</div>
                            <div class="key" @click="enter('=')">=</div>
                            <div class="key" @click="enter('.')">.</div>
                            <div class="key" @click="enter('')"></div>
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
    grid-template-columns: 50px 11fr 4fr;
    grid-template-rows: 1fr 16rem;
    grid-template-areas:
        "lsidebar main rsidebar"
        "bottombar bottombar rsidebar";
}

@media (max-width: 1080px) {
    .grid {
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

    display: flex;
    justify-content: center;

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

    .input-bar {
        position: absolute;

        min-height: 50px;
        max-height: 30%;

        left: 20%;
        right: 20%;
        bottom: 15px;

        border-radius: 10px;
        border: 1px solid var(--surface0);

        background-color: var(--mantle);

        .text-input {
            width: 100%;
            height: 100%;
            padding: 10px;

            font-size: 18px;
        }

        .pos-line {
            position: absolute;
            top: 10px;
            left: 10px;

            translate: calc(1.13ch * v-bind(pos)) 0;

            height: 20px;
            width: 1px;
            background-color: var(--text);
        }
    }
}

.rsidebar {
    grid-area: rsidebar;
    background-color: var(--mantle);
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
