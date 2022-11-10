const markdownEditor = ({
    defaultTab = 'write',
    disabledTab = false,
    caption = null,
    placeholder = null
}) => {
    return {
        activeTab : '', // 'write' | 'preview',
        disabledTab : false,
        body : '',
        markdown : '',
        caption : '',
        placeholder : '',
        async init() {
            this.activeTab = defaultTab
            this.disabledTab = disabledTab
            this.caption = caption
            this.placeholder = placeholder

            this.$watch('body', e => this.markdown = window.marked(e))
        },
        template : `<div class="w-full">
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
                                <li class="right w-full text-right pt-3">
                                    <span x-show="caption != null" class="pl-1 text-xs" x-text="caption"></span>
                                </li>
                            </ul>
                        </nav>
                        <div class="pb-6 pt-1 w-full">
                            <div x-show="activeTab === 'write'">
                                <div>
                                    <textarea :placeholder="placeholder" title="content" class="mt-1 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md h-44"
                                        x-model="body"></textarea>
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