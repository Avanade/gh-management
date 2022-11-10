const markdownEditor = ({
    defaultTab = 'write',
    disabledTab = false
}) => {
    return {
        activeTab : '', // 'write' | 'preview',
        disabledTab : false,
        body : '',
        markdown : '',
        async init() {
            this.activeTab = defaultTab
            this.disabledTab = disabledTab
        },
        preview() { 
            this.markdown = window.marked(this.body) 
        },
        template : `<script src="https://cdn.jsdelivr.net/npm/marked@2.1.3/marked.min.js"></script>
            <div class="p-5 w-full">
                <nav x-show="!disabledTab">
                    <ul class="flex space-x-4">
                        <li>
                            <button @click="activeTab = 'write'" class="px-3 py-2 font-medium text-sm rounded-md"
                                :class="activeTab === 'write' ? 'bg-[#fff2eb] text-[#FF5800]':'text-gray-500 hover:text-gray-700'">
                                Write
                            </button>
                        </li>
                        <li>
                            <button @click="activeTab = 'preview'" class="px-3 py-2 font-medium text-sm rounded-md"
                                :class="activeTab === 'preview' ? 'bg-[#fff2eb] text-[#FF5800]':'text-gray-500 hover:text-gray-700'">
                                Preview
                            </button>
                        </li>
                    </ul>
                </nav>
                <div class="py-6 w-full">
                    <div x-show="activeTab === 'write'">
                        <div>
                            <textarea title="content" class="mt-1 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md h-44"
                                x-model="body" x-on:input.change="preview"></textarea>
                        </div>
                    </div>
                    <div x-show="activeTab === 'preview'" class="w-full border border-gray-300 rounded-md p-3 h-44 overflow-y-auto">
                        <div x-html="markdown" class="preview prose max-w-none prose-img:rounded-md">
                        </div>
                    </div>
                </div>
            </div>`
    }
}