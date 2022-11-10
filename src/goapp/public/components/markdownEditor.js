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
                            <a @click="activeTab = 'write'" :class="{ 'font-bold': activeTab === 'write' }"
                                class="bg-amber-500 hover:bg-amber-700 text-white font-bold py-2 px-4 rounded" href="#">
                                Write
                            </a>
                        </li>
                        <li>
                            <a @click="activeTab = 'preview'" :class="{ 'font-bold': activeTab === 'preview' }"
                                class="bg-amber-500 hover:bg-amber-700 text-white font-bold py-2 px-4 rounded" href="#">
                                Preview
                            </a>
                        </li>
                    </ul>
                </nav>
                <div class="py-6 w-full">
                    <div x-show="activeTab === 'write'">
                        <div>
                            <textarea title="content" class="mt-1 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
                                x-model="body" x-on:input.change="preview"></textarea>
                        </div>
                    </div>
                    <div x-show="activeTab === 'preview'" class="w-full border border-gray-300 rounded-md p-3">
                        <div x-html="markdown" class="preview prose max-w-none prose-img:rounded-md">
                        </div>
                    </div>
                </div>
            </div>`
    }
}