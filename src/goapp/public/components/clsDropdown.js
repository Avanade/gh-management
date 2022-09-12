class clsDropdown {
    template(id, model, data, text, subText, value=text, onClickAction=null){
        return `
        <!-- ****** CO-OWNER CB ****** -->
        <div x-data="{
            ${id}_show: false,
            get filteredData(){
                return !this.${model} ? this.${data} : this.${data}.filter(d => {
                    return d['${text}'].toString().toLowerCase().includes(this.${model}.toLowerCase()) ||
                        d['${subText}'].toString().toLowerCase().includes(this.${model}.toLowerCase())
                })
            }
        }">
            <div class="relative mt-1" >
                <input @click="${id}_show = !${id}_show" x-model="${model}" type="text" class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border-1" role="combobox" aria-controls="options" aria-expanded="false" autocomplete="off">

                <button x-show="${model}" type="button" class="absolute inset-y-0 right-0 flex items-center rounded-r-md px-2 focus:outline-none">
                <svg id="${id}_clearBtn" class="h-5 w-5 text-gray-400" @click="${model} = ''" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path d="M256 8C119 8 8 119 8 256s111 248 248 248 248-111 248-248S393 8 256 8zm0 448c-110.5 0-200-89.5-200-200S145.5 56 256 56s200 89.5 200 200-89.5 200-200 200zm101.8-262.2L295.6 256l62.2 62.2c4.7 4.7 4.7 12.3 0 17l-22.6 22.6c-4.7 4.7-12.3 4.7-17 0L256 295.6l-62.2 62.2c-4.7 4.7-12.3 4.7-17 0l-22.6-22.6c-4.7-4.7-4.7-12.3 0-17l62.2-62.2-62.2-62.2c-4.7-4.7-4.7-12.3 0-17l22.6-22.6c4.7-4.7 12.3-4.7 17 0l62.2 62.2 62.2-62.2c4.7-4.7 12.3-4.7 17 0l22.6 22.6c4.7 4.7 4.7 12.3 0 17z"/></svg>
                </button>

                <ul class="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none sm:text-sm" id="${id}_options" x-show="${id}_show" role="listbox">

                    <template x-for="(i,n) in filteredData">
                    <li class="relative cursor-default select-none py-2 pl-3 pr-9 text-gray-900" :key="n" role="option" @click="${model}=i.${value}; ${id}_show = !${id}_show; ${onClickAction}" tabindex="-1">
                            <div class="flex" >
                                <span class="truncate" x-text="i.${text}"></span>
                                <span class="ml-2 truncate text-gray-500" x-show="!!i.${subText}" x-text="'&nbsp'+i.${subText}"></span>
                            </div>
                    
                            <span x-show="${model}==i.${value}" class="absolute inset-y-0 right-0 flex items-center pr-4 text-indigo-600" >
                            <svg  class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" >
                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                            </svg>
                            </span>


                        </li>
                    </template> 
                        
                </ul>
            </div>
        </div>
        <!-- END OF CO-OWNER CB... -->
        `
    }
}
